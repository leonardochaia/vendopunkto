package invoice

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/errors"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

// invoiceManager contains the business logic for handling invoices
type invoiceManager struct {
	repository    vendopunkto.InvoiceRepository
	logger        hclog.Logger
	pluginManager vendopunkto.PluginManager
	
	topic         vendopunkto.InvoiceTopic
}

func (mgr *invoiceManager) createAddressForInvoice(invoiceID string, currency string) (string, error) {
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

func (mgr *invoiceManager) getDefaultPaymentMethods() ([]string, error) {
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

func (mgr *invoiceManager) addPaymentMethodsToInvoice(
	invoice *vendopunkto.Invoice,
	paymentMethods []dtos.PaymentMethodCreationParams) error {
	exchange, err := mgr.pluginManager.GetConfiguredExchangeRatesPlugin()
	if err != nil {
		return err
	}

	currencies := []string{}
	for _, method := range paymentMethods {
		if method.Total.LessThanOrEqual(decimal.Zero) {
			currencies = append(currencies, method.Currency)
		}
	}

	rates, err := exchange.GetExchangeRates(invoice.Currency, currencies)
	if err != nil {
		return err
	}

	// convert the invoice's total to the paymentMethod's using the rates plugin
	for _, method := range paymentMethods {
		currency := method.Currency
		total := method.Total

		if total.LessThanOrEqual(decimal.Zero) {
			rate, ok := rates[currency]
			if !ok {
				mgr.logger.Warn("Failed to find a rate currency, method ignored",
					"coin", currency,
					"invoiceID", invoice.ID)
				continue
			}

			if rate.LessThanOrEqual(decimal.Zero) {
				mgr.logger.Warn("Invalid rate for currency was provided, method ignored",
					"coin", currency,
					"rate", rate,
					"total", invoice.Total,
					"invoiceID", invoice.ID)
				continue
			}
			total = invoice.Total.Mul(rate)
		}

		// we don't generate addresses ahead of time
		address := ""
		invoice.AddPaymentMethod(currency, address, total)
	}

	return nil
}

// GetInvoice finds an invoice by it's ID
func (mgr *invoiceManager) GetInvoice(
	ctx context.Context,
	id string) (*vendopunkto.Invoice, error) {
	const op errors.Op = "invoicemgr.getInvoice"
	inv, err := mgr.repository.FindByID(ctx, id)

	if err != nil {
		return nil, errors.E(op, err)
	}

	return inv, nil
}

// Search finds invoices matching filter
func (mgr *invoiceManager) Search(
	ctx context.Context,
	filter vendopunkto.InvoiceFilter) ([]vendopunkto.Invoice, error) {

	const op errors.Op = "invoicemgr.searchInvoices"
	list, err := mgr.repository.Search(ctx, filter)

	if err != nil {
		return nil, errors.E(op, err)
	}

	return list, nil
}

// GetInvoiceByAddress finds an invoice by it's the provided paymentMethod's
// address.
func (mgr *invoiceManager) GetInvoiceByAddress(
	ctx context.Context,
	address string) (*vendopunkto.Invoice, error) {
	const op errors.Op = "invoicemgr.getInvoiceByAddress"
	inv, err := mgr.repository.FindByAddress(ctx, address)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return inv, nil
}

// CreateInvoice creates an invoice with the provided total and currency.
// If no payment methods are provided, all supported currencies will be used
func (mgr *invoiceManager) CreateInvoice(
	ctx context.Context,
	params dtos.InvoiceCreationParams) (*vendopunkto.Invoice, error) {
	const op errors.Op = "invoicemgr.create"

	if params.Total.LessThanOrEqual(decimal.Zero) {
		return nil, errors.E(op, errors.Parameters,
			errors.Str("Total parameters must be provided and be positive"))
	}

	invoice := &vendopunkto.Invoice{
		ID:        xid.New().String(),
		Total:     params.Total,
		Currency:  strings.ToLower(params.Currency),
		CreatedAt: time.Now(),
	}

	path := errors.PathName("invoice/" + invoice.ID)

	logger := mgr.logger.With(
		"invoideID", invoice.ID,
		"currency", invoice.Currency,
		"total", invoice.Total.String(),
		"methodCount", len(params.PaymentMethods),
	)

	// populate payment methods using all currencies if none was provided
	if params.PaymentMethods == nil || len(params.PaymentMethods) == 0 {
		logger.Debug("No payment methods were provided. Using defaults.")

		methods, err := mgr.getDefaultPaymentMethods()
		if err != nil {
			return nil, errors.E(op, path, err)
		}
		for _, currency := range methods {
			params.PaymentMethods = append(params.PaymentMethods, dtos.PaymentMethodCreationParams{
				Currency: currency,
				Total:    decimal.Zero,
			})
		}
	}

	err := mgr.addPaymentMethodsToInvoice(invoice, params.PaymentMethods)
	if err != nil {
		return nil, errors.E(op, path, err)
	}

	if len(invoice.PaymentMethods) == 0 {
		err := errors.Str("Failed to create invoice: no payment methods where created")
		logger.Error(err.Error())
		return nil, errors.E(op, path, err)
	}

	// create the address for the default payment method ahead
	defaultMethod := invoice.FindDefaultPaymentMethod()

	if defaultMethod == nil {
		defaultMethod = invoice.PaymentMethods[0]
	}

	address, err := mgr.createAddressForInvoice(invoice.ID, defaultMethod.Currency)
	if err != nil {
		return nil, errors.E(op, path, err)
	}

	defaultMethod.Address = address

	err = mgr.repository.Create(ctx, invoice)

	if err != nil {
		return nil, errors.E(op, path, err)
	}

	logger.Info("Created new invoice")

	return invoice, nil
}

// CreateAddressForPaymentMethod will use the currency wallet to create a unique
// address to receive payments.
func (mgr *invoiceManager) CreateAddressForPaymentMethod(
	ctx context.Context,
	invoiceID string,
	currency string) (*vendopunkto.Invoice, error) {
	const op errors.Op = "invoicemgr.createAddressForPaymentMethod"
	path := errors.PathName("invoice/" + invoiceID + "/" + currency)

	invoice, err := mgr.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, errors.E(op, path, err)
	}

	method := invoice.FindPaymentMethodForCurrency(currency)
	if method == nil {
		return nil, errors.E(op, path,
			errors.Str("Provided invoice does not have a payment method for that currency"))
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
		return nil, errors.E(op, path, err)
	}

	method.Address = address

	err = mgr.repository.UpdatePaymentMethod(ctx, method)

	if err != nil {
		return nil, errors.E(op, path, err)
	}

	mgr.logger.Info("Generated new address for payment method",
		"id", invoice.ID,
		"currency", currency,
		"methodId", method.ID,
		"address", address,
	)

	mgr.topic.Send(*invoice)

	return invoice, nil
}

// ConfirmPayment either creates or updates a payment.
func (mgr *invoiceManager) ConfirmPayment(
	ctx context.Context,
	address string,
	confirmations uint64,
	amount decimal.Decimal,
	txHash string,
	blockHeight uint64) (*vendopunkto.Invoice, error) {
	const op errors.Op = "invoicemgr.confirmPayment"
	path := errors.PathName("address/" + address)

	invoice, err := mgr.GetInvoiceByAddress(ctx, address)

	if err != nil {
		return nil, errors.E(op, path, err)
	}

	method := invoice.FindPaymentMethodForAddress(address)
	if method == nil {
		return nil, errors.E(op, path, errors.NotExist,
			errors.Errorf("Method not found for address %s", address))
	}

	payment := method.FindPayment(txHash)

	logger := mgr.logger.With(
		"invoice", invoice.ID,
		"method", method.Currency,
		"address", address,
		"txHash", txHash,
		"amount", amount,
		"confirmations", confirmations,
		"blockHeight", blockHeight,
	)

	// Update the payment
	if payment != nil {
		if payment.Update(confirmations, blockHeight) {
			logger.Info("Updating payment")
			err = mgr.repository.UpdatePayment(ctx, payment)
		}
	} else {
		// New payment
		logger.Info("Received new payment")
		payment = method.AddPayment(txHash, amount, confirmations, blockHeight)
		err = mgr.repository.CreatePayment(ctx, payment)
	}

	if err != nil {
		return nil, errors.E(op, path, err)
	}

	mgr.topic.Send(*invoice)

	return invoice, nil
}
