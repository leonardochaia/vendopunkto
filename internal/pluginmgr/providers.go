package pluginmgr

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/clients"
)

var PluginProviders = wire.NewSet(NewManager)

func NewManager(
	logger hclog.Logger,
	client clients.HTTP) *Manager {
	return &Manager{
		logger:        logger.Named("pluginmgr"),
		wallets:       make(map[string]walletAndInfo),
		exchangeRates: make(map[string]exchangeRatesAndInfo),
		client:        client,
	}
}
