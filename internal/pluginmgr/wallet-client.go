package pluginmgr

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/leonardochaia/vendopunkto/plugin"
)

type coinWalletClientImpl struct {
	apiURL url.URL
	client http.Client
	info   plugin.PluginInfo
}

func NewWalletClient(
	url url.URL,
	info plugin.PluginInfo,
	client http.Client) plugin.WalletPlugin {
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
	u, err := c.getBaseURL(plugin.GenerateAddressWalletEndpoint)
	if err != nil {
		return "", err
	}

	final := c.apiURL.ResolveReference(u)

	params, err := json.Marshal(&plugin.CoinWalletAddressParams{
		InvoiceID: invoiceID,
	})

	if err != nil {
		return "", err
	}

	resp, err := c.client.Post(final.String(), "application/json", bytes.NewBuffer(params))

	if err != nil {
		return "", err
	}

	var result plugin.CoinWalletAddressResponse
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result.Address, err
}

func (c coinWalletClientImpl) GetPluginInfo() (plugin.PluginInfo, error) {
	return c.info, nil
}

func (c coinWalletClientImpl) GetWalletInfo() (plugin.WalletPluginInfo, error) {
	u, err := c.getBaseURL(plugin.WalletInfoEndpoint)
	if err != nil {
		return plugin.WalletPluginInfo{}, err
	}

	final := c.apiURL.ResolveReference(u)

	resp, err := c.client.Get(final.String())

	if err != nil {
		return plugin.WalletPluginInfo{}, err
	}

	var result plugin.WalletPluginInfo
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
