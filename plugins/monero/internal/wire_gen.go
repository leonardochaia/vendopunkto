// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package monero

import (
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/plugin"
)

// Injectors from wire.go:

func NewContainer(globalLogger hclog.Logger) (*Container, error) {
	client, err := newMoneroWalletClient(globalLogger)
	if err != nil {
		return nil, err
	}
	walletPlugin, err := newMoneroWalletPlugin(globalLogger, client)
	if err != nil {
		return nil, err
	}
	server := plugin.NewServer(globalLogger)
	container := newContainer(walletPlugin, server, client)
	return container, nil
}
