package pluginmgr

import (
	"bytes"
	"encoding/json"
	"github.com/leonardochaia/vendopunkto/plugin"
	"net/http"
	"net/url"
	"time"
)

type coinWalletClientImpl struct {
	apiURL url.URL
	client http.Client
	info   plugin.PluginInfo
}

func NewWalletClient(url url.URL, info plugin.PluginInfo) plugin.WalletPlugin {
	return &coinWalletClientImpl{
		apiURL: url,
		info:   info,
		client: http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c coinWalletClientImpl) GenerateNewAddress(invoiceID string) (string, error) {
	u, err := url.Parse(plugin.WalletMainEndpoint + plugin.GenerateAddressWalletEndpoint)
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
	u, err := url.Parse(plugin.WalletMainEndpoint + plugin.WalletInfoEndpoint)
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
