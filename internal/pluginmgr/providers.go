package pluginmgr

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
)

var PluginProviders = wire.NewSet(NewManager)

func NewManager(logger hclog.Logger) *Manager {
	return &Manager{
		logger:  logger.Named("pluginmgr"),
		wallets: make(map[string]walletAndInfo),
	}
}
