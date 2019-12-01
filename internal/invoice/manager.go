package invoice

import (
	"github.com/jinzhu/gorm"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/rs/xid"
)

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

	inv.db.Model(&invoice).Related(&invoice.Payments)

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
	}

	err = inv.db.Create(invoice).Error

	return invoice, err
}

func (inv *Manager) ConfirmPayment(
	address string,
	confirmations uint,
	amount uint64,
	txHash string) (*Invoice, error) {

	invoice, err := inv.GetInvoiceByAddress(address)

	if err != nil {
		return nil, err
	}

	payment := invoice.FindPayment(txHash)

	// Update the payment
	if payment != nil {
		payment.Update(confirmations)
	} else {
		// New payment
		payment = invoice.AddPayment(txHash, amount, confirmations)
	}

	inv.db.Save(payment)
	return invoice, nil
}
