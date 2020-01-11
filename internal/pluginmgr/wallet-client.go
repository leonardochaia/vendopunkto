package pluginmgr

import (
	"net/url"

	"github.com/leonardochaia/vendopunkto/clients"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/plugin"
)

type coinWalletClientImpl struct {
	apiURL url.URL
	client clients.HTTP
	info   plugin.PluginInfo
}

// newWalletClient constructs a new client implementing WalletPlugin
func newWalletClient(
	url url.URL,
	info plugin.PluginInfo,
	client clients.HTTP) plugin.WalletPlugin {
	return &coinWalletClientImpl{
		apiURL: url,
		info:   info,
		client: client,
	}
}

func (c coinWalletClientImpl) getBaseURL(end string) (*url.URL, error) {
	u, err := url.Parse(c.info.GetAddress() + plugin.WalletMainEndpoint + end)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c coinWalletClientImpl) GenerateNewAddress(invoiceID string) (string, error) {
	const op errors.Op = "wallet.generateNewAddress"
	u, err := c.getBaseURL(plugin.GenerateAddressWalletEndpoint)
	if err != nil {
		return "", errors.E(op, errors.Internal, err)
	}

	final := c.apiURL.ResolveReference(u).String()

	params := &plugin.CoinWalletAddressParams{
		InvoiceID: invoiceID,
	}

	var result plugin.CoinWalletAddressResponse
	_, err = c.client.PostJSON(final, params, &result)

	if err != nil {
		return "", errors.E(op, err)
	}

	return result.Address, nil
}

func (c coinWalletClientImpl) GetIncomingTransfers(params plugin.WalletPluginIncomingTransferParams) ([]plugin.WalletPluginIncomingTransferResult, error) {
	const op errors.Op = "wallet.getIncomingTransfers"
	u, err := c.getBaseURL(plugin.GetIncomingTransfersWalletEndpoint)
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	final := c.apiURL.ResolveReference(u).String()

	var result []plugin.WalletPluginIncomingTransferResult
	_, err = c.client.PostJSON(final, params, &result)

	if err != nil {
		return nil, errors.E(op, err)
	}

	return result, nil
}

func (c coinWalletClientImpl) GetPluginInfo() (plugin.PluginInfo, error) {
	return c.info, nil
}

func (c coinWalletClientImpl) GetWalletInfo() (plugin.WalletPluginInfo, error) {
	const op errors.Op = "wallet.getWalletInfo"
	u, err := c.getBaseURL(plugin.WalletInfoEndpoint)
	if err != nil {
		return plugin.WalletPluginInfo{}, errors.E(op, errors.Internal, err)
	}

	final := c.apiURL.ResolveReference(u)

	var result plugin.WalletPluginInfo
	_, err = c.client.GetJSON(final.String(), &result)

	if err != nil {
		return result, errors.E(op, err)
	}

	return result, nil
}
