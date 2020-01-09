package dtos

import (
	"time"

	"github.com/shopspring/decimal"
)

type InvoiceCreationParams struct {
	Total          decimal.Decimal               `json:"total"`
	Currency       string                        `json:"currency"`
	PaymentMethods []PaymentMethodCreationParams `json:"paymentMethods"`
}

type PaymentMethodCreationParams struct {
	Currency string          `json:"currency"`
	Total    decimal.Decimal `json:"total"`
}

type InvoiceConfirmPaymentsParams struct {
	TxHash        string          `json:"txHash"`
	Address       string          `json:"address"`
	Amount        decimal.Decimal `json:"amount"`
	Confirmations uint64          `json:"confirmations"`
	BlockHeight   uint64          `json:"blockHeight"`
}

type InvoiceGeneratePaymentMethodAddressParams struct {
	Currency string `json:"currency"`
}

type InvoiceDto struct {
	ID                string              `json:"id"`
	Total             decimal.Decimal     `json:"total"`
	Currency          string              `json:"currency"`
	CreatedAt         time.Time           `json:"createdAt"`
	PaymentMethods    []*PaymentMethodDto `json:"paymentMethods"`
	Status            uint                `json:"status"`
	PaymentPercentage float64             `json:"paymentPercentage"`
	Remaining         decimal.Decimal     `json:"remaining"`
	Payments          []*PaymentDto       `json:"payments"`
}

type PaymentMethodDto struct {
	ID        uint            `json:"id"`
	Total     decimal.Decimal `json:"total"`
	Currency  string          `json:"currency"`
	Address   string          `json:"address"`
	Remaining decimal.Decimal `json:"remaining"`
	QRCode    string          `json:"qrCode"`
}

type PaymentDto struct {
	TxHash        string          `json:"txHash"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	Confirmations uint64          `json:"confirmations"`
	BlockHeight   uint64          `json:"blockHeight"`
	ConfirmedAt   time.Time       `json:"confirmedAt"`
	CreatedAt     time.Time       `json:"createdAt"`
	Status        uint            `json:"status"`
}
