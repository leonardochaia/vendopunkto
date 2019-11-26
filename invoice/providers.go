package invoice

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
)

var InvoiceProviders = wire.NewSet(NewHandler, NewManager)

func NewHandler(manager *Manager, globalLogger hclog.Logger) *Handler {
	return &Handler{
		logger:  globalLogger.Named("invoice"),
		manager: manager,
	}
}

func NewManager(db *gorm.DB, wallet wallet.Client) (*Manager, error) {

	db.AutoMigrate(&Invoice{})

	return &Manager{
		wallet: wallet,
		db:     db,
	}, nil
}
