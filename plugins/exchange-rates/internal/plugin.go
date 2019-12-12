package rates

import (
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/shopspring/decimal"
	gecko "github.com/superoo7/go-gecko/v3"
)

type geckoExchangeRatesPlugin struct {
	client *gecko.Client
}

func (p geckoExchangeRatesPlugin) GetExchangeRates(
	currency string,
	currencies []string) ([]plugin.ExchangeRatesResult, error) {

	result, err := p.client.SimplePrice(currencies, []string{currency})

	if err != nil {
		return nil, err
	}

	output := []plugin.ExchangeRatesResult{}

	for coin, rates := range *result {
		output = append(output, plugin.ExchangeRatesResult{
			Currency: coin,
			Rate:     decimal.NewFromFloat32(rates[currency]),
		})
	}

	return output, nil
}

func (p geckoExchangeRatesPlugin) GetPluginInfo() (plugin.PluginInfo, error) {
	return plugin.PluginInfo{
		Name: "Gecko Exchange Rates",
		ID:   "gecko-exchange-rates",
		Type: plugin.PluginTypeExchangeRate,
	}, nil
}
