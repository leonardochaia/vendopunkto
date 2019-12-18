package pluginmgr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/spf13/viper"
)

type walletAndInfo struct {
	client plugin.WalletPlugin
	info   plugin.WalletPluginInfo
}

type exchangeRatesAndInfo struct {
	client plugin.ExchangeRatesPlugin
	info   plugin.PluginInfo
}

type Manager struct {
	logger hclog.Logger
	client http.Client

	wallets       map[string]walletAndInfo
	exchangeRates map[string]exchangeRatesAndInfo
}

func (manager *Manager) LoadPlugins() {
	plugins := viper.GetStringSlice("plugins.enabled")
	hostAddress := viper.GetString("plugins.server.plugin_host_address")

	for _, addr := range plugins {

		url, err := url.Parse(addr)

		if err != nil {
			manager.logger.Error("Failed to parse plugin URL", "error", err, "address", addr)
			manager.logger.Error("An example could be 'http://localhost:3333'")
			continue
		}

		err = manager.initializePlugin(*url, hostAddress)
		if err != nil {
			manager.logger.Error("Failed to communicate with plugin", "error", err, "URL", url.String())
			continue
		}
	}
}

func (manager *Manager) GetWallet(ID string) (plugin.WalletPlugin, error) {
	if w, ok := manager.wallets[ID]; ok {
		return w.client, nil
	}
	return nil, fmt.Errorf("Could not find a wallet with ID " + ID)
}

func (manager *Manager) GetWalletForCurrency(currency string) (plugin.WalletPlugin, error) {
	for _, wallet := range manager.wallets {
		if wallet.info.Currency.Symbol == currency {
			return wallet.client, nil
		}
	}
	return nil, fmt.Errorf("Could not find a wallet for currency " + currency)
}

func (manager *Manager) GetWalletInfoForCurrency(currency string) (plugin.WalletPluginInfo, error) {
	for _, wallet := range manager.wallets {
		if wallet.info.Currency.Symbol == currency {
			return wallet.info, nil
		}
	}
	return plugin.WalletPluginInfo{}, fmt.Errorf("Could not find a wallet for currency " + currency)
}

func (manager *Manager) GetAllCurrencies() ([]plugin.WalletPluginCurrency, error) {
	output := []plugin.WalletPluginCurrency{}
	for _, v := range manager.wallets {
		output = append(output, v.info.Currency)
	}

	return output, nil
}

func (manager *Manager) GetExchangeRatesPlugin(ID string) (plugin.ExchangeRatesPlugin, error) {
	if w, ok := manager.exchangeRates[ID]; ok {
		return w.client, nil
	}
	return nil, fmt.Errorf("Could not find an exchange rates plugin with ID " + ID)
}

func (manager *Manager) GetConfiguredExchangeRatesPlugin() (plugin.ExchangeRatesPlugin, error) {
	return manager.GetExchangeRatesPlugin(viper.GetString("plugins.default_exchange_rates"))
}

func (manager *Manager) initializePlugin(pluginURL url.URL, hostAddress string) error {

	infos, err := manager.activatePlugin(pluginURL, hostAddress)

	if err != nil {
		return err
	}

	for _, info := range infos {

		switch info.Type {
		case plugin.PluginTypeWallet:
			err = manager.initializeWalletPlugin(pluginURL, info)
		case plugin.PluginTypeExchangeRate:
			err = manager.initializeExchangeRatesPlugin(pluginURL, info)
		default:
			err = fmt.Errorf("Plugin type is not supported %s", info.Type)
		}

		if err != nil {
			manager.logger.Error("Failed to initialize plugin",
				"id", info.ID,
				"name", info.Name,
				"type", info.Type,
				"address", pluginURL.String(),
				"error", err)
			continue
		}
	}

	return nil
}

func (manager *Manager) initializeExchangeRatesPlugin(
	pluginURL url.URL,
	info plugin.PluginInfo) error {

	ratesClient := NewExchangeRatesClient(pluginURL, info, manager.client)

	info, err := ratesClient.GetPluginInfo()
	if err != nil {
		return err
	}

	manager.exchangeRates[info.ID] = exchangeRatesAndInfo{
		client: ratesClient,
		info:   info,
	}

	manager.logger.Info("Initialized Exchange rates Plugin",
		"id", info.ID,
		"name", info.Name,
		"address", pluginURL.String()+info.GetAddress())

	return nil
}

func (manager *Manager) initializeWalletPlugin(pluginURL url.URL, info plugin.PluginInfo) error {
	walletClient := NewWalletClient(pluginURL, info, manager.client)

	walletInfo, err := walletClient.GetWalletInfo()
	if err != nil {
		return err
	}

	walletInfo.Currency.Symbol = strings.ToLower(walletInfo.Currency.Symbol)

	manager.wallets[info.ID] = walletAndInfo{
		client: walletClient,
		info:   walletInfo,
	}

	manager.logger.Info("Initialized Wallet Plugin",
		"id", info.ID,
		"name", info.Name,
		"currency", walletInfo.Currency.Symbol,
		"address", pluginURL.String()+info.GetAddress())

	return nil
}

// activatePlugin does the initial handshake where the plugin returns its basic
// info while the host address is provided so the plugin can reach back
func (manager *Manager) activatePlugin(apiURL url.URL, hostAddress string) ([]plugin.PluginInfo, error) {
	u, err := url.Parse(plugin.ActivatePluginEndpoint)
	if err != nil {
		return nil, err
	}

	final := apiURL.ResolveReference(u)

	if err != nil {
		return nil, err
	}

	params, err := json.Marshal(&plugin.ActivatePluginParams{
		HostAddress: hostAddress,
	})

	if err != nil {
		return nil, err
	}

	resp, err := manager.client.Post(final.String(), "application/json", bytes.NewBuffer(params))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid response. Got status " + resp.Status)
	}

	var result []plugin.PluginInfo
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
