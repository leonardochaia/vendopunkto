package development

import (
	"github.com/google/wire"
	"github.com/leonardochaia/vendopunkto/plugin"
)

type Container struct {
	FakeBitcoinWallet     plugin.WalletPlugin
	FakeMoneroWallet      plugin.WalletPlugin
	FakeBitcoinCashWallet plugin.WalletPlugin
	Server                *plugin.Server
}

var Providers = wire.NewSet(
	plugin.NewServer,
	newContainer)

func newContainer(
	server *plugin.Server,
) *Container {
	return &Container{
		FakeMoneroWallet:      newFakeMoneroWalletPlugin(),
		FakeBitcoinWallet:     newFakeBitcoinWalletPlugin(),
		FakeBitcoinCashWallet: newFakeBitcoinCashWalletPlugin(),
		Server:                server,
	}
}

func newFakeMoneroWalletPlugin() plugin.WalletPlugin {
	return fakeWalletPlugin{
		name:   "Monero",
		symbol: "xmr",
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
