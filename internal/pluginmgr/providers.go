package pluginmgr

import (
	"net/http"

	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/currency"
)

var PluginProviders = wire.NewSet(NewManager)

func NewManager(
	logger hclog.Logger,
	currency currency.Manager,
	client http.Client) *Manager {
	return &Manager{
		logger:        logger.Named("pluginmgr"),
		wallets:       make(map[string]walletAndInfo),
		exchangeRates: make(map[string]exchangeRatesAndInfo),
		client:        client,
		currency:      currency,
	}
}
