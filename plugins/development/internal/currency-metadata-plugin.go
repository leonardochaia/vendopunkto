package development

import (
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
)

var defaultCurrencies = []plugin.CurrencyMetadata{
	plugin.CurrencyMetadata{
		Name:         "Monero",
		Symbol:       "xmr",
		LogoImageURL: "https://assets.coingecko.com/coins/images/69/large/monero_logo.png",
	},
	plugin.CurrencyMetadata{
		Name:         "Bitcoin",
		Symbol:       "btc",
		LogoImageURL: "https://assets.coingecko.com/coins/images/1/large/bitcoin.png?1547033579",
	},
	plugin.CurrencyMetadata{
		Name:         "Bitcoin Cash",
		Symbol:       "bch",
		LogoImageURL: "https://assets.coingecko.com/coins/images/780/large/bitcoin_cash.png?1547034541",
	},
	plugin.CurrencyMetadata{
		Name:   "US Dollars",
		Symbol: "usd",
	},
	plugin.CurrencyMetadata{
		Name:   "Argentine Peso",
		Symbol: "ars",
	},
}

type fakeCurrencyMetadata struct {
	logger hclog.Logger
}

func (p fakeCurrencyMetadata) GetCurrencies(
	currencies []string) ([]plugin.CurrencyMetadata, error) {

	return defaultCurrencies, nil
}

func (p fakeCurrencyMetadata) GetPluginInfo() (plugin.PluginInfo, error) {
	return plugin.PluginInfo{
		Name: "Fake Currency Metadata",
		ID:   "fake-currency-metadata",
		Type: plugin.PluginTypeCurrencyMetadata,
	}, nil
}
