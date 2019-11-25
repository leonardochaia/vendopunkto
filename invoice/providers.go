package invoice

import (
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
	"go.uber.org/zap"
)

var InvoiceProviders = wire.NewSet(NewHandler, NewManager)

func NewHandler(manager *Manager) *Handler {
	return &Handler{
		logger:  zap.S().With("package", "invoice"),
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
