package development

import (
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/shopspring/decimal"
)

// rates is a map of crypto to USD
var rates = make(plugin.ExchangeRatesResult)

func init() {
	rates["usd"] = decimal.NewFromInt(1)
	rates["xmr"] = decimal.NewFromInt(50)
	rates["btc"] = decimal.NewFromInt(10000)
	rates["bch"] = decimal.NewFromInt(5000)
	rates["ars"] = decimal.NewFromFloat(0.2)
}

type fakeExchangeRatesPlugin struct {
	logger hclog.Logger
}

func (p fakeExchangeRatesPlugin) SearchSupportedCurrencies(term string) ([]dtos.BasicCurrencyDto, error) {
	result := make([]dtos.BasicCurrencyDto, len(rates))

	i := 0
	for k := range rates {
		result[i] = dtos.BasicCurrencyDto{
			Symbol: k,
		}
		i++
	}
	return result, nil
}

func (p fakeExchangeRatesPlugin) GetExchangeRates(
	currency string,
	currencies []string) (plugin.ExchangeRatesResult, error) {

	output := make(plugin.ExchangeRatesResult)

	currency = strings.ToLower(currency)

	converted := rates[currency]

	for coin, rate := range rates {
		newRate := converted.Div(rate)
		p.logger.Info("Obtained exchange rate",
			"source", currency,
			"currency", coin,
			"rate", newRate.String())

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
