package plugin

import (
	"github.com/go-chi/chi"
)

type CurrencyMetadata struct {
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	LogoImageURL string `json:"logoImageUrl"`
}

// CurrencyMetadataPlugin providers currency information
type CurrencyMetadataPlugin interface {
	VendoPunktoPlugin
	GetCurrencies(currencies []string) ([]CurrencyMetadata, error)
}

// currencyMetadataServerPlugin mounts the router and provides the actual plugin
// implementation to the handler.
type currencyMetadataServerPlugin struct {
	Impl CurrencyMetadataPlugin
}

func BuildCurrencyMetadataPlugin(impl CurrencyMetadataPlugin) ServerPlugin {
	return &currencyMetadataServerPlugin{
		Impl: impl,
	}
}

func (serverPlugin *currencyMetadataServerPlugin) initializeRouter(router *chi.Mux) error {
	handler := NewCurrencyMetadataHandler(serverPlugin.Impl, serverPlugin)

	router.Mount(CurrencyMetadataMainEndpoint, handler)
	return nil
}

func (serverPlugin *currencyMetadataServerPlugin) GetPluginImpl() (VendoPunktoPlugin, error) {
	return serverPlugin.Impl, nil
}

const (
	CurrencyMetadataMainEndpoint = "/vp/currency-metadata"
)
