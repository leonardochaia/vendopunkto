package monero

import (
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/shopspring/decimal"
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

func (p moneroWalletPlugin) GetIncomingTransfers(params plugin.WalletPluginIncomingTransferParams) (
	[]plugin.WalletPluginIncomingTransferResult, error) {

	resp, err := p.client.GetTransfers(&wallet.RequestGetTransfers{
		In:             true,
		Pool:           true,
		FilterByHeight: params.MinBlockHeight > 0,
		MinHeight:      params.MinBlockHeight,
		// TODO: this is a hack until https://github.com/monero-ecosystem/go-monero-rpc-client/issues/4
		MaxHeight: 9223372036854775807,
	})

	if err != nil {
		return nil, err
	}

	p.logger.Info("Obtained TXs",
		"pool", len(resp.Pool),
		"in", len(resp.In),
		"minHeight", params.MinBlockHeight)

	output := []plugin.WalletPluginIncomingTransferResult{}
	for _, transfer := range append(resp.In, resp.Pool...) {

		addrResp, err := p.client.MakeIntegratedAddress(&wallet.RequestMakeIntegratedAddress{
			StandardAddress: transfer.Address,
			PaymentID:       transfer.PaymentID,
		})

		if err != nil {
			return nil, err
		}

		tx := plugin.WalletPluginIncomingTransferResult{
			TxHash:        transfer.TxID,
			Address:       addrResp.IntegratedAddress,
			BlockHeight:   transfer.Height,
			Confirmations: transfer.Confirmations,
			Amount:        unit.AtomicUnit(transfer.Amount),
		}
		output = append(output, tx)
	}

	return output, nil
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
