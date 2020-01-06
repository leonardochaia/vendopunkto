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
	expected := fmt.Sprintf("monero:%s?tx_amount=%d", address, amount)

	wi := WalletPluginInfo{
		Currency: WalletPluginCurrency{
			Name:           "Monero",
			Symbol:         "XMR",
			QRCodeTemplate: "monero:{{.Address}}?tx_amount={{.Amount}}",
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
