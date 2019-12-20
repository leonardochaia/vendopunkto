package dtos

import "github.com/leonardochaia/vendopunkto/unit"

type AtomicUnitDTO struct {
	Value          uint64 `json:"value"`
	ValueFormatted string `json:"valueFormatted"`
}

func NewAtomicUnitDTO(u unit.AtomicUnit) AtomicUnitDTO {
	return AtomicUnitDTO{
		Value:          uint64(u),
		ValueFormatted: u.Formatted(),
	}
}
