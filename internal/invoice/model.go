package invoice

import (
	"encoding/json"
	"time"
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
	ID             string    `json:"id"`
	Amount         uint64    `json:"amount" gorm:"type:BIGINT"`
	Currency       string    `json:"currency"`
	PaymentAddress string    `json:"address" gorm:"index:inv_addr"`
	CreatedAt      time.Time `json:"createdAt"`
	Payments       []Payment `json:"payments"`
}

type Payment struct {
	InvoiceID     string    `json:"invoiceId" gorm:"index"`
	TxHash        string    `json:"txHash" gorm:"primary_key"`
	Amount        uint64    `json:"amount" gorm:"type:BIGINT"`
	Confirmations uint      `json:"confirmations"`
	ConfirmedAt   time.Time `json:"confirmedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (invoice *Invoice) Status() InvoiceStatus {
	if invoice.Payments != nil {
		confirmedAmount := invoice.SumConfirmedAmount()
		if confirmedAmount >= invoice.Amount {
			return Completed
		}
	}

	return Pending
}

func (invoice *Invoice) SumConfirmedAmount() uint64 {
	var confirmedAmount uint64
	confirmedAmount = 0
	if invoice.Payments != nil {
		for _, payment := range invoice.Payments {
			if payment.Status() == Confirmed {
				confirmedAmount += payment.Amount
			}
		}
	}
	return confirmedAmount
}

func (invoice *Invoice) FindPayment(txHash string) *Payment {
	if invoice.Payments != nil {
		for _, payment := range invoice.Payments {
			if payment.TxHash == txHash {
				return &payment
			}
		}
	}
	return nil
}

func (invoice *Invoice) AddPayment(
	txHash string,
	amount uint64,
	confirmations uint,
) *Payment {

	payment := &Payment{
		TxHash:        txHash,
		InvoiceID:     invoice.ID,
		Amount:        amount,
		Confirmations: confirmations,
	}

	invoice.Payments = append(invoice.Payments, *payment)

	if payment.Confirmations > 0 {
		payment.ConfirmedAt = time.Now()
	}

	return payment
}

func (invoice *Invoice) MarshalJSON() ([]byte, error) {
	type Alias Invoice
	return json.Marshal(&struct {
		Status InvoiceStatus `json:"status"`
		*Alias
	}{
		Status: invoice.Status(),
		Alias:  (*Alias)(invoice),
	})
}

func (payment *Payment) MarshalJSON() ([]byte, error) {
	type Alias Payment
	return json.Marshal(&struct {
		Status PaymentStatus `json:"status"`
		*Alias
	}{
		Status: payment.Status(),
		Alias:  (*Alias)(payment),
	})
}

func (payment *Payment) Status() PaymentStatus {
	if payment.Confirmations > 0 {
		return Confirmed
	}
	return Mempool
}

func (payment *Payment) Update(confirmations uint) {

	if confirmations > 0 && payment.Confirmations == 0 {
		payment.ConfirmedAt = time.Now()
	}

	// Wallet always win. In order to support reorgs.
	payment.Confirmations = confirmations
}
