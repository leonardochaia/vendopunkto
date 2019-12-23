package invoice

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
)

// Handler exposes APIs for interacting with invoices
type Handler struct {
	manager   *Manager
	pluginMgr *pluginmgr.Manager
	logger    hclog.Logger
}

// Routes are the public routes, added to the public API
func (handler *Handler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Post("/{id}/payment-method/address", errors.WrapHandler(handler.generatePaymentMethodAddress))
	router.Get("/{id}", errors.WrapHandler(handler.getInvoice))

	return router
}

// InternalRoutes are the internal routes, added to the internal API
func (handler *Handler) InternalRoutes() chi.Router {
	router := chi.NewRouter()

	router.Post("/", errors.WrapHandler(handler.createInvoice))

	router.Post("/payments/confirm", errors.WrapHandler(handler.confirmPayment))

	return router
}

func (handler *Handler) createInvoice(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.invoice.create"
	var params = new(dtos.InvoiceCreationParams)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	if params.Total == 0 {
		return errors.E(op, errors.Parameters, errors.Str("A total parameters must be provided"))
	}

	invoice, err := handler.manager.CreateInvoice(r.Context(),
		params.Total, params.Currency, params.PaymentMethods)

	if err != nil {
		return errors.E(op, err)
	}

	return handler.renderInvoiceDto(w, r, *invoice)
}

func (handler *Handler) getInvoice(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.invoice.get"

	invoiceID := chi.URLParam(r, "id")
	if invoiceID == "" {
		return errors.E(op, errors.Parameters, errors.Str("No ID was provided"))
	}

	invoice, err := handler.manager.GetInvoice(r.Context(), invoiceID)
	if err != nil {
		return errors.E(op, err)
	}

	return handler.renderInvoiceDto(w, r, *invoice)
}

func (handler *Handler) generatePaymentMethodAddress(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.invoice.generatePaymentMethodAddress"

	invoiceID := chi.URLParam(r, "id")
	if invoiceID == "" {
		return errors.E(op, errors.Parameters, errors.Str("No ID was provided"))
	}

	var params = new(dtos.InvoiceGeneratePaymentMethodAddressParams)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	invoice, err := handler.manager.CreateAddressForPaymentMethod(r.Context(),
		invoiceID, params.Currency)

	if err != nil {
		return errors.E(op, err)
	}

	return handler.renderInvoiceDto(w, r, *invoice)
}

// confirmPayment is an internal endpoint that confirms an invoice has been paid
func (handler *Handler) confirmPayment(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.invoice.confirmPayment"

	var params = new(dtos.InvoiceConfirmPaymentsParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	_, err := handler.manager.ConfirmPayment(r.Context(),
		params.Address,
		params.Confirmations,
		params.Amount,
		params.TxHash,
	)

	if err != nil {
		return errors.E(op, err)
	}

	render.NoContent(w, r)
	return nil
}

func (handler *Handler) renderInvoiceDto(
	w http.ResponseWriter,
	r *http.Request,
	invoice Invoice) error {
	const op errors.Op = "api.invoice.renderInvoiceDto"

	dto, err := convertInvoiceToDto(invoice, handler.pluginMgr)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	render.JSON(w, r, dto)
	return nil
}
