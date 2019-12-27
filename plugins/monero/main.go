package main

import (
	"net"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	monero "github.com/leonardochaia/vendopunkto/plugins/monero/internal"
	"github.com/spf13/viper"
)

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "monero-plugin",
		Output: os.Stdout,
		Level:  hclog.LevelFromString(strings.ToUpper(viper.GetString("logger.level"))),
	})

	container, err := monero.NewContainer(logger)

	if err != nil {
		panic(err.Error())
	}

	err = container.Server.AddPlugin(plugin.BuildWalletPlugin(container.Plugin))
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
