package invoice

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/rs/xid"
)

type Status int

const (
	// Pending when the invoice has been created but the wallet has not
	// received payment yet.
	Pending Status = 1
	// Mempool when the wallet has received the payment, and it's waiting
	// to be included in a block.
	Mempool Status = 2
	// Confirmed when the wallet has received enough confirmations.
	Confirmed Status = 3
	// Failed when the wallet fails to process the transaction for
	// whatever reason
	Failed Status = 4
)

type Invoice struct {
	ID             string    `json:"id"`
	Amount         uint64    `json:"amount" gorm:"type:BIGINT"`
	Currency       string    `json:"currency"`
	PaymentAddress string    `json:"address" gorm:"index:inv_addr"`
	Status         Status    `json:"status"`
	Confirmations  uint      `json:"confirmations"`
	TxHash         string    `json:"txHash"`
	ConfirmedAt    time.Time `json:"confirmedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Manager struct {
	db            *gorm.DB
	pluginManager *pluginmgr.Manager
}

func (inv *Manager) findInvoice(key string, value string) (*Invoice, error) {
	var invoice Invoice

	result := inv.db.First(&invoice, key+" = ?", value)

	if result.RecordNotFound() {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &invoice, nil
}

func (inv *Manager) GetInvoice(id string) (*Invoice, error) {
	return inv.findInvoice("ID", id)
}

func (inv *Manager) GetInvoiceByAddress(address string) (*Invoice, error) {
	return inv.findInvoice("payment_address", address)
}

func (inv *Manager) CreateInvoice(amount uint64, currency string) (*Invoice, error) {

	newID := xid.New().String()

	wallet, err := inv.pluginManager.GetWalletForCurrency(currency)

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
		Currency:       currency,
		Status:         Pending,
	}

	err = inv.db.Create(invoice).Error

	return invoice, err
}

func (inv *Manager) ConfirmPayment(
	address string,
	confirmations uint,
	txHash string) (*Invoice, error) {

	invoice, err := inv.GetInvoiceByAddress(address)

	if err != nil {
		return nil, err
	}

	if invoice.TxHash == "" {
		invoice.TxHash = txHash
	} else {
		if invoice.TxHash != txHash {
			return nil, fmt.Errorf("Invoice TX Hash does not match provided TX Hash")
		}
	}

	// Wallet's always win.
	// This is done with the intention to support chain reorgs
	invoice.Confirmations = confirmations

	if invoice.Confirmations == 0 {
		invoice.Status = Mempool
	} else
	// TODO: min confirmation amount should be a setting
	if invoice.Confirmations > 0 {
		if invoice.Status != Confirmed {
			invoice.ConfirmedAt = time.Now()
		}
		invoice.Status = Confirmed
	}

	inv.db.Save(invoice)
	return invoice, nil
}
