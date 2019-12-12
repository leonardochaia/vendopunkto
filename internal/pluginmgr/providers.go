package pluginmgr

import (
	"net/http"
	"time"

	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/currency"
)

var PluginProviders = wire.NewSet(NewManager)

func NewManager(logger hclog.Logger, currency currency.Manager) *Manager {
	return &Manager{
		logger:        logger.Named("pluginmgr"),
		wallets:       make(map[string]walletAndInfo),
		exchangeRates: make(map[string]exchangeRatesAndInfo),
		http: http.Client{
			Timeout: 15 * time.Second,
		},
		currency: currency,
	}
}
