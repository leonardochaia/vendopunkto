package plugin

import (
	"github.com/go-chi/chi"
)

type WalletPluginInfo struct {
	ID         string
	Name       string
	Currencies []string
}

// WalletPlugin must be implemented for a Coin to be supported by vendopunkto
type WalletPlugin interface {
	GetPluginInfo() (*WalletPluginInfo, error)
	GenerateNewAddress(invoiceID string) (string, error)
}

// walletServerPlugin mounts the router and provides the actual plugin
// implementation to the handler.
type walletServerPlugin struct {
	Impl WalletPlugin
}

func BuildWalletPlugin(impl WalletPlugin) ServerPlugin {
	return walletServerPlugin{
		Impl: impl,
	}
}

func (builder walletServerPlugin) initializeRouter(router *chi.Mux) {
	handler := NewWalletHandler(builder.Impl)

	router.Mount(WalletMainEndpoint, handler)
}

const (
	WalletMainEndpoint            = "/vp/wallet"
	GenerateAddressWalletEndpoint = "/address"
	PluginInfoEndpoint            = "/info"
)
