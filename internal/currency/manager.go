package currency

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type Manager struct {
	db     *gorm.DB
	logger hclog.Logger
}

func (mgr *Manager) findCurrency(value string) (*Currency, error) {
	var currency Currency

	result := mgr.db.Where("name = ? OR symbol = ?", value, value).First(&currency)

	if result.RecordNotFound() {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	// mgr.db.Model(&currency).Related(&currency.Payments)

	return &currency, nil
}

func (mgr *Manager) GetAll() ([]Currency, error) {
	var currencies []Currency
	result := mgr.db.Select(&currencies)

	if result.Error != nil {
		return nil, result.Error
	}

	return currencies, nil
}

func (mgr *Manager) UpdateRate(coin string, rate decimal.Decimal) (*Currency, error) {

	currency, err := mgr.findCurrency(coin)

	if err != nil {
		return nil, err
	}

	if currency == nil {
		mgr.logger.Error("Coun't find coin with name/symbol", "coin", coin)
		return nil, fmt.Errorf("Couldn't find coin with name/symbol " + coin)
	}

	mgr.logger.Info("Updating currency rate", "coin", coin,
		"oldRate", currency.Rate,
		"newRate", rate)

	currency.Rate = rate

	mgr.db.Save(currency)
	return currency, nil
}

func (mgr *Manager) RegisterCurrency(name string, symbol string) (*Currency, error) {
	currency, err := mgr.findCurrency(symbol)

	if err != nil {
		return nil, err
	}

	if currency != nil {
		return currency, nil
	}

	mgr.logger.Info("Registering new currency", "currency", name, "symbol", symbol)

	currency = &Currency{
		Name:   name,
		Symbol: symbol,
	}

	err = mgr.db.Create(currency).Error

	if err != nil {
		return nil, err
	}

	return currency, nil
}
