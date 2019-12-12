package currency

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
)

var CurrencyProviders = wire.NewSet(NewManager)

func NewManager(db *gorm.DB, globalLogger hclog.Logger) (Manager, error) {

	db.AutoMigrate(&Currency{})

	return Manager{
		logger: globalLogger.Named("currency-manager"),
		db:     db,
	}, nil
}
