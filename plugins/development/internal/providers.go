package development

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
)

type Container struct {
	FakeBitcoinWallet       plugin.WalletPlugin
	FakeMoneroWallet        plugin.WalletPlugin
	FakeBitcoinCashWallet   plugin.WalletPlugin
	FakeExchangeRatesPlugin plugin.ExchangeRatesPlugin
	Server                  *plugin.Server
}

var Providers = wire.NewSet(
	plugin.NewServer,
	newFakeExchangeRatesPlugin,
	newContainer)

func newContainer(
	server *plugin.Server,
	exchangeRates plugin.ExchangeRatesPlugin,
) *Container {
	return &Container{
		FakeMoneroWallet:        newFakeMoneroWalletPlugin(),
		FakeBitcoinWallet:       newFakeBitcoinWalletPlugin(),
		FakeBitcoinCashWallet:   newFakeBitcoinCashWalletPlugin(),
		Server:                  server,
		FakeExchangeRatesPlugin: exchangeRates,
	}
}

func newFakeMoneroWalletPlugin() plugin.WalletPlugin {
	return fakeWalletPlugin{
		name:           "Monero",
		symbol:         "xmr",
		qrCodeTemplate: "monero:{{.Address}}?tx_amount={{.Amount}}",
	}
}

func newFakeBitcoinWalletPlugin() plugin.WalletPlugin {
	return fakeWalletPlugin{
		name:   "Bitcoin",
		symbol: "btc",
	}
}

func newFakeBitcoinCashWalletPlugin() plugin.WalletPlugin {
	return fakeWalletPlugin{
		name:   "Bitcoin Cash",
		symbol: "bch",
	}
}

func newFakeExchangeRatesPlugin(logger hclog.Logger) plugin.ExchangeRatesPlugin {
	return fakeExchangeRatesPlugin{
		logger: logger,
	}
}
