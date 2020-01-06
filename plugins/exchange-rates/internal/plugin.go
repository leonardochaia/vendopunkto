package rates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"
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

func (p *geckoExchangeRatesPlugin) getGeckoCurrencyIDForSymbol(symbol string) (string, error) {
	symbol = strings.ToLower(symbol)

	if p.coins == nil {
		allCoins, err := p.client.CoinsList()
		if err != nil {
			return "", err
		}

		p.coins = *allCoins
		p.logger.Info("Obtained coin list", "amount", len(p.coins))
	}

	for _, coin := range p.coins {
		if coin.Symbol == symbol {
			return coin.ID, nil
		}
	}

	return "", fmt.Errorf("Couldn't find a coin for symbol " + symbol)
}

func (p geckoExchangeRatesPlugin) GetExchangeRates(
	currency string,
	currencies []string) (plugin.ExchangeRatesResult, error) {

	source, err := p.getGeckoCurrencyIDForSymbol(currency)
	if err != nil {
		return nil, err
	}

	p.logger.Info("Requesting Gecko Simple Price",
		"source", source,
		"currencies", currencies)
	result, err := p.client.SimplePrice([]string{source}, currencies)

	if err != nil {
		return nil, err
	}

	output := make(plugin.ExchangeRatesResult)

	for _, rates := range *result {

		for _, target := range currencies {
			rate := decimal.NewFromFloat32(rates[target])
			if target == currency {
				rate = decimal.NewFromInt(1)
			}

			p.logger.Info("Obtained exchange rate",
				"source", source,
				"currency", target,
				"rate", rate)
			output[target] = rate
		}
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
