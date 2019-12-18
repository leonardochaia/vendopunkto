package invoice

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/leonardochaia/vendopunkto/unit"
)

type Handler struct {
	manager   *Manager
	pluginMgr *pluginmgr.Manager
	logger    hclog.Logger
}

func (handler *Handler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Post("/", errors.WrapHandler(handler.createInvoice))
	router.Post("/{id}/payment-method/address", errors.WrapHandler(handler.generatePaymentMethodAddress))
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
		Total          unit.AtomicUnit `json:"total"`
		Currency       string          `json:"currency"`
		PaymentMethods []string        `json:"paymentMethods"`
	}

	var params = new(creationParams)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.InvalidRequestParams(err)
	}

	invoice, err := handler.manager.CreateInvoice(r.Context(),
		params.Total, params.Currency, params.PaymentMethods)

	if err != nil {
		return errors.InternalServerError(err)
	}

	return handler.renderInvoiceDto(w, r, *invoice)
}

func (handler *Handler) getInvoice(w http.ResponseWriter, r *http.Request) *errors.APIError {

	invoiceID := chi.URLParam(r, "id")
	if invoiceID == "" {
		return errors.InvalidRequestParams(fmt.Errorf("No ID was provided"))
	}

	invoice, err := handler.manager.GetInvoice(r.Context(), invoiceID)
	if err != nil {
		return errors.InternalServerError(err)
	}

	if invoice == nil {
		return errors.ResourceNotFound()
	}

	return handler.renderInvoiceDto(w, r, *invoice)
}

func (handler *Handler) generatePaymentMethodAddress(w http.ResponseWriter, r *http.Request) *errors.APIError {

	invoiceID := chi.URLParam(r, "id")
	if invoiceID == "" {
		return errors.InvalidRequestParams(fmt.Errorf("No ID was provided"))
	}

	type inputParams struct {
		Currency string `json:"currency"`
	}

	var params = new(inputParams)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.InvalidRequestParams(err)
	}

	invoice, err := handler.manager.CreateAddressForPaymentMethod(r.Context(),
		invoiceID, params.Currency)

	if err != nil {
		return errors.InternalServerError(err)
	}

	return handler.renderInvoiceDto(w, r, *invoice)
}

// confirmPayment is an internal endpoint that confirms an invoice has been paid
func (handler *Handler) confirmPayment(w http.ResponseWriter, r *http.Request) *errors.APIError {

	type confirmPaymentsParams struct {
		TxHash        string          `json:"txHash"`
		Address       string          `json:"address"`
		Amount        unit.AtomicUnit `json:"amount"`
		Confirmations uint64          `json:"confirmations"`
	}
	var params = new(confirmPaymentsParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.InvalidRequestParams(err)
	}

	_, err := handler.manager.ConfirmPayment(r.Context(),
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

func (handler *Handler) renderInvoiceDto(
	w http.ResponseWriter,
	r *http.Request,
	invoice Invoice) *errors.APIError {

	dto, err := convertInvoiceToDto(invoice, handler.pluginMgr)
	if err != nil {
		return errors.InternalServerError(err)
	}

	render.JSON(w, r, dto)
	return nil
}
