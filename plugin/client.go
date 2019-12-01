package plugin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// PluginWalletClient for the internal plugin server hosted by vendopunkto
// Used by plugins to "talk back" to the host.
type PluginWalletClient interface {
	ConfirmPayment(address string, amount uint64, txHash string, confirmations uint) error
}

type pluginWalletClientImpl struct {
	apiURL url.URL
	client http.Client
}

func NewWalletClient(url url.URL) PluginWalletClient {
	return &pluginWalletClientImpl{
		apiURL: url,
		client: http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ConfirmPayment should be called when a payment has been confirmed
// on the wallet. Ideally this should be called with 0 confirmations
// when the transaction appears on the mempool, and again when it is confirmed.
func (c pluginWalletClientImpl) ConfirmPayment(
	address string,
	amount uint64,
	txHash string,
	confirmations uint) error {
	u, err := url.Parse("/v1/invoices/payments/confirm")
	if err != nil {
		return err
	}

	final := c.apiURL.ResolveReference(u)

	type confirmPaymentsParams struct {
		TxHash        string `json:"txHash"`
		Amount        uint64 `json:"amount"`
		Address       string `json:"address"`
		Confirmations uint   `json:"confirmations"`
	}

	params, err := json.Marshal(&confirmPaymentsParams{
		Address:       address,
		Amount:        amount,
		TxHash:        txHash,
		Confirmations: confirmations,
	})

	if err != nil {
		return err
	}

	resp, err := c.client.Post(final.String(), "application/json", bytes.NewBuffer(params))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return err
}
