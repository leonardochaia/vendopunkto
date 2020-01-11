package vendopunkto

import (
	"context"

	"github.com/leonardochaia/vendopunkto/plugin"
)

// PluginManager is the first line of control for plugins
type PluginManager interface {
	LoadPlugins(ctx context.Context)
	GetWallet(ID string) (plugin.WalletPlugin, error)
	GetWalletForCurrency(currency string) (plugin.WalletPlugin, error)
	GetWalletInfoForCurrency(currency string) (plugin.WalletPluginInfo, error)
	GetAllCurrencies() ([]plugin.WalletPluginCurrency, error)
	GetExchangeRatesPlugin(ID string) (plugin.ExchangeRatesPlugin, error)
	GetConfiguredExchangeRatesPlugin() (plugin.ExchangeRatesPlugin, error)
}
