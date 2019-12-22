package clients

import (
	"net/url"

	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/unit"
)

// InternalClient for the internal plugin server hosted by vendopunkto
// Used by plugins to "talk back" to the host.
type InternalClient interface {
	ConfirmPayment(address string, amount unit.AtomicUnit, txHash string, confirmations uint64) error
}

type internalClientImpl struct {
	apiURL url.URL
	client HTTP
}

// NewInternalClient creates an InternalClient
func NewInternalClient(hostAddress string, client HTTP) (InternalClient, error) {
	apiURL, err := url.Parse(hostAddress)
	if err != nil {
		return nil, err
	}
	return &internalClientImpl{
		apiURL: *apiURL,
		client: client,
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

	const op errors.Op = "internalClient.confirmPayment"

	u, err := url.Parse("/v1/invoices/payments/confirm")
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	final := c.apiURL.ResolveReference(u).String()

	params := dtos.InvoiceConfirmPaymentsParams{
		Address:       address,
		Amount:        amount,
		TxHash:        txHash,
		Confirmations: confirmations,
	}

	_, err = c.client.PostJSON(final, params, nil)

	if err != nil {
		return errors.E(op, err)
	}

	return nil
}
