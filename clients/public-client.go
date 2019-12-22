package clients

import (
	"net/url"

	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/unit"
)

// PublicClient for the internal plugin server hosted by vendopunkto
// Used by plugins to "talk back" to the host.
type PublicClient interface {
	CreateInvoice(total unit.AtomicUnit, currency string) (*dtos.InvoiceDto, error)
	GetInvoice(ID string) (*dtos.InvoiceDto, error)
	GeneratePaymentMethodAdress(invoiceID string, currency string) (*dtos.InvoiceDto, error)
}

type publicClientImpl struct {
	apiURL url.URL
	client HTTP
}

// NewPublicClient creates a new PublicClient
func NewPublicClient(hostAddress string, client HTTP) (PublicClient, error) {
	apiURL, err := url.Parse(hostAddress)
	if err != nil {
		return nil, err
	}
	return &publicClientImpl{
		apiURL: *apiURL,
		client: client,
	}, nil
}

func (c publicClientImpl) getAPIURL(suffix string) (string, error) {
	u, err := url.Parse(suffix)
	if err != nil {
		return "", err
	}

	final := c.apiURL.ResolveReference(u)
	return final.String(), nil
}

func (c publicClientImpl) CreateInvoice(
	total unit.AtomicUnit,
	currency string) (*dtos.InvoiceDto, error) {
	const op errors.Op = "publicClient.createInvoice"

	url, err := c.getAPIURL("/v1/invoices")
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	params := dtos.InvoiceCreationParams{
		Total:    total,
		Currency: currency,
		// PaymentMethods: ,
	}

	var result *dtos.InvoiceDto
	_, err = c.client.PostJSON(url, params, &result)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return result, nil
}

func (c publicClientImpl) GetInvoice(ID string) (*dtos.InvoiceDto, error) {
	const op errors.Op = "publicClient.getInvoice"

	url, err := c.getAPIURL("/v1/invoices/" + ID)
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	var result *dtos.InvoiceDto
	_, err = c.client.GetJSON(url, &result)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return result, nil

}

func (c publicClientImpl) GeneratePaymentMethodAdress(invoiceID string, currency string) (*dtos.InvoiceDto, error) {
	const op errors.Op = "publicClient.getInvoice"

	url, err := c.getAPIURL("/v1/invoices/" + invoiceID + "/payment-method/address")
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	params := dtos.InvoiceGeneratePaymentMethodAddressParams{
		Currency: currency,
	}

	var result *dtos.InvoiceDto
	_, err = c.client.PostJSON(url, params, &result)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return result, nil
}
