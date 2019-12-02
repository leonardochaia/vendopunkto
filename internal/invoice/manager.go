package invoice

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/rs/xid"
)

type Manager struct {
	db            *gorm.DB
	logger        hclog.Logger
	pluginManager *pluginmgr.Manager
}

func (mgr *Manager) findInvoice(key string, value string) (*Invoice, error) {
	var invoice Invoice

	result := mgr.db.First(&invoice, key+" = ?", value)
	if result.RecordNotFound() {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	mgr.db.Model(&invoice).Related(&invoice.Payments)

	return &invoice, nil
}

func (mgr *Manager) GetInvoice(id string) (*Invoice, error) {
	return mgr.findInvoice("ID", id)
}

func (mgr *Manager) GetInvoiceByAddress(address string) (*Invoice, error) {
	return mgr.findInvoice("payment_address", address)
}

func (mgr *Manager) CreateInvoice(amount uint64, currency string) (*Invoice, error) {

	newID := xid.New().String()

	wallet, err := mgr.pluginManager.GetWalletForCurrency(currency)

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

	err = mgr.db.Create(invoice).Error

	return invoice, err
}

func (mgr *Manager) ConfirmPayment(
	address string,
	confirmations uint64,
	amount uint64,
	txHash string) (*Invoice, error) {

	invoice, err := mgr.GetInvoiceByAddress(address)

	if err != nil {
		return nil, err
	}

	if invoice == nil {
		mgr.logger.Error("Coun't find invoice for address", "address", address,
			"txHash", txHash, "amount", amount)
		return nil, fmt.Errorf("Couldn't find invoice for address " + address)
	}

	mgr.logger.Info("Received payment confirmation", "invoice", invoice.ID,
		"address", address, "txHash", txHash, "amount", amount)
	payment := invoice.FindPayment(txHash)

	// Update the payment
	if payment != nil {
		payment.Update(confirmations)
	} else {
		// New payment
		payment = invoice.AddPayment(txHash, amount, confirmations)
	}

	mgr.db.Save(payment)
	return invoice, nil
}
