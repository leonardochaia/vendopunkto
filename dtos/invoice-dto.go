package dtos

import (
	"time"

	"github.com/leonardochaia/vendopunkto/unit"
)

type InvoiceCreationParams struct {
	Total          unit.AtomicUnit `json:"total"`
	Currency       string          `json:"currency"`
	PaymentMethods []string        `json:"paymentMethods"`
}

type InvoiceConfirmPaymentsParams struct {
	TxHash        string          `json:"txHash"`
	Address       string          `json:"address"`
	Amount        unit.AtomicUnit `json:"amount"`
	Confirmations uint64          `json:"confirmations"`
	BlockHeight   uint64          `json:"blockHeight"`
}

type InvoiceGeneratePaymentMethodAddressParams struct {
	Currency string `json:"currency"`
}

type InvoiceDto struct {
	ID                string              `json:"id"`
	Total             AtomicUnitDTO       `json:"total"`
	Currency          string              `json:"currency"`
	CreatedAt         time.Time           `json:"createdAt"`
	PaymentMethods    []*PaymentMethodDto `json:"paymentMethods"`
	Status            uint                `json:"status"`
	PaymentPercentage float64             `json:"paymentPercentage"`
	Remaining         AtomicUnitDTO       `json:"remaining"`
	Payments          []*PaymentDto       `json:"payments"`
}

type PaymentMethodDto struct {
	ID        uint          `json:"id"`
	Total     AtomicUnitDTO `json:"total"`
	Currency  string        `json:"currency"`
	Address   string        `json:"address"`
	Remaining AtomicUnitDTO `json:"remaining"`
	QRCode    string        `json:"qrCode"`
}

type PaymentDto struct {
	TxHash        string        `json:"txHash"`
	Amount        AtomicUnitDTO `json:"amount"`
	Currency      string        `json:"currency"`
	Confirmations uint64        `json:"confirmations"`
	BlockHeight   uint64        `json:"blockHeight"`
	ConfirmedAt   time.Time     `json:"confirmedAt"`
	CreatedAt     time.Time     `json:"createdAt"`
	Status        uint          `json:"status"`
}
