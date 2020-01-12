package main

import (
	"net"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	development "github.com/leonardochaia/vendopunkto/plugins/development/internal"
	"github.com/spf13/viper"
)

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "development-plugin",
		Output: os.Stdout,
		Level:  hclog.LevelFromString(strings.ToUpper(viper.GetString("logger.level"))),
	})

	err := run(logger)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func run(logger hclog.Logger) error {
	container, err := development.NewContainer(logger)

	if err != nil {
		return err
	}

	wallets := []plugin.WalletPlugin{
		container.FakeMoneroWallet,
		container.FakeBitcoinWallet,
		container.FakeBitcoinCashWallet,
	}

	for _, wallet := range wallets {
		err = container.Server.AddPlugin(plugin.BuildWalletPlugin(wallet))
		if err != nil {
			return err
		}
	}

	container.Server.AddPlugin(plugin.BuildExchangeRatesPlugin(container.FakeExchangeRatesPlugin))

	container.Server.AddPlugin(plugin.BuildCurrencyMetadataPlugin(container.FakeCurrencyMetadataPlugion))

	addr := net.JoinHostPort(
		viper.GetString("plugin.server.host"),
		viper.GetString("plugin.server.port"))

	logger.Info("Starting plugin", "address", addr)

	err = container.Server.Start(addr)

	if err != nil {
		return err
	}

	return nil
}
