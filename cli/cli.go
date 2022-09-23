package cli

import (
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/tomato-net/vault-agent/config"
)

func New(cfg config.Config, log logr.Logger) *cobra.Command {
	return &cobra.Command{
		Use: "vault-agent",
		RunE: func(c *cobra.Command, args []string) error {
			log.Info("got config", "vault_url", cfg.VaultAddress())
			return nil
		},
	}
}
