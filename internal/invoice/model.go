package invoice

import (
	"time"

	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/leonardochaia/vendopunkto/unit"
)

type InvoiceStatus int

const (
	// Pending when the invoice has been created but the wallet has not
	// received payment yet or the payment amount is not enough
	Pending InvoiceStatus = 1
	// Completed when the invoice has been payid completely.
	Completed InvoiceStatus = 2
	// Failed when the payments fails for whatever reason
	Failed InvoiceStatus = 3
)

type PaymentStatus int

const (
	// Mempool when the wallet has received the payment, and it's waiting
	// to be included in a block.
	Mempool PaymentStatus = 1
	// Confirmed when the wallet has received enough confirmations.
	Confirmed PaymentStatus = 2
)

type Invoice struct {
	ID             string
	Total          unit.AtomicUnit
	Currency       string
	CreatedAt      time.Time
	PaymentMethods []*PaymentMethod
}

type PaymentMethod struct {
	ID        uint
	InvoiceID string
	Total     unit.AtomicUnit
	Currency  string
	Address   string
	Payments  []*Payment
}

type Payment struct {
	ID              uint
	TxHash          string
	PaymentMethodID uint
	Amount          unit.AtomicUnit
	Confirmations   uint64
	BlockHeight     uint64
	ConfirmedAt     time.Time
	CreatedAt       time.Time
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

	payed := invoice.CalculateTotalPayedAmount().Float64()

	return (payed * 100) / invoice.Total.Float64()
}

// CalculateTotalPayedAmount returns the total amount, of all payments
// converted to the invoice's currency
func (invoice *Invoice) CalculateTotalPayedAmount() unit.AtomicUnit {
	var total unit.AtomicUnit
	total = 0

	for _, method := range invoice.PaymentMethods {
		for _, payment := range method.Payments {
			if payment.Status() == Confirmed {
				// payment amount converted to invoice's currency
				total += convertCurrencyWithTotals(method.Total, invoice.Total, payment.Amount)
			}
		}
	}

	return total
}

// CalculateRemainingAmount returns how much is needed to fully pay this invoice
// in the invoice's currency
func (invoice *Invoice) CalculateRemainingAmount() unit.AtomicUnit {
	totalPayed := invoice.CalculateTotalPayedAmount()

	if totalPayed > invoice.Total {
		return 0
	}

	return invoice.Total - totalPayed
}

// CalculatePaymentMethodRemaining returns how much is remaining in the method's
// currency to fully pay this invoice
func (invoice Invoice) CalculatePaymentMethodRemaining(method PaymentMethod) unit.AtomicUnit {

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
	amount unit.AtomicUnit,
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
	amount unit.AtomicUnit,
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
	aTotal unit.AtomicUnit,
	bTotal unit.AtomicUnit,
	aAmount unit.AtomicUnit) unit.AtomicUnit {

	// exchange rate of invoice's currency to method's currency
	exchangeRate := aTotal.Float64() / bTotal.Float64()

	// the amount converted to invoice's currency
	converted := aAmount.Float64() / exchangeRate

	return unit.NewFromFloat(converted)
}

func convertInvoiceToDto(invoice Invoice, pluginMgr *pluginmgr.Manager) (dtos.InvoiceDto, error) {

	dto := &dtos.InvoiceDto{
		ID:                invoice.ID,
		Total:             dtos.NewAtomicUnitDTO(invoice.Total),
		Currency:          invoice.Currency,
		CreatedAt:         invoice.CreatedAt,
		Status:            uint(invoice.Status()),
		PaymentPercentage: invoice.CalculatePaymentPercentage(),
		Remaining:         dtos.NewAtomicUnitDTO(invoice.CalculateRemainingAmount()),
		PaymentMethods:    []*dtos.PaymentMethodDto{},
		Payments:          []*dtos.PaymentDto{},
	}

	for _, method := range invoice.PaymentMethods {

		var qrCode string
		if method.Address != "" {
			info, err := pluginMgr.GetWalletInfoForCurrency(method.Currency)
			if err != nil {
				return dtos.InvoiceDto{}, err
			}
			qrCode, err = info.BuildQRCode(method.Address, method.Total)
			if err != nil {
				return dtos.InvoiceDto{}, err
			}
		}

		methodDto := &dtos.PaymentMethodDto{
			ID:        method.ID,
			Total:     dtos.NewAtomicUnitDTO(method.Total),
			Currency:  method.Currency,
			Address:   method.Address,
			Remaining: dtos.NewAtomicUnitDTO(invoice.CalculatePaymentMethodRemaining(*method)),
			QRCode:    qrCode,
		}

		for _, payment := range method.Payments {
			paymentDto := &dtos.PaymentDto{
				TxHash:        payment.TxHash,
				Amount:        dtos.NewAtomicUnitDTO(payment.Amount),
				Confirmations: payment.Confirmations,
				ConfirmedAt:   payment.ConfirmedAt,
				CreatedAt:     payment.CreatedAt,
				Status:        uint(payment.Status()),
				Currency:      method.Currency,
				BlockHeight:   payment.BlockHeight,
			}
			dto.Payments = append(dto.Payments, paymentDto)
		}

		dto.PaymentMethods = append(dto.PaymentMethods, methodDto)
	}

	return *dto, nil
}
