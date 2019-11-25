package invoice

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
	"github.com/rs/xid"
)

type Invoice struct {
	ID             string    `json:"id"`
	Amount         uint      `json:"amount" gorm:"type:BIGINT"`
	Denomination   string    `json:"denomination"`
	PaymentAddress string    `json:"address"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Manager struct {
	wallet wallet.Client
	db     *gorm.DB
}

func (inv *Manager) GetInvoice(id string) (*Invoice, error) {
	var invoice Invoice
	err := inv.db.First(&invoice, "ID = ?", id).Error
	return &invoice, err
}

func (inv *Manager) CreateInvoice(amount uint, denomination string) (*Invoice, error) {

	address, err := inv.wallet.CreateAddress(&wallet.RequestCreateAddress{})
	if err != nil {
		return nil, err
	}

	invoice := &Invoice{
		ID:             xid.New().String(),
		PaymentAddress: address.Address,
		Amount:         amount,
		Denomination:   denomination,
	}

	err = inv.db.Create(invoice).Error

	return invoice, err
}
