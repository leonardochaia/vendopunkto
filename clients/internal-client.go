package clients

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/unit"
)

// InternalClient for the internal plugin server hosted by vendopunkto
// Used by plugins to "talk back" to the host.
type InternalClient interface {
	ConfirmPayment(address string, amount unit.AtomicUnit, txHash string, confirmations uint64) error
}

type internalClientImpl struct {
	apiURL url.URL
	client http.Client
}

func NewInternalClient(hostAddress string) (InternalClient, error) {
	apiURL, err := url.Parse(hostAddress)
	if err != nil {
		return nil, err
	}
	return &internalClientImpl{
		apiURL: *apiURL,
		client: http.Client{
			Timeout: 15 * time.Second,
		},
	}, nil
}

// ConfirmPayment should be called when a payment has been confirmed
// on the wallet. Ideally this should be called with 0 confirmations
// when the transaction appears on the mempool, and again when it is confirmed.
func (c internalClientImpl) ConfirmPayment(
	address string,
	amount unit.AtomicUnit,
	txHash string,
	confirmations uint64) error {
	u, err := url.Parse("/v1/invoices/payments/confirm")
	if err != nil {
		return err
	}

	final := c.apiURL.ResolveReference(u)

	params, err := json.Marshal(&dtos.InvoiceConfirmPaymentsParams{
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
