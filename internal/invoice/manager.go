package invoice

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/leonardochaia/vendopunkto/unit"
	"github.com/rs/xid"
)

type Manager struct {
	db            *gorm.DB
	logger        hclog.Logger
	pluginManager *pluginmgr.Manager
}

func (mgr *Manager) findInvoice(key string, value string) (*Invoice, error) {
	var invoice Invoice

	result := mgr.db.Debug().
		Preload("PaymentMethods").
		Preload("PaymentMethods.Payments").
		First(&invoice, key+" = ?", value)

	if result.RecordNotFound() {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &invoice, nil
}

func (mgr *Manager) createAddressForInvoice(invoiceID string, currency string) (string, error) {
	wallet, err := mgr.pluginManager.GetWalletForCurrency(currency)

	if err != nil {
		return "", err
	}

	address, err := wallet.GenerateNewAddress(invoiceID)
	if err != nil {
		return "", err
	}

	return address, nil
}

func (mgr *Manager) getDefaultPaymentMethods() ([]string, error) {
	wallets, err := mgr.pluginManager.GetAllCurrencies()
	if err != nil {
		return nil, err
	}

	output := []string{}
	for _, wallet := range wallets {
		output = append(output, wallet.Symbol)
	}

	return output, nil
}

func (mgr *Manager) addPaymentMethodsToInvoice(
	invoice *Invoice,
	paymentMethods []string) error {

	// TODO: default rates
	exchange, err := mgr.pluginManager.GetExchangeRatesPlugin("gecko-exchange-rates")
	if err != nil {
		return err
	}

	rates, err := exchange.GetExchangeRates(invoice.Currency, paymentMethods)
	if err != nil {
		return err
	}

	// convert the invoice's total to the paymentMethod's using the rates plugin
	for _, coin := range paymentMethods {
		rate, ok := rates[coin]
		if !ok {
			mgr.logger.Warn("Failed to find a rate currency, method ignored",
				"coin", coin,
				"invoiceID", invoice.ID)
			continue
		}

		if rate == 0 {
			mgr.logger.Warn("Invalid rate for currency was provided, method ignored",
				"coin", coin,
				"rate", rate,
				"total", invoice.Total,
				"invoiceID", invoice.ID)
			continue
		}

		// we don't generate addresses ahead of time
		address := ""
		totalConverted := invoice.Total.Float64() * rate
		invoice.AddPaymentMethod(coin, address, unit.NewFromFloat(totalConverted))
	}

	return nil
}

func (mgr *Manager) GetInvoice(id string) (*Invoice, error) {
	return mgr.findInvoice("ID", id)
}

func (mgr *Manager) GetInvoiceByAddress(address string) (*Invoice, error) {
	var po PaymentMethod
	result := mgr.db.Where("address = ?", address).First((&po))

	if result.RecordNotFound() {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return mgr.findInvoice("ID", po.InvoiceID)
}

func (mgr *Manager) CreateInvoice(
	total unit.AtomicUnit,
	currency string,
	paymentMethods []string) (*Invoice, error) {

	currency = strings.ToLower(currency)

	invoice := &Invoice{
		ID:       xid.New().String(),
		Total:    total,
		Currency: currency,
	}

	// populate payment methods using all currencies if none was provided
	if paymentMethods == nil || len(paymentMethods) == 0 {
		methods, err := mgr.getDefaultPaymentMethods()
		if err != nil {
			return nil, err
		}
		paymentMethods = methods
	}

	err := mgr.addPaymentMethodsToInvoice(invoice, paymentMethods)
	if err != nil {
		return nil, err
	}

	if len(invoice.PaymentMethods) == 0 {
		err := fmt.Errorf("Failed to create invoice: no payment methods where created")
		mgr.logger.Error(err.Error(), "coin", currency, "methods", paymentMethods)
		return nil, err
	}

	// create the address for the default payment method ahead
	defaultMethod := invoice.FindDefaultPaymentMethod()

	if defaultMethod == nil {
		defaultMethod = invoice.PaymentMethods[0]
	}

	address, err := mgr.createAddressForInvoice(invoice.ID, defaultMethod.Currency)
	if err != nil {
		return nil, err
	}

	defaultMethod.Address = address

	err = mgr.db.Create(invoice).Error

	if err != nil {
		return nil, err
	}

	mgr.logger.Info("Created new invoice",
		"id", invoice.ID,
		"total", invoice.Total.String(),
		"currency", invoice.Currency,
		"paymentMethods", paymentMethods)

	return mgr.GetInvoice(invoice.ID)
}

func (mgr *Manager) CreateAddressForPaymentMethod(invoiceID string, currency string) (*Invoice, error) {
	invoice, err := mgr.GetInvoice(invoiceID)
	if err != nil {
		return nil, err
	}

	method := invoice.FindPaymentMethodForCurrency(currency)
	if method == nil {
		mgr.logger.Error("Invalid currency for invoice",
			"id", invoice.ID,
			"currency", currency)
		return nil, fmt.Errorf("Provided invoice does not have a payment method for that currency")
	}

	if method.Address != "" {
		mgr.logger.Warn("Requested address generation for method that already has address",
			"id", invoice.ID,
			"currency", currency,
			"methodId", method.ID)

		return invoice, nil
	}

	address, err := mgr.createAddressForInvoice(invoice.ID, currency)
	if err != nil {
		return nil, err
	}

	method.Address = address

	err = mgr.db.Save(method).Error

	if err != nil {
		return nil, err
	}

	mgr.logger.Info("Generated new address for payment method",
		"id", invoice.ID,
		"currency", currency,
		"methodId", method.ID,
		"address", address,
	)

	return invoice, nil
}

func (mgr *Manager) ConfirmPayment(
	address string,
	confirmations uint64,
	amount unit.AtomicUnit,
	txHash string) (*Invoice, error) {

	invoice, err := mgr.GetInvoiceByAddress(address)

	if err != nil {
		return nil, err
	}

	if invoice == nil {
		mgr.logger.Error("Couldn't find invoice for address", "address", address,
			"txHash", txHash, "amount", amount)
		return nil, fmt.Errorf("Couldn't find invoice for address " + address)
	}

	method := invoice.FindPaymentMethodForAddress(address)
	payment := method.FindPayment(txHash)

	mgr.logger.Info("Received payment confirmation",
		"invoice", invoice.ID,
		"method", method.Currency,
		"address", address,
		"txHash", txHash,
		"amount", amount)

	// Update the payment
	if payment != nil {
		payment.Update(confirmations)
	} else {
		// New payment
		payment = method.AddPayment(txHash, amount, confirmations)
	}

	mgr.db.Save(payment)
	return invoice, nil
}
