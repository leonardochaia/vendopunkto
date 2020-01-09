package dtos

import "github.com/shopspring/decimal"

type ExchangeRatesCurrency struct {
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	SupportsPayments bool   `json:"supportsPayments"`
}

type GetCurrencyRatesParams struct {
	FromCurrency string   `json:"fromCurrency"`
	ToCurrencies []string `json:"toCurrencies"`
}

type GetExchangeParams struct {
	Amount       decimal.Decimal `json:"amount"`
	FromCurrency string          `json:"fromCurrency"`
	ToCurrencies []string        `json:"toCurrencies"`
}

type GetExchangeResult map[string]decimal.Decimal
