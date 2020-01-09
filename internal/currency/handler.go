package currency

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/internal/pluginmgr"
	"github.com/shopspring/decimal"
)

type Handler struct {
	manager *pluginmgr.Manager
	logger  hclog.Logger
}

func (handler *Handler) InternalRoutes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", errors.WrapHandler(handler.getAllCurrencies))
	router.Post("/rates", errors.WrapHandler(handler.getRates))
	router.Post("/rates/convert", errors.WrapHandler(handler.getExchange))

	return router
}

func (handler *Handler) getAllCurrencies(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.currencies.getAllCurrencies"

	currencies, err := handler.manager.GetAllCurrencies()
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	result := []dtos.ExchangeRatesCurrency{
		dtos.ExchangeRatesCurrency{
			Name:             "US Dollars",
			Symbol:           "usd",
			SupportsPayments: false,
		},
	}

	for _, currency := range currencies {
		result = append(result, dtos.ExchangeRatesCurrency{
			Name:             currency.Name,
			Symbol:           currency.Symbol,
			SupportsPayments: true,
		})
	}

	render.JSON(w, r, result)
	return nil
}

func (handler *Handler) getRates(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.currencies.getRates"
	var params = new(dtos.GetCurrencyRatesParams)
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	exchange, err := handler.manager.GetConfiguredExchangeRatesPlugin()
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	rates, err := exchange.GetExchangeRates(params.FromCurrency, params.ToCurrencies)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	render.JSON(w, r, rates)
	return nil
}

func (handler *Handler) getExchange(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.currencies.getRates"
	var params = new(dtos.GetExchangeParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	if params.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.E(op, errors.Parameters, errors.Str("Amount parameter must be provided"))
	}

	exchange, err := handler.manager.GetConfiguredExchangeRatesPlugin()
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
