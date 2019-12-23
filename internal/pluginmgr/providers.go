package pluginmgr

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/clients"
	"github.com/leonardochaia/vendopunkto/internal/conf"
)

var PluginProviders = wire.NewSet(NewManager)

func NewManager(
	logger hclog.Logger,
	client clients.HTTP,
	startupConf conf.Startup) *Manager {
	return &Manager{
		logger:        logger.Named("pluginmgr"),
		wallets:       make(map[string]walletAndInfo),
		exchangeRates: make(map[string]exchangeRatesAndInfo),
		client:        client,
		startupConf:   startupConf,
	}
}
