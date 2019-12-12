package currency

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type Currency struct {
	gorm.Model
	Name   string          `json:"name"`
	Symbol string          `json:"symbol"`
	Rate   decimal.Decimal `json:"rate" sql:"type:decimal(20,8);"`
}
