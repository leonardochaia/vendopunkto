package development

import (
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/rs/xid"
)

type fakeWalletPlugin struct {
	symbol string
	name   string
}

func (p fakeWalletPlugin) GenerateNewAddress(invoiceID string) (string, error) {
	return p.symbol + "-fake-" + xid.New().String(), nil
}

func (p fakeWalletPlugin) GetWalletInfo() (plugin.WalletPluginInfo, error) {
	return plugin.WalletPluginInfo{
		Currency: plugin.WalletPluginCurrency{
			Name:   p.name,
			Symbol: p.symbol,
		},
	}, nil
}

func (p fakeWalletPlugin) GetPluginInfo() (plugin.PluginInfo, error) {
	return plugin.PluginInfo{
		Name: p.name + " Fake Wallet",
		ID:   p.symbol + "-fake-wallet",
		Type: plugin.PluginTypeWallet,
	}, nil
}
