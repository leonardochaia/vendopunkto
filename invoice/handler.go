package invoice

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/store"
	"go.uber.org/zap"
)

type Handler struct {
	manager *Manager
	logger  *zap.SugaredLogger
}

func NewHandler(manager *Manager) *Handler {
	return &Handler{
		logger:  zap.S().With("package", "invoice"),
		manager: manager,
	}
}

func (handler *Handler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Post("/", handler.CreateInvoice)
	router.Get("/{id}", handler.GetInvoice)

	return router
}

func (handler *Handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {

	type creationParams struct {
		Amount       uint   `json:"amount"`
		Denomination string `json:"denomination"`
	}
	var params = new(creationParams)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		render.Render(w, r, errors.ErrInvalidRequest(err))
		return
	}

	invoice, err := handler.manager.CreateInvoice(r.Context(), params.Amount, params.Denomination)
	if err != nil {
		render.Render(w, r, errors.ErrInternal(err))
		return
	}
	render.JSON(w, r, invoice)
}

func (handler *Handler) GetInvoice(w http.ResponseWriter, r *http.Request) {

	invoiceID := chi.URLParam(r, "id")
	if invoiceID == "" {
		render.Render(w, r, errors.ErrInvalidRequest(fmt.Errorf("Invalid ID")))
		return
	}

	b, err := handler.manager.GetInvoice(r.Context(), invoiceID)

	if err == store.ErrNotFound {
		handler.logger.Infow("Provided unexistant invoice ID", "ID", invoiceID)
		render.Render(w, r, errors.ErrNotFound)
		return
	} else if err != nil {
		render.Render(w, r, errors.ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, b)
}
