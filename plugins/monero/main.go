package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
	monero "github.com/leonardochaia/vendopunkto/plugins/monero/internal"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
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

	// This is the server that listen for monero-wallet-rpc tx-notify
	err = startMoneroHTTPServer(
		logger,
		container.WalletClient,
		container.Handler.Routes())

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

func startMoneroHTTPServer(
	logger hclog.Logger,
	client wallet.Client,
	handler http.Handler) error {

	server := &http.Server{
		Addr:    net.JoinHostPort(viper.GetString("monero.server.host"), viper.GetString("monero.server.port")),
		Handler: handler,
	}

	// Listen
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("Could not listen on %s: %v", server.Addr, err)
	}

	go func() {
		if err = server.Serve(listener); err != nil {
			logger.Error("API Listen error", "error", err, "address", server.Addr)
			os.Exit(1)
		}
	}()
	logger.Info("Plugin Server Listening", "address", server.Addr)

	return nil
}
