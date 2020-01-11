package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/errors"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
	"github.com/shopspring/decimal"
)

// CurrencyHandler for currencies
type CurrencyHandler struct {
	plugins      vendopunkto.PluginManager
	currencyRepo vendopunkto.CurrencyRepository
	logger       hclog.Logger
}

// InternalRoutes creates a router for the internal API
func (handler *CurrencyHandler) InternalRoutes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", errors.WrapHandler(handler.getAllCurrencies))
	router.Post("/rates/convert", errors.WrapHandler(handler.getExchange))

	return router
}

func (handler *CurrencyHandler) getAllCurrencies(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.currencies.getAllCurrencies"
	currencies, err := handler.currencyRepo.Search(r.Context())
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	result := []dtos.CurrencyDto{}

	for _, currency := range currencies {
		_, err := handler.plugins.GetWalletInfoForCurrency(currency.Symbol)
		result = append(result, dtos.CurrencyDto{
			Name:             currency.Name,
			Symbol:           currency.Symbol,
			LogoImageURL:     currency.LogoImageURL,
			SupportsPayments: err == nil,
		})
	}

	render.JSON(w, r, result)
	return nil
}

func (handler *CurrencyHandler) getExchange(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.currencies.getExchange"
	var params = new(dtos.GetExchangeParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	if params.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.E(op, errors.Parameters, errors.Str("Amount parameter must be provided"))
	}

	exchange, err := handler.plugins.GetConfiguredExchangeRatesPlugin()
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	rates, err := exchange.GetExchangeRates(params.FromCurrency, params.ToCurrencies)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	result := make(dtos.GetExchangeResult)

	for _, coin := range params.ToCurrencies {
		if rate, ok := rates[coin]; ok {
			result[coin] = rate.Mul(params.Amount)
		}
	}

	render.JSON(w, r, result)
	return nil
}
