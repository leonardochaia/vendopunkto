package dtos

import "github.com/shopspring/decimal"

type AtomicUnitDTO struct {
	Value          decimal.Decimal `json:"value"`
	ValueFormatted string          `json:"valueFormatted"`
}

func NewAtomicUnitDTO(u decimal.Decimal) AtomicUnitDTO {
	return AtomicUnitDTO{
		Value:          u,
		ValueFormatted: u.String(),
	}
}
