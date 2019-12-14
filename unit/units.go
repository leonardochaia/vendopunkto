package unit

import (
	"encoding/json"
	"fmt"
)

type AtomicUnit uint64

func (u AtomicUnit) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Value          uint64 `json:"value"`
		ValueFormatted string `json:"valueFormatted"`
	}{
		Value:          uint64(u),
		ValueFormatted: u.String(),
	})
}

func (u AtomicUnit) Float64() float64 {
	return float64(u) / 1e12
}

func (u AtomicUnit) String() string {
	// Credits to https://github.com/monero-ecosystem/go-monero-rpc-client/blob/master/wallet/util.go
	str0 := fmt.Sprintf("%013d", u)
	l := len(str0)
	return str0[:l-12] + "." + str0[l-12:]
}

// NewFromFloat converts a float64 to VendoPunkto atomic units
func NewFromFloat(units float64) AtomicUnit {
	return AtomicUnit(units * 1e12)
}
