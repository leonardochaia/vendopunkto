package plugin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type coinWalletClientImpl struct {
	apiURL url.URL
	client http.Client
}

func NewWalletClient(url url.URL) WalletPlugin {
	return &coinWalletClientImpl{
		apiURL: url,
		client: http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c coinWalletClientImpl) GenerateNewAddress(invoiceID string) (string, error) {
	u, err := url.Parse(WalletMainEndpoint + GenerateAddressWalletEndpoint)
	if err != nil {
		return "", err
	}

	final := c.apiURL.ResolveReference(u)

	params, err := json.Marshal(&CoinWalletAddressParams{
		InvoiceID: invoiceID,
	})

	if err != nil {
		return "", err
	}

	resp, err := c.client.Post(final.String(), "application/json", bytes.NewBuffer(params))

	if err != nil {
		return "", err
	}

	var result CoinWalletAddressResponse
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result.Address, err
}

func (c coinWalletClientImpl) GetPluginInfo() (*PluginInfo, error) {
	u, err := url.Parse(WalletMainEndpoint + PluginInfoEndpoint)
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

	var result PluginInfo
	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}
