package development

import (
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
)

// rates is a map of crypto to USD
var rates = make(plugin.ExchangeRatesResult)

func init() {
	rates["xmr"] = 50
	rates["btc"] = 10000
	rates["bch"] = 5000
	rates["usd"] = 1
}

type fakeExchangeRatesPlugin struct {
	logger hclog.Logger
}

func (p fakeExchangeRatesPlugin) GetExchangeRates(
	currency string,
	currencies []string) (plugin.ExchangeRatesResult, error) {

	output := make(plugin.ExchangeRatesResult)

	if currency == "usd" {
		return rates, nil
	}

	converted := rates[currency]

	for coin, rate := range rates {
		newRate := converted / rate
		p.logger.Info("Obtained exchange rate",
			"source", currency,
			"currency", coin,
			"rate", newRate)

		output[coin] = newRate
	}

	return output, nil
}

func (p fakeExchangeRatesPlugin) GetPluginInfo() (plugin.PluginInfo, error) {
	return plugin.PluginInfo{
		Name: "Fake Exchange Rates",
		ID:   "fake-exchange-rates",
		Type: plugin.PluginTypeExchangeRate,
	}, nil
}
