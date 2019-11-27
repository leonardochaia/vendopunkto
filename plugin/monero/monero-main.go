package main

import (
	"errors"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
	"github.com/spf13/viper"
)

type MoneroWalletPlugin struct {
	client wallet.Client
}

func (p MoneroWalletPlugin) GenerateNewAddress(invoiceID string) (string, error) {
	result, err := p.client.MakeIntegratedAddress(&wallet.RequestMakeIntegratedAddress{})

	if err != nil {
		return "", err
	}

	return result.IntegratedAddress, nil
}

func (p MoneroWalletPlugin) GetPluginInfo() (*plugin.WalletPluginInfo, error) {
	return &plugin.WalletPluginInfo{
		Name:       "Monero Wallet",
		ID:         "monero-wallet",
		Currencies: []string{"XMR"},
	}, nil
}

func main() {

	// Sets up the config file, environment etc
	viper.SetTypeByDefaultValue(true)                      // If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	viper.AutomaticEnv()                                   // Automatically use environment variables where available
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Environement variables use underscores instead of periods

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "monero-plugin",
		Output: os.Stdout,
		Level:  hclog.LevelFromString(strings.ToUpper(viper.GetString("logger.level"))),
	})

	client, err := newClient(logger)

	if err != nil {
		panic(err.Error())
	}

	server := plugin.NewServer()

	monero := MoneroWalletPlugin{
		client: client,
	}

	err = server.AddPlugin(plugin.BuildWalletPlugin(monero))
	if err != nil {
		panic(err.Error())
	}

	err = server.Start("0.0.0.0:3333")

	if err != nil {
		panic(err.Error())
	}
}

func newClient(globalLogger hclog.Logger) (wallet.Client, error) {
	baseURL := viper.GetString("monero.wallet_rpc_url")
	logger := globalLogger.Named("monero")

	if baseURL == "" {
		return nil, errors.New("A Monero Wallet RPC url must be provided")
	}

	client := wallet.New(wallet.Config{
		Address: baseURL + "/json_rpc",
	})

	version, err := client.GetVersion()

	if err != nil {
		logger.Warn("Failed to connect to Monero Wallet RPC", "url", baseURL)
		return nil, err
	}

	logger.Info("Monero Wallet RPC Test Success", "version", version.Version)

	return client, nil
}
