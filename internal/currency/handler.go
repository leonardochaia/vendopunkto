package currency

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
)

type Handler struct {
	manager *pluginmgr.Manager
	logger  hclog.Logger
}

func (handler *Handler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", errors.WrapHandler(handler.getAllCurrencies))

	return router
}

func (handler *Handler) getAllCurrencies(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.currencies.getAllCurrencies"

	currencies, err := handler.manager.GetAllCurrencies()
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	render.JSON(w, r, currencies)
	return nil
}
