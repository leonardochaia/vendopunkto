package vendopunkto

import (
	"context"
	"time"

	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/shopspring/decimal"
)

// InvoiceStatus determines the current status for the invoice
type InvoiceStatus int

const (
	// Pending when the invoice has been created but the wallet has not
	// received payment yet or the payment amount is not enough
	Pending InvoiceStatus = 1
	// Completed when the invoice has been payed completely.
	Completed InvoiceStatus = 2
	// Failed when the payments fails for whatever reason
	Failed InvoiceStatus = 3
)

// PaymentStatus of the payment
type PaymentStatus int

const (
	// Mempool when the wallet has received the payment, and it's waiting
	// to be included in a block.
	Mempool PaymentStatus = 1
	// Confirmed when the wallet has received enough confirmations.
	Confirmed PaymentStatus = 2
)

type Invoice struct {
	ID             string          `sql:",pk"`
	Total          decimal.Decimal `sql:",notnull"`
	Currency       string          `sql:",notnull"`
	CreatedAt      time.Time       `sql:",notnull"`
	PaymentMethods []*PaymentMethod
}

type PaymentMethod struct {
	ID        uint            `sql:",pk"`
	InvoiceID string          `sql:",notnull"`
	Total     decimal.Decimal `sql:",notnull"`
	Currency  string          `sql:",notnull"`
	Address   string          `sql:",notnull"`
	Payments  []*Payment
}

type Payment struct {
	ID              uint            `sql:",pk"`
	TxHash          string          `sql:",notnull"`
	PaymentMethodID uint            `sql:",notnull"`
	Amount          decimal.Decimal `sql:",notnull"`
	Confirmations   uint64
	BlockHeight     uint64
	ConfirmedAt     time.Time
	CreatedAt       time.Time `sql:",notnull"`
}

func (invoice *Invoice) Status() InvoiceStatus {
	p := invoice.CalculatePaymentPercentage()
	if p >= 100 {
		return Completed
	}

	return Pending
}

// CalculatePaymentPercentage returns how much of this invoice has been payed
// in percentage.
// Given that an invoice could be payed using multiple currencies, this tells us
// what percentage of the invoice has been payed
func (invoice *Invoice) CalculatePaymentPercentage() float64 {

	h := decimal.NewFromInt(100)
	payed := invoice.CalculateTotalPayedAmount()

	f, _ := payed.Mul(h).Div(invoice.Total).Float64()
	return f
}

// CalculateTotalPayedAmount returns the total amount, of all payments
// converted to the invoice's currency
func (invoice *Invoice) CalculateTotalPayedAmount() decimal.Decimal {
	total := decimal.Zero

	for _, method := range invoice.PaymentMethods {
		for _, payment := range method.Payments {
			if payment.Status() == Confirmed {
				// payment amount converted to invoice's currency
				converted := convertCurrencyWithTotals(method.Total, invoice.Total, payment.Amount)
				total = decimal.Sum(total, converted)
			}
		}
	}

	return total
}

// CalculateRemainingAmount returns how much is needed to fully pay this invoice
// in the invoice's currency
func (invoice *Invoice) CalculateRemainingAmount() decimal.Decimal {
	totalPayed := invoice.CalculateTotalPayedAmount()

	if totalPayed.GreaterThan(invoice.Total) {
		return decimal.Zero
	}

	return invoice.Total.Sub(totalPayed)
}

// CalculatePaymentMethodRemaining returns how much is remaining in the method's
// currency to fully pay this invoice
func (invoice Invoice) CalculatePaymentMethodRemaining(method PaymentMethod) decimal.Decimal {

	remainingInInvoiceCurrency := invoice.CalculateRemainingAmount()

	// convert it to method's currency
	return convertCurrencyWithTotals(invoice.Total, method.Total, remainingInInvoiceCurrency)
}

func (invoice *Invoice) FindPaymentMethodForAddress(address string) *PaymentMethod {
	if invoice.PaymentMethods != nil {
		for _, method := range invoice.PaymentMethods {
			if method.Address == address {
				return method
			}
		}
	}
	return nil
}

func (invoice *Invoice) FindPaymentMethodForCurrency(currency string) *PaymentMethod {
	if invoice.PaymentMethods != nil {
		for _, method := range invoice.PaymentMethods {
			if method.Currency == currency {
				return method
			}
		}
	}
	return nil
}

func (invoice *Invoice) FindDefaultPaymentMethod() *PaymentMethod {
	return invoice.FindPaymentMethodForCurrency(invoice.Currency)
}

func (invoice *Invoice) AddPaymentMethod(
	currency string,
	address string,
	amount decimal.Decimal,
) *PaymentMethod {

	method := &PaymentMethod{
		InvoiceID: invoice.ID,
		Currency:  currency,
		Address:   address,
		Total:     amount,
	}

	invoice.PaymentMethods = append(invoice.PaymentMethods, method)

	return method
}

func (payment *Payment) Status() PaymentStatus {
	if payment.Confirmations > 0 {
		return Confirmed
	}
	return Mempool
}

func (payment *Payment) Update(confirmations uint64, blockHeight uint64) bool {

	if confirmations == payment.Confirmations && blockHeight == payment.BlockHeight {
		return false
	}

	if confirmations > 0 && payment.Confirmations == 0 {
		payment.ConfirmedAt = time.Now()
	}

	// Wallet always win. In order to support reorgs.
	payment.Confirmations = confirmations
	payment.BlockHeight = blockHeight
	return true
}

func (method *PaymentMethod) AddPayment(
	txHash string,
	amount decimal.Decimal,
	confirmations uint64,
	blockHeight uint64,
) *Payment {

	payment := &Payment{
		TxHash:          txHash,
		PaymentMethodID: method.ID,
		Amount:          amount,
		Confirmations:   confirmations,
		BlockHeight:     blockHeight,
	}

	method.Payments = append(method.Payments, payment)

	if payment.Confirmations > 0 {
		payment.ConfirmedAt = time.Now()
	}

	return payment
}

func (method *PaymentMethod) FindPayment(txHash string) *Payment {
	for _, payment := range method.Payments {
		if payment.TxHash == txHash {
			return payment
		}
	}
	return nil
}

// convertCurrencyWithTotals returns the conversion of the provided amount to the
// invoice's currency.
func convertCurrencyWithTotals(
	aTotal decimal.Decimal,
	bTotal decimal.Decimal,
	aAmount decimal.Decimal) decimal.Decimal {

	// exchange rate of invoice's currency to method's currency
	exchangeRate := aTotal.Div(bTotal)

	// the amount converted to invoice's currency
	converted := aAmount.Div(exchangeRate)

	return converted
}

type InvoiceFilter struct {
}

// InvoiceManager is the logic abstraction
type InvoiceManager interface {
	GetInvoice(ctx context.Context, id string) (*Invoice, error)
	Search(ctx context.Context, filter InvoiceFilter) ([]Invoice, error)
	GetInvoiceByAddress(ctx context.Context, address string) (*Invoice, error)
	CreateInvoice(ctx context.Context, params dtos.InvoiceCreationParams) (*Invoice, error)
	CreateAddressForPaymentMethod(ctx context.Context, invoiceID string,
		currency string) (*Invoice, error)
	ConfirmPayment(ctx context.Context, address string, confirmations uint64,
		amount decimal.Decimal, txHash string, blockHeight uint64) (*Invoice, error)
}

// InvoiceRepository is the abstraction for handling Invoices database
type InvoiceRepository interface {
	Search(ctx context.Context, filter InvoiceFilter) ([]Invoice, error)
	FindByID(ctx context.Context, id string) (*Invoice, error)
	FindByAddress(ctx context.Context, address string) (*Invoice, error)
	Create(ctx context.Context, invoice *Invoice) error
	UpdatePaymentMethod(ctx context.Context, method *PaymentMethod) error
	CreatePayment(ctx context.Context, payment *Payment) error
	UpdatePayment(ctx context.Context, payment *Payment) error

	// GetMaxBlockHeightForCurrencies returns a map of currencies and the last
	// known block height. It must return all currencies that are awaiting
	// payment, even if they don't have a known last block height.
	GetMaxBlockHeightForCurrencies(ctx context.Context) (map[string]uint64, error)
}

// InvoiceTopic is a pub-sub mechanism where consumers can Register to
// receive messages sent to using Send.
// credits to https://github.com/tv42/topic/blob/master/topic.go
type InvoiceTopic interface {
	Register(invoiceID string) <-chan Invoice
	Unregister(invoiceID string)
	Send(msg Invoice)
	Close()
}
