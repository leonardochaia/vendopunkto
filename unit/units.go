package unit

import (
	"strconv"
)

type AtomicUnit uint64

func (u AtomicUnit) Float64() float64 {
	return float64(u) / 1e12
}

func (u AtomicUnit) Formatted() string {
	return strconv.FormatFloat(float64(u)/1e12, 'f', -1, 64)
}

// NewFromFloat converts a float64 to VendoPunkto atomic units
func NewFromFloat(units float64) AtomicUnit {
	return AtomicUnit(units * 1e12)
}
