package monero

import (
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
)

type MoneroWalletPlugin struct {
	client wallet.Client
}

func (p MoneroWalletPlugin) GenerateNewAddress(invoiceID string) (string, error) {
	result, err := p.client.MakeIntegratedAddress(&wallet.RequestMakeIntegratedAddress{})

	if err != nil {
		return "", err
	}

	return result.IntegratedAddress, nil
}

func (p MoneroWalletPlugin) GetPluginInfo() (*plugin.WalletPluginInfo, error) {
	return &plugin.WalletPluginInfo{
		Name:       "Monero Wallet",
		ID:         "monero-wallet",
		Currencies: []string{"XMR"},
	}, nil
}
