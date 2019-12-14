package pluginmgr

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/leonardochaia/vendopunkto/plugin"
)

type exchangeRatesClientImpl struct {
	apiURL url.URL
	client http.Client
	info   plugin.PluginInfo
}

func NewExchangeRatesClient(
	url url.URL,
	info plugin.PluginInfo,
	client http.Client) plugin.ExchangeRatesPlugin {
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

	u, err := c.getBaseURL("")
	if err != nil {
		return nil, err
	}

	final := c.apiURL.ResolveReference(u)

	params, err := json.Marshal(&plugin.GetExchangeRatesParams{
		Currency:   currency,
		Currencies: currencies,
	})

	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(final.String(), "application/json", bytes.NewBuffer(params))

	if err != nil {
		return nil, err
	}

	var result plugin.ExchangeRatesResult
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}

func (c exchangeRatesClientImpl) GetPluginInfo() (plugin.PluginInfo, error) {
	return c.info, nil
}
