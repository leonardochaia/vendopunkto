package plugin

import (
	"github.com/go-chi/chi"
)

type ExchangeRatesResult map[string]float64

// ExchangeRatesPlugin must be implemented in order to provide coin exchange
// rates
type ExchangeRatesPlugin interface {
	VendoPunktoPlugin
	GetExchangeRates(base string, currencies []string) (ExchangeRatesResult, error)
}

// exchangeRatesServerPlugin mounts the router and provides the actual plugin
// implementation to the handler.
type exchangeRatesServerPlugin struct {
	Impl ExchangeRatesPlugin
}

func BuildExchangeRatesPlugin(impl ExchangeRatesPlugin) ServerPlugin {
	return &exchangeRatesServerPlugin{
		Impl: impl,
	}
}

func (serverPlugin *exchangeRatesServerPlugin) initializeRouter(router *chi.Mux) error {
	handler := NewHandler(serverPlugin.Impl, serverPlugin)

	router.Mount(ExchangeRatesMainEndpoint, handler)
	return nil
}

func (serverPlugin *exchangeRatesServerPlugin) GetPluginImpl() (VendoPunktoPlugin, error) {
	return serverPlugin.Impl, nil
}

const (
	ExchangeRatesMainEndpoint = "/vp/exchange-rates"
)
