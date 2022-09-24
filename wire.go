//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/spf13/cobra"
	"github.com/tomato-net/vault-agent/cli"
	"github.com/tomato-net/vault-agent/config"
	"github.com/tomato-net/vault-agent/logger"
	"github.com/tomato-net/vault-agent/renewer"
)

func ProvideCLI() (*cobra.Command, error) {
	wire.Build(
		logger.New,
		config.New,
		renewer.New,
		cli.New,
	)
	return nil, nil
}
