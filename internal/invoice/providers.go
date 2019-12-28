package invoice

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
)

var InvoiceProviders = wire.NewSet(NewHandler, NewManager, NewTopic)

func NewHandler(
	manager *Manager,
	globalLogger hclog.Logger,
	pluginMgr *pluginmgr.Manager,
	topic Topic) *Handler {
	return &Handler{
		logger:    globalLogger.Named("invoice-handler"),
		manager:   manager,
		pluginMgr: pluginMgr,
		topic:     topic,
	}
}

func NewManager(
	repository InvoiceRepository,
	pluginManager *pluginmgr.Manager,
	globalLogger hclog.Logger,
	topic Topic) (*Manager, error) {
	return &Manager{
		logger:        globalLogger.Named("invoice-manager"),
		pluginManager: pluginManager,
		repository:    repository,
		topic:         topic,
	}, nil
}
