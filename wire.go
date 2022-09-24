//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/spf13/cobra"
	"github.com/tomato-net/vault-agent/cli"
	"github.com/tomato-net/vault-agent/config"
	"github.com/tomato-net/vault-agent/credentials"
	"github.com/tomato-net/vault-agent/logger"
	"github.com/tomato-net/vault-agent/token"
)

func ProvideCLI() (*cobra.Command, error) {
	wire.Build(
		logger.New,
		config.New,
		credentials.NewKeychainAccessor,
		wire.Bind(new(credentials.Reader), new(*credentials.KeychainAccessor)),
		wire.Bind(new(credentials.ReadWriter), new(*credentials.KeychainAccessor)),
		token.NewClient,
		token.NewProviderLDAP,
		wire.Bind(new(token.Provider), new(*token.ProviderLDAP)),
		token.NewRenewer,
		cli.New,
	)
	return nil, nil
}
