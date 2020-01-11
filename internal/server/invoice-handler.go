package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/errors"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
)

// InvoiceHandler exposes APIs for interacting with invoices
type InvoiceHandler struct {
	manager   vendopunkto.InvoiceManager
	pluginMgr vendopunkto.PluginManager
	logger    hclog.Logger
	topic     vendopunkto.InvoiceTopic
}

// Routes are the public routes, added to the public API
func (handler *InvoiceHandler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Post("/{id}/payment-method/address", errors.WrapHandler(handler.generatePaymentMethodAddress))
	router.Get("/{id}", errors.WrapHandler(handler.getInvoice))
	router.Get("/{id}/ws", errors.WrapHandler(handler.invoiceWebSocket))

	return router
}

// InternalRoutes are the internal routes, added to the internal API
func (handler *InvoiceHandler) InternalRoutes() chi.Router {
	router := chi.NewRouter()

	router.Get("/{id}", errors.WrapHandler(handler.getInvoice))

	router.Post("/search", errors.WrapHandler(handler.searchInvoices))
	router.Post("/", errors.WrapHandler(handler.createInvoice))

	return router
}

func (handler *InvoiceHandler) createInvoice(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.invoice.create"

	var params = new(dtos.InvoiceCreationParams)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	invoice, err := handler.manager.CreateInvoice(r.Context(), *params)

	if err != nil {
		return errors.E(op, err)
	}

	return handler.renderInvoiceDto(w, r, *invoice)
}

func (handler *InvoiceHandler) getInvoice(w http.ResponseWriter, r *http.Request) error {
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

func (handler *InvoiceHandler) searchInvoices(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.invoice.get"

	var params = new(vendopunkto.InvoiceFilter)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	invoices, err := handler.manager.Search(r.Context(), *params)
	if err != nil {
		return errors.E(op, err)
	}

	return handler.renderInvoiceListDto(w, r, invoices)
}

func (handler *InvoiceHandler) generatePaymentMethodAddress(w http.ResponseWriter, r *http.Request) error {
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

func (handler *InvoiceHandler) renderInvoiceDto(
	w http.ResponseWriter,
	r *http.Request,
	invoice vendopunkto.Invoice) error {
	const op errors.Op = "api.invoice.renderInvoiceDto"

	dto, err := convertInvoiceToDto(invoice, handler.pluginMgr)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	render.JSON(w, r, dto)
	return nil
}

func (handler *InvoiceHandler) renderInvoiceListDto(
	w http.ResponseWriter,
	r *http.Request,
	invoices []vendopunkto.Invoice) error {
	const op errors.Op = "api.invoice.renderInvoiceListDto"

	result := []dtos.InvoiceDto{}
	for _, invoice := range invoices {
		dto, err := convertInvoiceToDto(invoice, handler.pluginMgr)
		if err != nil {
			return errors.E(op, errors.Internal, err)
		}
		result = append(result, dto)
	}

	render.JSON(w, r, result)
	return nil
}

func convertInvoiceToDto(
	invoice vendopunkto.Invoice,
	pluginMgr vendopunkto.PluginManager) (dtos.InvoiceDto, error) {

	dto := &dtos.InvoiceDto{
		ID:                invoice.ID,
		Total:             invoice.Total,
		Currency:          invoice.Currency,
		CreatedAt:         invoice.CreatedAt,
		Status:            uint(invoice.Status()),
		PaymentPercentage: invoice.CalculatePaymentPercentage(),
		Remaining:         invoice.CalculateRemainingAmount(),
		PaymentMethods:    []*dtos.PaymentMethodDto{},
		Payments:          []*dtos.PaymentDto{},
	}

	for _, method := range invoice.PaymentMethods {

		var qrCode string
		if method.Address != "" {
			info, err := pluginMgr.GetWalletInfoForCurrency(method.Currency)
			if err != nil {
				return dtos.InvoiceDto{}, err
			}
			qrCode, err = info.BuildQRCode(method.Address, method.Total)
			if err != nil {
				return dtos.InvoiceDto{}, err
			}
		}

		methodDto := &dtos.PaymentMethodDto{
			ID:        method.ID,
			Total:     method.Total,
			Currency:  method.Currency,
			Address:   method.Address,
			Remaining: invoice.CalculatePaymentMethodRemaining(*method),
			QRCode:    qrCode,
		}

		for _, payment := range method.Payments {
			paymentDto := &dtos.PaymentDto{
				TxHash:        payment.TxHash,
				Amount:        payment.Amount,
				Confirmations: payment.Confirmations,
				ConfirmedAt:   payment.ConfirmedAt,
				CreatedAt:     payment.CreatedAt,
				Status:        uint(payment.Status()),
				Currency:      method.Currency,
				BlockHeight:   payment.BlockHeight,
			}
			dto.Payments = append(dto.Payments, paymentDto)
		}

		dto.PaymentMethods = append(dto.PaymentMethods, methodDto)
	}

	return *dto, nil
}
