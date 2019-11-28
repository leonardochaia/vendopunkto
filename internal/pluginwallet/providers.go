package pluginwallet

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
)

var PluginWalletProviders = wire.NewSet(NewHandler, NewRouter)

func NewHandler(globalLogger hclog.Logger) *Handler {
	return &Handler{
		logger: globalLogger.Named("plugin-wallet"),
	}
}
