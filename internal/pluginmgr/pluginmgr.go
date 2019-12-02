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

type Manager struct {
	logger hclog.Logger
	http   http.Client

	wallets map[string]walletAndInfo
}

func (manager *Manager) LoadPlugins() {
	plugins := viper.GetStringSlice("plugins.enabled")
	hostAddress := viper.GetString("plugins.server.plugin_host_address")

	for _, addr := range plugins {
		split := strings.Split(addr, "|")

		if split[0] != "wallet" {
			manager.logger.Error("Failed to load plugin. Only 'wallet' type is supported",
				"plugin", addr)
			manager.logger.Error("An example could be 'wallet|http://localhost:3333'")
			continue
		}

		url, err := url.Parse(split[1])

		if err != nil {
			manager.logger.Error("Failed to parse plugin URL", "error", err)
			manager.logger.Error("An example could be 'wallet|http://localhost:3333'")
			continue
		}

		err = manager.initializeWalletPlugin(*url, hostAddress)
		if err != nil {
			manager.logger.Error("Failed to communicate with plugin", "error", err, "URL", url.String())
			continue
		}
	}
}

func (manager *Manager) initializeWalletPlugin(pluginURL url.URL, hostAddress string) error {

	info, err := manager.activatePlugin(pluginURL, hostAddress)

	if err != nil {
		return err
	}

	manager.logger.Info("Registering wallet plugin", "name", info.Name, "ID", info.ID)
	manager.wallets[info.ID] = walletAndInfo{
		client: NewWalletClient(pluginURL, *info),
		info:   *info,
	}
	return nil
}

func (manager *Manager) GetWallet(ID string) (plugin.WalletPlugin, error) {
	if w, ok := manager.wallets[ID]; ok {
		return w.client, nil
	}
	return nil, fmt.Errorf("Could not find a wallet with ID " + ID)
}

func (manager *Manager) GetWalletForCurrency(currency string) (plugin.WalletPlugin, error) {
	for _, wallet := range manager.wallets {
		for _, c := range wallet.info.Currencies {
			if c == currency {
				return wallet.client, nil
			}
		}
	}
	return nil, fmt.Errorf("Could not find a wallet for currency " + currency)
}

// activatePlugin does the initial handshake where the plugin returns its basic
// info while the host address is provided so the plugin can reach back
func (manager *Manager) activatePlugin(apiURL url.URL, hostAddress string) (*plugin.WalletPluginInfo, error) {
	u, err := url.Parse(plugin.WalletMainEndpoint + plugin.ActivatePluginEndpoint)
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

	resp, err := manager.http.Post(final.String(), "application/json", bytes.NewBuffer(params))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid response. Got status " + resp.Status)
	}

	var result plugin.WalletPluginInfo
	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}
