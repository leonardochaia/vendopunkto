package clients

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/unit"
)

// VendoPunktoInternalClient for the internal plugin server hosted by vendopunkto
// Used by plugins to "talk back" to the host.
type PublicClient interface {
	CreateInvoice(total unit.AtomicUnit, currency string) (*dtos.InvoiceDto, error)
	GetInvoice(ID string) (*dtos.InvoiceDto, error)
	GeneratePaymentMethodAdress(invoiceID string, currency string) (*dtos.InvoiceDto, error)
}

type publicClientImpl struct {
	apiURL url.URL
	client http.Client
}

func NewPublicClient(hostAddress string) (PublicClient, error) {
	apiURL, err := url.Parse(hostAddress)
	if err != nil {
		return nil, err
	}
	return &publicClientImpl{
		apiURL: *apiURL,
		client: http.Client{
			Timeout: 15 * time.Second,
		},
	}, nil
}

func checkAPIResponse(resp *http.Response) error {
	if resp.StatusCode >= 400 {
		return errors.DecodeError(resp)
	}

	return nil
}

func decodeInvoice(resp *http.Response) (*dtos.InvoiceDto, error) {
	inv := &dtos.InvoiceDto{}
	err := json.NewDecoder(resp.Body).Decode(&inv)
	if err != nil {
		return nil, err
	}
	return inv, nil
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

	url, err := c.getAPIURL("/v1/invoices")
	if err != nil {
		return nil, err
	}

	params, err := json.Marshal(&dtos.InvoiceCreationParams{
		Total:    total,
		Currency: currency,
		// PaymentMethods: ,
	})

	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err = checkAPIResponse(resp); err != nil {
		return nil, err
	}

	return decodeInvoice(resp)
}

func (c publicClientImpl) GetInvoice(ID string) (*dtos.InvoiceDto, error) {

	url, err := c.getAPIURL("/v1/invoices/" + ID)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err = checkAPIResponse(resp); err != nil {
		return nil, err
	}

	return decodeInvoice(resp)
}

func (c publicClientImpl) GeneratePaymentMethodAdress(invoiceID string, currency string) (*dtos.InvoiceDto, error) {
	url, err := c.getAPIURL("/v1/invoices/" + invoiceID + "/payment-method/address")
	if err != nil {
		return nil, err
	}

	params, err := json.Marshal(&dtos.InvoiceGeneratePaymentMethodAddressParams{
		Currency: currency,
	})

	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err = checkAPIResponse(resp); err != nil {
		return nil, err
	}

	return decodeInvoice(resp)
}
