package monero

import (
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
)

type moneroWalletPlugin struct {
	client wallet.Client
}

func (p moneroWalletPlugin) GenerateNewAddress(invoiceID string) (string, error) {
	result, err := p.client.MakeIntegratedAddress(&wallet.RequestMakeIntegratedAddress{})

	if err != nil {
		return "", err
	}

	return result.IntegratedAddress, nil
}

func (p moneroWalletPlugin) GetPluginInfo() (*plugin.WalletPluginInfo, error) {
	return &plugin.WalletPluginInfo{
		Name:       "Monero Wallet",
		ID:         "monero-wallet",
		Currencies: []string{"XMR"},
	}, nil
}
