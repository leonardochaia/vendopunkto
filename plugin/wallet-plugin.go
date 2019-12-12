package plugin

import (
	"github.com/go-chi/chi"
)

type WalletPluginCurrency struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type WalletPluginInfo struct {
	Currency WalletPluginCurrency `json:"currency"`
}

// WalletPlugin must be implemented for a Coin to be supported by vendopunkto
type WalletPlugin interface {
	VendoPunktoPlugin
	GetWalletInfo() (WalletPluginInfo, error)
	GenerateNewAddress(invoiceID string) (string, error)
}

// walletServerPlugin mounts the router and provides the actual plugin
// implementation to the handler.
type walletServerPlugin struct {
	Impl WalletPlugin
}

func BuildWalletPlugin(impl WalletPlugin) ServerPlugin {
	return &walletServerPlugin{
		Impl: impl,
	}
}

func (serverPlugin *walletServerPlugin) initializeRouter(router *chi.Mux) error {
	handler := NewWalletHandler(serverPlugin.Impl, serverPlugin)

	router.Mount(WalletMainEndpoint, handler)
	return nil
}

func (serverPlugin *walletServerPlugin) GetPluginImpl() (VendoPunktoPlugin, error) {
	return serverPlugin.Impl, nil
}

const (
	WalletMainEndpoint            = "/vp/wallet"
	GenerateAddressWalletEndpoint = "/address"
	WalletInfoEndpoint            = "/info"
)
