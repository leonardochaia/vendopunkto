package plugin

import (
	"fmt"
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
	Impl     WalletPlugin
	VPClient PluginWalletClient
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

func (serverPlugin *walletServerPlugin) initializePlugin(hostAddr string) error {
	client, err := NewWalletClient(hostAddr)
	if err != nil {
		return err
	}

	serverPlugin.VPClient = client
	return nil
}

func (serverPlugin *walletServerPlugin) GetWalletClient() (PluginWalletClient, error) {
	if serverPlugin.VPClient != nil {
		return serverPlugin.VPClient, nil
	}

	return nil, fmt.Errorf("Plugin has not been initialized yet")
}

const (
	WalletMainEndpoint            = "/vp/wallet"
	GenerateAddressWalletEndpoint = "/address"
	WalletInfoEndpoint            = "/info"
)
