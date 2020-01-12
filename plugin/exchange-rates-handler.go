package plugin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/errors"
)

// ExchangeRatesHandler exposes the plugin through HTTP
type ExchangeRatesHandler struct {
	plugin       ExchangeRatesPlugin
	serverPlugin ServerPlugin
}

type GetExchangeRatesParams struct {
	Currency   string   `json:"currency"`
	Currencies []string `json:"currencies"`
}

func NewExchangeRatesHandler(plugin ExchangeRatesPlugin, serverPlugin ServerPlugin) *chi.Mux {
	router := chi.NewRouter()

	handler := &ExchangeRatesHandler{
		plugin:       plugin,
		serverPlugin: serverPlugin,
	}

	router.Post("/", errors.WrapHandler(handler.getExchangeRates))

	return router
}

func (handler *ExchangeRatesHandler) getExchangeRates(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "plugin.api.getExchangeRates"
	var params = new(GetExchangeRatesParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	res, err := handler.plugin.GetExchangeRates(params.Currency, params.Currencies)

	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	render.JSON(w, r, res)
	return nil
}
