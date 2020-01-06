package plugin

import (
	"fmt"
	"testing"

	"github.com/leonardochaia/vendopunkto/testutils"
	"github.com/shopspring/decimal"
)

func TestQRGeneration(t *testing.T) {

	address := "xmr-fake-addr"
	amount := decimal.NewFromFloat(1)
	expected := fmt.Sprintf("monero:%s?tx_amount=%s", address, "1000000000000")

	wi := WalletPluginInfo{
		Currency: WalletPluginCurrency{
			Name:           "Monero",
			Symbol:         "XMR",
			QRCodeTemplate: "{{$t:= newDecimal 1000000000000}}monero:{{.Address}}?tx_amount={{.Amount.Mul $t}}",
		},
	}

	result, err := wi.BuildQRCode(address, amount)
	testutils.Ok(t, err)

	testutils.Equals(t, expected, result)
}

func TestQRBIP21Generation(t *testing.T) {

	address := "btc-fake-addr"
	amount := decimal.NewFromFloat(50)
	expected := fmt.Sprintf("bitcoin:%s?amount=%s", address, amount.String())

	wi := WalletPluginInfo{
		Currency: WalletPluginCurrency{
			Name:   "Bitcoin",
			Symbol: "BTC",
		},
	}

	result, err := wi.BuildQRCode(address, amount)
	testutils.Ok(t, err)

	testutils.Equals(t, expected, result)
}
