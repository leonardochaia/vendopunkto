package monero

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
)

// Handler is Monero specific. This will listen for the wallet cli calls
type Handler struct {
	client wallet.Client
	logger hclog.Logger
	server plugin.Server
}

func (handler *Handler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Post("/txNotify", errors.WrapHandler(handler.txNotify))

	return router
}

// txNotify to be called on monero-wallet-cli "txNotify" flag
func (handler *Handler) txNotify(w http.ResponseWriter, r *http.Request) *errors.APIError {
	txHash := r.URL.Query().Get("txHash")
	if txHash == "" {
		return errors.InternalServerError(fmt.Errorf("Coldn't get txHash from query params"))
	}

	resp, err := handler.client.GetTransferByTxID(&wallet.RequestGetTransferByTxID{
		TxID: txHash,
	})

	if err != nil {
		handler.logger.Error("Failed to obtained TX from wallet RPC", "error", err.Error())
		return errors.InternalServerError(err)
	}

	addrResp, err := handler.client.MakeIntegratedAddress(&wallet.RequestMakeIntegratedAddress{
		StandardAddress: resp.Transfer.Address,
		PaymentID:       resp.Transfer.PaymentID,
	})

	if err != nil {
		handler.logger.Error("Failed to make integrated address", "error", err.Error())
		return errors.InternalServerError(err)
	}

	var (
		addr          = addrResp.IntegratedAddress
		amount        = resp.Transfer.Amount
		confirmations = resp.Transfer.Confirmations
	)
	handler.logger.Info("Obtained TX", "address", addr, "amount", resp.Transfer.Amount)

	client, err := handler.server.GetInternalClient()

	if err != nil {
		handler.logger.Error("Failed to obtain host client", "error", err.Error())
		return errors.InternalServerError(err)
	}

	err = client.ConfirmPayment(addr, amount, txHash, confirmations)
	if err != nil {
		handler.logger.Error("Failed to confirm payment with host", "error", err.Error())
		return errors.InternalServerError(err)
	}

	handler.logger.Info("Payment confirmed", "txHash", txHash,
		"amount", amount, "conf", confirmations)

	render.NoContent(w, r)
	return nil
}
