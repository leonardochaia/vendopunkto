package invoice

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
)

type Handler struct {
	manager *Manager
	logger  hclog.Logger
}

func (handler *Handler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Post("/", errors.WrapHandler(handler.createInvoice))
	router.Get("/{id}", errors.WrapHandler(handler.getInvoice))

	return router
}

func (handler *Handler) InternalRoutes() chi.Router {
	router := chi.NewRouter()

	router.Post("/payments/confirm", errors.WrapHandler(handler.confirmPayment))

	return router
}

func (handler *Handler) createInvoice(w http.ResponseWriter, r *http.Request) *errors.APIError {

	type creationParams struct {
		Amount   uint64 `json:"amount"`
		Currency string `json:"currency"`
	}

	var params = new(creationParams)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.InvalidRequestParams(err)
	}

	invoice, err := handler.manager.CreateInvoice(params.Amount, params.Currency)
	if err != nil {
		return errors.InternalServerError(err)
	}

	render.JSON(w, r, invoice)
	return nil
}

func (handler *Handler) getInvoice(w http.ResponseWriter, r *http.Request) *errors.APIError {

	invoiceID := chi.URLParam(r, "id")
	if invoiceID == "" {
		return errors.InvalidRequestParams(fmt.Errorf("No ID was provided"))
	}

	invoice, err := handler.manager.GetInvoice(invoiceID)
	if err != nil {
		return errors.InternalServerError(err)
	}
	if invoice == nil {
		return errors.ResourceNotFound()
	}

	render.JSON(w, r, invoice)
	return nil
}

// confirmPayment is an internal endpoint that confirms an invoice has been paid
func (handler *Handler) confirmPayment(w http.ResponseWriter, r *http.Request) *errors.APIError {

	type confirmPaymentsParams struct {
		TxHash        string `json:"txHash"`
		Address       string `json:"address"`
		Amount        uint64 `json:"amount"`
		Confirmations uint   `json:"confirmations"`
	}
	var params = new(confirmPaymentsParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.InvalidRequestParams(err)
	}

	_, err := handler.manager.ConfirmPayment(
		params.Address,
		params.Confirmations,
		params.Amount,
		params.TxHash,
	)

	if err != nil {
		return errors.InternalServerError(err)
	}

	render.NoContent(w, r)
	return nil
}
