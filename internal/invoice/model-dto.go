package invoice

import (
	"time"

	"github.com/leonardochaia/vendopunkto/unit"
)

type InvoiceDto struct {
	ID                string              `json:"id"`
	Total             unit.AtomicUnit     `json:"total"`
	Currency          string              `json:"currency"`
	CreatedAt         time.Time           `json:"createdAt"`
	PaymentMethods    []*PaymentMethodDto `json:"paymentMethods"`
	Status            InvoiceStatus       `json:"status"`
	PaymentPercentage float64             `json:"paymentPercentage"`
	Remaining         unit.AtomicUnit     `json:"remaining"`
	Payments          []*PaymentDto       `json:"payments"`
}

type PaymentMethodDto struct {
	ID        uint            `json:"id"`
	Total     unit.AtomicUnit `json:"total"`
	Currency  string          `json:"currency"`
	Address   string          `json:"address"`
	Remaining unit.AtomicUnit `json:"remaining"`
}

type PaymentDto struct {
	TxHash        string          `json:"txHash"`
	Amount        unit.AtomicUnit `json:"amount"`
	Currency      string          `json:"currency"`
	Confirmations uint64          `json:"confirmations"`
	ConfirmedAt   time.Time       `json:"confirmedAt"`
	CreatedAt     time.Time       `json:"createdAt"`
	Status        PaymentStatus   `json:"status"`
}

func (invoice Invoice) ToDto() (InvoiceDto, error) {

	dto := &InvoiceDto{
		ID:                invoice.ID,
		Total:             invoice.Total,
		Currency:          invoice.Currency,
		CreatedAt:         invoice.CreatedAt,
		Status:            invoice.Status(),
		PaymentPercentage: invoice.CalculatePaymentPercentage(),
		Remaining:         invoice.CalculateRemainingAmount(),
		PaymentMethods:    []*PaymentMethodDto{},
		Payments:          []*PaymentDto{},
	}

	for _, method := range invoice.PaymentMethods {
		methodDto := &PaymentMethodDto{
			ID:        method.ID,
			Total:     method.Total,
			Currency:  method.Currency,
			Address:   method.Address,
			Remaining: invoice.CalculatePaymentMethodRemaining(*method),
		}

		for _, payment := range method.Payments {
			paymentDto := &PaymentDto{
				TxHash:        payment.TxHash,
				Amount:        payment.Amount,
				Confirmations: payment.Confirmations,
				ConfirmedAt:   payment.ConfirmedAt,
				CreatedAt:     payment.CreatedAt,
				Status:        payment.Status(),
				Currency:      method.Currency,
			}
			dto.Payments = append(dto.Payments, paymentDto)
		}

		dto.PaymentMethods = append(dto.PaymentMethods, methodDto)
	}

	return *dto, nil
}
