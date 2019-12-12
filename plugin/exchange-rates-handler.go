package plugin

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/errors"
)

// ExchangeRatesHandler is the Wallet Plugin HTTP API implementation
// which runs on the plugin server and will be called
// by vendopunkto
// It exposes the wallet interface provided by the
// plugin developer into HTTP endpoints
type ExchangeRatesHandler struct {
	plugin       ExchangeRatesPlugin
	serverPlugin ServerPlugin
}

type GetExchangeRatesParams struct {
	Currency   string   `json:"currency"`
	Currencies []string `json:"currencies"`
}

func NewHandler(plugin ExchangeRatesPlugin, serverPlugin ServerPlugin) *chi.Mux {
	router := chi.NewRouter()

	handler := &ExchangeRatesHandler{
		plugin:       plugin,
		serverPlugin: serverPlugin,
	}

	router.Post("/", errors.WrapHandler(handler.getExchangeRates))

	return router
}

func (handler *ExchangeRatesHandler) getExchangeRates(w http.ResponseWriter, r *http.Request) *errors.APIError {
	var params = new(GetExchangeRatesParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.InvalidRequestParams(err)
	}

	res, err := handler.plugin.GetExchangeRates(params.Currency, params.Currencies)

	if err != nil {
		return errors.InternalServerError(err)
	}

	render.JSON(w, r, res)
	return nil
}
