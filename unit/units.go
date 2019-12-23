// Package unit provides utilites for working with unit.AtomicUnits
package unit

import (
	"strconv"
)

// AtomicUnit is the representation used to store cryptocurrency amounts
// in VendoPunto. One atomic unit is the smallest fraction that VendoPunkto
// can store. It's underlying implementation is a uint64.
// 1.0 <CryptoAmount> equals 1e12 AtomicUnit
// This is based on Monero
type AtomicUnit uint64

// Float64 returns a float64 representation of the value.
func (u AtomicUnit) Float64() float64 {
	return float64(u) / 1e12
}

// Formatted returns a string representation of the value.
func (u AtomicUnit) Formatted() string {
	return strconv.FormatFloat(float64(u)/1e12, 'f', -1, 64)
}

// NewFromFloat converts a float64 to VendoPunkto atomic units
func NewFromFloat(units float64) AtomicUnit {
	return AtomicUnit(units * 1e12)
}
