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
}

func NewWalletClient(url url.URL) plugin.WalletPlugin {
	return &coinWalletClientImpl{
		apiURL: url,
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

func (c coinWalletClientImpl) GetPluginInfo() (*plugin.WalletPluginInfo, error) {
	u, err := url.Parse(plugin.WalletMainEndpoint + plugin.PluginInfoEndpoint)
	if err != nil {
		return nil, err
	}

	final := c.apiURL.ResolveReference(u)

	if err != nil {
		return nil, err
	}

	resp, err := c.client.Get(final.String())

	if err != nil {
		return nil, err
	}

	var result plugin.WalletPluginInfo
	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}
