package rates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/shopspring/decimal"
	gecko "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
)

type geckoExchangeRatesPlugin struct {
	client *gecko.Client
	coins  types.CoinList
	logger hclog.Logger
}

func (p *geckoExchangeRatesPlugin) loadAllCoins() error {
	if p.coins == nil {
		allCoins, err := p.client.CoinsList()
		if err != nil {
			return err
		}

		p.coins = *allCoins
		p.logger.Info("Obtained coin list", "amount", len(p.coins))
	}
	return nil
}

func (p *geckoExchangeRatesPlugin) getGeckoCurrencyIDForSymbol(symbol string) (string, error) {
	symbol = strings.ToLower(symbol)

	if err := p.loadAllCoins(); err != nil {
		return "", err
	}

	for _, coin := range p.coins {
		if coin.Symbol == symbol {
			return coin.ID, nil
		}
	}

	return "", fmt.Errorf("Couldn't find a coin for symbol " + symbol)
}

func (p geckoExchangeRatesPlugin) SearchSupportedCurrencies(term string) ([]dtos.BasicCurrencyDto, error) {

	currencies, err := p.client.SimpleSupportedVSCurrencies()
	if err != nil {
		return nil, err
	}

	result := make([]dtos.BasicCurrencyDto, len(*currencies))

	i := 0
	for _, c := range *currencies {
		result[i] = dtos.BasicCurrencyDto{
			Symbol: c,
			Name:   strings.ToUpper(c),
		}
		i++
	}

	return result, nil
}

func (p geckoExchangeRatesPlugin) GetExchangeRates(
	currency string,
	currencies []string) (plugin.ExchangeRatesResult, error) {

	sources := make([]string, len(currencies))
	sourceMap := make(map[string]string)
	i := 0
	for _, c := range currencies {
		source, err := p.getGeckoCurrencyIDForSymbol(c)
		if err != nil {
			return nil, err
		}
		sourceMap[source] = c
		sources[i] = source
		i++
	}

	p.logger.Info("Requesting Gecko Simple Price",
		"source", currency,
		"currencies", sources)
	result, err := p.client.SimplePrice(sources, []string{currency})

	if err != nil {
		return nil, err
	}

	output := make(plugin.ExchangeRatesResult)

	for _, target := range sources {
		rate := decimal.NewFromFloat32((*result)[target][currency])
		if target == currency {
			rate = decimal.NewFromInt(1)
		}

		p.logger.Info("Obtained exchange rate",
			"source", currency,
			"currency", sourceMap[target],
			"rate", rate)
		output[sourceMap[target]] = decimal.NewFromInt(1).Div(rate)
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
