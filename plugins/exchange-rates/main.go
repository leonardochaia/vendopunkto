package main

import (
	"net"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	rates "github.com/leonardochaia/vendopunkto/plugins/exchange-rates/internal"
	"github.com/spf13/viper"
)

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "exchange-rates-plugin",
		Output: os.Stdout,
		Level:  hclog.LevelFromString(strings.ToUpper(viper.GetString("logger.level"))),
	})

	container, err := rates.NewContainer(logger)

	if err != nil {
		panic(err.Error())
	}

	err = container.Server.AddPlugin(plugin.BuildExchangeRatesPlugin(container.Plugin))
	if err != nil {
		panic(err.Error())
	}

	addr := net.JoinHostPort(
		viper.GetString("plugin.server.host"),
		viper.GetString("plugin.server.port"))

	logger.Info("Starting plugin", "address", addr)

	err = container.Server.Start(addr)

	if err != nil {
		panic(err.Error())
	}
}
