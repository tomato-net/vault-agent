//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/spf13/cobra"
	"github.com/tomato-net/vault-agent/cli"
	"github.com/tomato-net/vault-agent/config"
	"github.com/tomato-net/vault-agent/logger"
	"github.com/tomato-net/vault-agent/token"
)

func ProvideCLI() (*cobra.Command, error) {
	wire.Build(
		logger.New,
		config.New,
		token.NewClient,
		token.NewProviderLDAP,
		wire.Bind(new(token.Provider), new(*token.ProviderLDAP)),
		token.NewRenewer,
		cli.New,
	)
	return nil, nil
}
