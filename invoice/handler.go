package invoice

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/errors"
	"go.uber.org/zap"
)

type Handler struct {
	manager *Manager
	logger  *zap.SugaredLogger
}

func (handler *Handler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Post("/", errors.WrapHandler(handler.createInvoice))
	router.Get("/{id}", errors.WrapHandler(handler.getInvoice))

	return router
}

func (handler *Handler) createInvoice(w http.ResponseWriter, r *http.Request) *errors.APIError {

	type creationParams struct {
		Amount       uint64 `json:"amount"`
		Denomination string `json:"denomination"`
	}

	var params = new(creationParams)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.InvalidRequestParams(err)
	}

	invoice, err := handler.manager.CreateInvoice(params.Amount, params.Denomination)
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
