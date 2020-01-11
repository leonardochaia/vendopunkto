package pluginmgr

import (
	"net/url"

	"github.com/leonardochaia/vendopunkto/clients"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/plugin"
)

type exchangeRatesClientImpl struct {
	apiURL url.URL
	client clients.HTTP
	info   plugin.PluginInfo
}

func newExchangeRatesClient(
	url url.URL,
	info plugin.PluginInfo,
	client clients.HTTP) plugin.ExchangeRatesPlugin {
	return &exchangeRatesClientImpl{
		apiURL: url,
		info:   info,
		client: client,
	}
}

func (c exchangeRatesClientImpl) getBaseURL(end string) (*url.URL, error) {
	u, err := url.Parse(c.info.GetAddress() + plugin.ExchangeRatesMainEndpoint + end)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c exchangeRatesClientImpl) GetExchangeRates(
	currency string,
	currencies []string) (plugin.ExchangeRatesResult, error) {
	const op errors.Op = "exchangeRatesClient.getExchangeRates"

	u, err := c.getBaseURL("")
	if err != nil {
		return nil, errors.E(op, errors.Internal, err)
	}

	final := c.apiURL.ResolveReference(u).String()

	params := &plugin.GetExchangeRatesParams{
		Currency:   currency,
		Currencies: currencies,
	}

	var result plugin.ExchangeRatesResult
	_, err = c.client.PostJSON(final, params, &result)

	if err != nil {
		return nil, errors.E(op, err)
	}

	return result, nil
}

func (c exchangeRatesClientImpl) GetPluginInfo() (plugin.PluginInfo, error) {
	return c.info, nil
}
