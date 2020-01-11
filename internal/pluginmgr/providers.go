package pluginmgr

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/clients"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
	"github.com/leonardochaia/vendopunkto/internal/conf"
)

// Providers implementations
var Providers = wire.NewSet(NewPluginManager)

// NewPluginManager creates the PluginManager implementation
func NewPluginManager(
	logger hclog.Logger,
	client clients.HTTP,
	currencyRepo vendopunkto.CurrencyRepository,
	startupConf conf.Startup) vendopunkto.PluginManager {
	return &pluginManager{
		logger:        logger.Named("plugin-manager"),
		wallets:       make(map[string]walletAndInfo),
		exchangeRates: make(map[string]exchangeRatesAndInfo),
		client:        client,
		currencyRepo:  currencyRepo,
		startupConf:   startupConf,
	}
}
