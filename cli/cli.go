package cli

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/tomato-net/vault-agent/token"
)

func New(r *token.Renewer, log logr.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use: "vault-agent",
		RunE: func(c *cobra.Command, args []string) error {
			log.Info("starting renewer")
			defer r.Stop()
			if err := r.Start(); err != nil {
				log.Error(err, "renewer failed")
				os.Exit(1)
			}

			log.Info("finishing")
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
