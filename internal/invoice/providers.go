package invoice

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
)

var InvoiceProviders = wire.NewSet(NewHandler, NewManager)

func NewHandler(manager *Manager, globalLogger hclog.Logger) *Handler {
	return &Handler{
		logger:  globalLogger.Named("invoice"),
		manager: manager,
	}
}

func NewManager(db *gorm.DB, pluginManager *pluginmgr.Manager) (*Manager, error) {

	db.AutoMigrate(&Invoice{})
	db.AutoMigrate(&Payment{})

	return &Manager{
		pluginManager: pluginManager,
		db:            db,
	}, nil
}
