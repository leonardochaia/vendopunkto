package server

import (
	"net/http"
	"strings"

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
	router.Post("/pricing/supported", errors.WrapHandler(handler.getSupportedPricingCurrencies))
	router.Get("/payment-methods", errors.WrapHandler(handler.getPaymentMethodCurrencies))
	router.Post("/rates/convert", errors.WrapHandler(handler.getExchange))

	return router
}

func (handler *CurrencyHandler) getPricingCurrencies(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.currencies.getPricingCurrencies"

	pricingCurrencies := handler.runtimeConf.GetPricingCurrencies()
	if len(pricingCurrencies) == 0 {
		render.JSON(w, r, []dtos.CurrencyDto{})
		return nil
	}

	currencies, err := handler.currencyRepo.FindBySymbols(r.Context(), pricingCurrencies)

	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	result := handler.convertToCurrencyDto(pricingCurrencies, currencies)

	render.JSON(w, r, result)
	return nil
}

func (handler *CurrencyHandler) getSupportedPricingCurrencies(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.currencies.getSupportedPricingCurrencies"

	var params = new(dtos.SearchSupportedCurrenciesParams)

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		return errors.E(op, errors.Parameters, err)
	}

	exchangeRatesPlugin, err := handler.plugins.GetConfiguredExchangeRatesPlugin()
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	supportedCurrencies, err := exchangeRatesPlugin.SearchSupportedCurrencies(params.Term)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	if len(supportedCurrencies) == 0 {
		render.JSON(w, r, []dtos.CurrencyDto{})
		return nil
	}

	symbols := make([]string, len(supportedCurrencies))
	i := 0
	for _, c := range supportedCurrencies {
		symbols[i] = c.Symbol
		i++
	}

	currencies, err := handler.currencyRepo.FindBySymbols(r.Context(), symbols)

	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	result := handler.basicConvertToCurrencyDto(supportedCurrencies, currencies)
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

func (handler *CurrencyHandler) convertToCurrencyDto(
	symbols []string,
	currencies []*vendopunkto.Currency) []dtos.CurrencyDto {
	result := []dtos.CurrencyDto{}
	for _, symbol := range symbols {
		var currency *vendopunkto.Currency
		for _, c := range currencies {
			if strings.ToLower(c.Symbol) == strings.ToLower(symbol) {
				currency = c
				break
			}
		}
		_, err := handler.plugins.GetWalletInfoForCurrency(symbol)
		if currency != nil {
			result = append(result, dtos.CurrencyDto{
				Name:             currency.Name,
				Symbol:           currency.Symbol,
				LogoImageURL:     currency.LogoImageURL,
				SupportsPayments: err != nil,
			})
		} else {
			result = append(result, dtos.CurrencyDto{
				Name:             strings.ToUpper(symbol),
				Symbol:           symbol,
				SupportsPayments: err != nil,
			})
		}
	}
	return result
}

func (handler *CurrencyHandler) basicConvertToCurrencyDto(
	basics []dtos.BasicCurrencyDto,
	currencies []*vendopunkto.Currency) []dtos.CurrencyDto {
	result := []dtos.CurrencyDto{}
	for _, basic := range basics {
		var currency *vendopunkto.Currency
		for _, c := range currencies {
			if strings.ToLower(c.Symbol) == strings.ToLower(basic.Symbol) {
				currency = c
				break
			}
		}
		_, err := handler.plugins.GetWalletInfoForCurrency(basic.Symbol)
		if currency != nil {
			result = append(result, dtos.CurrencyDto{
				Name:             currency.Name,
				Symbol:           currency.Symbol,
				LogoImageURL:     currency.LogoImageURL,
				SupportsPayments: err != nil,
			})
		} else {
			name := basic.Name
			if name == "" {
				name = strings.ToUpper(basic.Symbol)
			}
			result = append(result, dtos.CurrencyDto{
				Name:             name,
				Symbol:           basic.Symbol,
				LogoImageURL:     basic.LogoImageURL,
				SupportsPayments: err != nil,
			})
		}
	}
	return result
}
