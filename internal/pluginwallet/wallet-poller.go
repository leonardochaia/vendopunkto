package pluginwallet

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
	"github.com/leonardochaia/vendopunkto/internal/conf"
	"github.com/leonardochaia/vendopunkto/internal/store"
	"github.com/leonardochaia/vendopunkto/plugin"
)

// WalletPoller will poll pending invoice's wallets to obtain payment confirmations
type WalletPoller struct {
	ticker *time.Ticker // periodic ticker

	logger      hclog.Logger
	txBuilder   store.TransactionBuilder
	pluginMgr   vendopunkto.PluginManager
	invoiceRepo vendopunkto.InvoiceRepository
	invoiceMgr  vendopunkto.InvoiceManager
	running     bool
}

// NewPoller creates a new WalletPoller
func NewPoller(
	logger hclog.Logger,
	runtimeConf *conf.Runtime,
	txBuilder store.TransactionBuilder,
	pluginMgr vendopunkto.PluginManager,
	invoiceRepo vendopunkto.InvoiceRepository,
	invoiceMgr vendopunkto.InvoiceManager) (*WalletPoller, error) {

	interval := runtimeConf.GetWalletPollingInterval()

	poller := &WalletPoller{
		ticker:      time.NewTicker(interval),
		txBuilder:   txBuilder,
		pluginMgr:   pluginMgr,
		invoiceRepo: invoiceRepo,
		invoiceMgr:  invoiceMgr,
		logger:      logger.Named("wallet-poller"),
	}
	return poller, nil
}

// Start starts the poller
func (poller *WalletPoller) Start() {
	poller.logger.Info("Wallet poller started")
	for {
		select {
		case <-poller.ticker.C:
			if poller.running {
				poller.logger.Warn("Wallet poller already running, skipping.")
				continue
			}
			poller.running = true
			err := poller.doPoll()
			poller.running = false
			if err != nil {
				poller.logger.Error("Wallet Poller errored", "error", err)
			}
		}
	}
}

// Stop stops the ticker
func (poller *WalletPoller) Stop() {
	poller.logger.Info("Wallet poller stopping")
	poller.ticker.Stop()
}

func (poller *WalletPoller) doPoll() error {

	ctx, tx, err := poller.txBuilder.BuildLazyTransactionContext(context.TODO())
	if err != nil {
		return err
	}

	heights, err := poller.invoiceRepo.GetMaxBlockHeightForCurrencies(ctx)
	if err != nil {
		return err
	}

	for currency, height := range heights {
		wallet, err := poller.pluginMgr.GetWalletForCurrency(currency)
		info, err := wallet.GetWalletInfo()
		if err != nil {
			poller.logger.Error("Error while obtaining wallet info",
				"error", err)
			continue
		}

		params := plugin.WalletPluginIncomingTransferParams{
			MinBlockHeight: height - 10, // TODO: Settings
		}

		transfers, err := wallet.GetIncomingTransfers(params)
		if err != nil {
			poller.logger.Error("Error while obtaining incoming transfer from wallet",
				"error", err)
			continue
		}

		if len(transfers) == 0 {
			continue
		}

		// match transfers to payments
		poller.logger.Info("Got transfers from wallet",
			"wallet", info.Currency.Name,
			"minHeight", height,
			"transfers", len(transfers))

		for _, transfer := range transfers {
			_, err := poller.invoiceMgr.ConfirmPayment(ctx, transfer.Address,
				transfer.Confirmations, transfer.Amount,
				transfer.TxHash, transfer.BlockHeight)
			if err != nil {
				if vpErr, ok := err.(*errors.Error); ok && vpErr.Kind == errors.NotExist {
					poller.logger.Debug("Received wallet transfer for unknown invoice",
						"txHash", transfer.TxHash,
						"address", transfer.Address,
						"error", err,
					)
					continue
				}
				poller.logger.Error("Error confirming payment",
					"txHash", transfer.TxHash,
					"address", transfer.Address,
					"error", err)
			}
		}
	}

	return tx.CommitIfNeeded()
}
