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
	CreateInvoice(total unit.AtomicUnit, currency string) (*dtos.InvoiceDto, error)
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
	vURL, err := url.Parse("/api/v1/")
	if err != nil {
		return nil, err
	}

	return &internalClientImpl{
		apiURL: *apiURL.ResolveReference(vURL),
		client: client,
	}, nil
}

func (c internalClientImpl) getAPIURL(suffix string) (string, error) {
	u, err := url.Parse(suffix)
	if err != nil {
		return "", err
	}

	final := c.apiURL.ResolveReference(u)
	return final.String(), nil
}

// CreateInvoice will create a new invoice with the provided total and currency
// all payment methods will be added
func (c internalClientImpl) CreateInvoice(
	total unit.AtomicUnit,
	currency string) (*dtos.InvoiceDto, error) {
	const op errors.Op = "internalClient.createInvoice"

	url, err := c.getAPIURL("invoices")
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	params := dtos.InvoiceCreationParams{
		Total:    total,
		Currency: currency,
	}

	var result *dtos.InvoiceDto
	_, err = c.client.PostJSON(url, params, &result)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return result, nil
}
