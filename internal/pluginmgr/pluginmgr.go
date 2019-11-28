package pluginmgr

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/internal/pluginwallet"
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/spf13/viper"
)

type walletAndInfo struct {
	client plugin.WalletPlugin
	info   plugin.WalletPluginInfo
}

type Manager struct {
	logger hclog.Logger

	walletRouter *pluginwallet.Router

	wallets map[string]walletAndInfo
}

func (manager *Manager) LoadPlugins() {
	plugins := viper.GetStringSlice("plugins.enabled")

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

		err = manager.initializeWalletPlugin(*url)
		if err != nil {
			manager.logger.Error("Failed to communicate with plugin", "error", err, "URL", url.String())
			continue
		}
	}
}

func (manager *Manager) initializeWalletPlugin(url url.URL) error {

	client := NewWalletClient(url)
	info, err := client.GetPluginInfo()

	if err != nil {
		return err
	}

	manager.logger.Info("Registering wallet plugin", "name", info.Name, "ID", info.ID)
	manager.wallets[info.ID] = walletAndInfo{
		client: client,
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
