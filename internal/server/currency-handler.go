package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/dtos"
	"github.com/leonardochaia/vendopunkto/errors"
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
	"github.com/leonardochaia/vendopunkto/internal/conf"
	"github.com/shopspring/decimal"
)

// CurrencyHandler for currencies
type CurrencyHandler struct {
	plugins      vendopunkto.PluginManager
	currencyRepo vendopunkto.CurrencyRepository
	logger       hclog.Logger
	runtimeConf  *conf.Runtime
}

// InternalRoutes creates a router for the internal API
func (handler *CurrencyHandler) InternalRoutes() chi.Router {
	router := chi.NewRouter()

	router.Get("/pricing", errors.WrapHandler(handler.getPricingCurrencies))
	router.Get("/payment-methods", errors.WrapHandler(handler.getPaymentMethodCurrencies))
	router.Post("/rates/convert", errors.WrapHandler(handler.getExchange))

	return router
}

func (handler *CurrencyHandler) getPricingCurrencies(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.currencies.getPricingCurrencies"

	pricingCurrencies := handler.runtimeConf.GetPricingCurrencies()
	result := []dtos.CurrencyDto{}
	if len(pricingCurrencies) == 0 {
		render.JSON(w, r, result)
		return nil
	}

	currencies, err := handler.currencyRepo.FindBySymbols(r.Context(), pricingCurrencies)

	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

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

func (handler *CurrencyHandler) getPaymentMethodCurrencies(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.currencies.getPaymentMethodCurrencies"

	paymentMethods := handler.runtimeConf.GetPaymentMethods()
	result := []dtos.CurrencyDto{}
	if len(paymentMethods) == 0 {
		render.JSON(w, r, result)
		return nil
	}

	currencies, err := handler.currencyRepo.FindBySymbols(r.Context(), paymentMethods)

	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	for _, currency := range currencies {
		result = append(result, dtos.CurrencyDto{
			Name:         currency.Name,
			Symbol:       currency.Symbol,
			LogoImageURL: currency.LogoImageURL,
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
