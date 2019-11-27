package invoice

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/leonardochaia/vendopunkto/pluginmgr"
	"github.com/rs/xid"
)

type Invoice struct {
	ID             string    `json:"id"`
	Amount         uint64    `json:"amount" gorm:"type:BIGINT"`
	Denomination   string    `json:"denomination"`
	PaymentAddress string    `json:"address"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Manager struct {
	db            *gorm.DB
	pluginManager *pluginmgr.Manager
}

func (inv *Manager) GetInvoice(id string) (*Invoice, error) {
	var invoice Invoice

	result := inv.db.First(&invoice, "ID = ?", id)

	if result.RecordNotFound() {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &invoice, nil
}

func (inv *Manager) GetInvoiceByAddress(address string) (*Invoice, error) {
	var invoice Invoice

	result := inv.db.First(&invoice, "payment_address = ?", address)

	if result.RecordNotFound() {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &invoice, nil
}

func (inv *Manager) CreateInvoice(amount uint64, denomination string) (*Invoice, error) {

	newID := xid.New().String()

	wallet, err := inv.pluginManager.GetWallet("monero-wallet")

	if err != nil {
		return nil, err
	}

	address, err := wallet.GenerateNewAddress(newID)
	if err != nil {
		return nil, err
	}

	invoice := &Invoice{
		ID:             newID,
		PaymentAddress: address,
		Amount:         amount,
		Denomination:   denomination,
	}

	err = inv.db.Create(invoice).Error

	return invoice, err
}
