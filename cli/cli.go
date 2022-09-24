package cli

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/tomato-net/vault-agent/credentials"
	"github.com/tomato-net/vault-agent/token"
)

func New(log logr.Logger, renewer *token.Renewer, creds credentials.ReadWriter) *cobra.Command {
	cmd := &cobra.Command{
		Use: "vault-agent",
		RunE: func(c *cobra.Command, args []string) error {
			log.Info("starting renewer")
			if has, err := creds.Has(); !has || err != nil {
				log.Error(fmt.Errorf("no creds found: %w", err), "no creds found in provider")
				os.Exit(1)
			}

			start := func() error {
				if err := renewer.Start(); err != nil {
					return fmt.Errorf("renewer failed: %w", err)
				}

				defer renewer.Stop()

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
	cmd.AddCommand(NewCredentials(log, creds))

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

func NewCredentials(log logr.Logger, creds credentials.ReadWriter) *cobra.Command {
	cmd := &cobra.Command{
		Use: "credentials",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:           "add",
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(c *cobra.Command, args []string) error {
				var password string
				fmt.Print("password: ")
				fmt.Scanf("%s", &password)
				if err := creds.Write(credentials.PasswordCredential{Data: password}); err != nil {
					log.Error(err, "failed to add credentials")
					os.Exit(1)
				}

				return nil
			},
		},
		&cobra.Command{
			Use:           "update",
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(c *cobra.Command, args []string) error {
				var password string
				fmt.Print("password: ")
				fmt.Scanf("%s", &password)

				if err := creds.Delete(); err != nil {
					log.Error(err, "failed deleting existing credentials to update, trying to write new creds anyway")
				}

				if err := creds.Write(credentials.PasswordCredential{Data: password}); err != nil {
					log.Error(err, "failed to add credentials, you may have no creds anymore, that sucks")
					os.Exit(1)
				}

				return nil
			},
		},
	)

	return cmd
}
