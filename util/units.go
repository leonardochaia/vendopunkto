package util

import "fmt"

// Amounts are stored as integeres same as Monero does
// Credits to https://github.com/monero-ecosystem/go-monero-rpc-client/blob/master/wallet/util.go

// VendoPunktoToDecimal converts a raw atomic VendoPunkto balance to a more
// human readable format.
func VendoPunktoToDecimal(units uint64) string {
	str0 := fmt.Sprintf("%013d", units)
	l := len(str0)
	return str0[:l-12] + "." + str0[l-12:]
}

// VendoPunktoToFloat64 converts raw atomic VendoPunkto to a float64
func VendoPunktoToFloat64(units uint64) float64 {
	return float64(units) / 1e12
}

// Float64ToVendoPunkto converts raw atomic VendoPunkto to a float64
func Float64ToVendoPunkto(units float64) uint64 {
	return uint64(units * 1e12)
}
