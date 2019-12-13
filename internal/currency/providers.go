package currency

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
)

var CurrencyProviders = wire.NewSet(NewHandler)

func NewHandler(manager *pluginmgr.Manager, globalLogger hclog.Logger) (Handler, error) {
	return Handler{
		logger:  globalLogger.Named("currency-handler"),
		manager: manager,
	}, nil
}
