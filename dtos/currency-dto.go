package dtos

import "github.com/shopspring/decimal"

type CurrencyDto struct {
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	LogoImageURL     string `json:"logoImageUrl"`
	SupportsPayments bool   `json:"supportsPayments"`
}

type GetExchangeParams struct {
	Amount       decimal.Decimal `json:"amount"`
	FromCurrency string          `json:"fromCurrency"`
	ToCurrencies []string        `json:"toCurrencies"`
}

type GetExchangeResult map[string]decimal.Decimal
