package pluginmgr

import (
	"net/http"

	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
)

var PluginProviders = wire.NewSet(NewManager)

func NewManager(
	logger hclog.Logger,
	client http.Client) *Manager {
	return &Manager{
		logger:        logger.Named("pluginmgr"),
		wallets:       make(map[string]walletAndInfo),
		exchangeRates: make(map[string]exchangeRatesAndInfo),
		client:        client,
	}
}
