package pluginmgr

import (
	"net/http"
	"time"

	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
)

var PluginProviders = wire.NewSet(NewManager)

func NewManager(logger hclog.Logger) *Manager {
	return &Manager{
		logger:  logger.Named("pluginmgr"),
		wallets: make(map[string]walletAndInfo),
		http: http.Client{
			Timeout: 15 * time.Second,
		},
	}
}
