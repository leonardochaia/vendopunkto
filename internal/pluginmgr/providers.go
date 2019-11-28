package pluginmgr

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/pluginwallet"
)

var PluginProviders = wire.NewSet(NewManager)

func NewManager(
	logger hclog.Logger,
	walletRouter *pluginwallet.Router) *Manager {
	return &Manager{
		logger:       logger.Named("pluginmgr"),
		wallets:      make(map[string]walletAndInfo),
		walletRouter: walletRouter,
	}
}
