package cli

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/tomato-net/vault-agent/token"
)

func New(r *token.Renewer, log logr.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use: "vault-agent",
		RunE: func(c *cobra.Command, args []string) error {
			log.Info("starting renewer")

			start := func() error {
				if err := r.Start(); err != nil {
					return fmt.Errorf("renewer failed: %w", err)
				}

				defer r.Stop()

				return nil
			}

			for {
				log.Info("generating new token renewer")

				if err := start(); err != nil {
					log.Error(err, "token renewer failed")
				}
			}

			log.Info("vault-agent shutting down")
			return nil
		},
	}

	cmd.AddCommand(NewStatus(log))

	return cmd
}

func NewStatus(log logr.Logger) *cobra.Command {
	return &cobra.Command{
		Use: "status",
		RunE: func(c *cobra.Command, args []string) error {
			log.Info("running!")
			return nil
		},
	}
}
