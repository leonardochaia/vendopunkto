package plugin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/errors"
)

// CurrencyMetadataHandler exposes the plugin interface through HTTP
type CurrencyMetadataHandler struct {
	plugin       CurrencyMetadataPlugin
	serverPlugin ServerPlugin
}

type GetCurrencyParams struct {
	Currencies []string `json:"currencies"`
}

func NewCurrencyMetadataHandler(plugin CurrencyMetadataPlugin, serverPlugin ServerPlugin) *chi.Mux {
	router := chi.NewRouter()

	handler := &CurrencyMetadataHandler{
		plugin:       plugin,
		serverPlugin: serverPlugin,
	}

	router.Post("/", errors.WrapHandler(handler.getCurrencies))

	return router
}

func (handler *CurrencyMetadataHandler) getCurrencies(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "plugin.api.getCurrencies"
	var params = new(GetCurrencyParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	res, err := handler.plugin.GetCurrencies(params.Currencies)

	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	render.JSON(w, r, res)
	return nil
}
