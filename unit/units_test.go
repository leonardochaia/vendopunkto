package unit

import (
	"testing"

	"github.com/leonardochaia/vendopunkto/testutils"
)

func TestUnitCreationtFromFloat64(t *testing.T) {
	expected := AtomicUnit(1500000000000)
	result := NewFromFloat(1.5)
	testutils.Equals(t, expected, result)

	expected = AtomicUnit(1255230000)
	result = NewFromFloat(0.00125523)
	testutils.Equals(t, expected, result)
}

func TestUnitConversionToFloat64(t *testing.T) {
	expected := 1.5
	result := AtomicUnit(1500000000000).Float64()
	testutils.Equals(t, expected, result)

	expected = 0.00125523
	result = AtomicUnit(1255230000).Float64()
	testutils.Equals(t, expected, result)
}

func TestUnitConversionToString(t *testing.T) {
	expected := "1.5"
	result := AtomicUnit(1500000000000).Formatted()
	testutils.Equals(t, expected, result)

	expected = "0.00125523"
	result = AtomicUnit(1255230000).Formatted()
	testutils.Equals(t, expected, result)
}
