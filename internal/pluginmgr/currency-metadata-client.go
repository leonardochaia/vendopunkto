package pluginmgr

import (
	"net/url"

	"github.com/leonardochaia/vendopunkto/clients"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/plugin"
)

type currencyMetadataClientImpl struct {
	apiURL url.URL
	client clients.HTTP
	info   plugin.PluginInfo
}

func newCurrencyMetadataClient(
	url url.URL,
	info plugin.PluginInfo,
	client clients.HTTP) plugin.CurrencyMetadataPlugin {
	return &currencyMetadataClientImpl{
		apiURL: url,
		info:   info,
		client: client,
	}
}

func (c currencyMetadataClientImpl) getBaseURL(end string) (*url.URL, error) {
	u, err := url.Parse(c.info.GetAddress() + plugin.CurrencyMetadataMainEndpoint + end)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c currencyMetadataClientImpl) GetCurrencies(
	currencies []string) ([]plugin.CurrencyMetadata, error) {
	const op errors.Op = "currencyMetadataClient.GetCurrencies"

	u, err := c.getBaseURL("")
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	final := c.apiURL.ResolveReference(u).String()

	params := &plugin.GetCurrencyParams{
		Currencies: currencies,
	}

	var result []plugin.CurrencyMetadata
	_, err = c.client.PostJSON(final, params, &result)

	if err != nil {
		return nil, errors.E(op, err)
	}

	return result, nil
}

func (c currencyMetadataClientImpl) GetPluginInfo() (plugin.PluginInfo, error) {
	return c.info, nil
}
