package invoice

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
)

var InvoiceProviders = wire.NewSet(NewHandler, NewManager)

func NewHandler(manager *Manager, globalLogger hclog.Logger) *Handler {
	return &Handler{
		logger:  globalLogger.Named("invoice-handler"),
		manager: manager,
	}
}

func NewManager(
	repository InvoiceRepository,
	pluginManager *pluginmgr.Manager,
	globalLogger hclog.Logger) (*Manager, error) {
	return &Manager{
		logger:        globalLogger.Named("invoice-manager"),
		pluginManager: pluginManager,
		repository:    repository,
	}, nil
}
