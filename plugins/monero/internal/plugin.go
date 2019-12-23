package monero

import (
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
)

type moneroWalletPlugin struct {
	client wallet.Client
	logger hclog.Logger
}

func (p moneroWalletPlugin) GenerateNewAddress(invoiceID string) (string, error) {
	result, err := p.client.MakeIntegratedAddress(&wallet.RequestMakeIntegratedAddress{})

	if err != nil {
		return "", err
	}

	p.logger.Info("Generated new address", "address", result.IntegratedAddress)
	return result.IntegratedAddress, nil
}

func (p moneroWalletPlugin) GetPluginInfo() (plugin.PluginInfo, error) {
	return plugin.PluginInfo{
		Name: "Monero Wallet",
		ID:   "monero-wallet",
		Type: plugin.PluginTypeWallet,
	}, nil
}

func (p moneroWalletPlugin) GetWalletInfo() (plugin.WalletPluginInfo, error) {
	return plugin.WalletPluginInfo{
		Currency: plugin.WalletPluginCurrency{
			Name:           "Monero",
			Symbol:         "XMR",
			QRCodeTemplate: "monero:{{.Address}}?tx_amount={{.Amount}}",
		},
	}, nil
}
