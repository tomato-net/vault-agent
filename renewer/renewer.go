package renewer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/ldap"
	"github.com/tomato-net/vault-agent/config"
)

// Renewer handles the renewal of a connection to vault
type Renewer struct {
	vault *api.Client
	c     config.Config
	w     *api.LifetimeWatcher
	l     logr.Logger
}

func New(c config.Config, l logr.Logger) (*Renewer, error) {
	client, err := api.NewClient(&api.Config{
		Address: c.VaultAddress(),
		Timeout: 10 * time.Second,
	})

	if err != nil {
		return nil, fmt.Errorf("vault new client: %w", err)
	}

	return &Renewer{vault: client, l: l, c: c}, nil
}

func (r *Renewer) Start() error {
	ldapAuth, err := ldap.NewLDAPAuth(r.c.Username(), &ldap.Password{FromString: r.c.Password()})
	if err != nil {
		return fmt.Errorf("new ldap auth: %w", err)
	}

	secret, err := r.vault.Auth().Login(context.Background(), ldapAuth)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	r.l.Info("logged in", "token", secret.Auth.ClientToken)

	watcher, err := r.vault.NewLifetimeWatcher(&api.LifetimeWatcherInput{Secret: secret})
	if err != nil {
		return fmt.Errorf("new watcher: %w", err)
	}

	r.w = watcher
	go r.w.Start()

	writeToken := func(token string) error {
		f, err := os.OpenFile(r.c.TokenFile(), os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return fmt.Errorf("opening token file: %w", err)
		}

		defer f.Close()

		if _, err := f.Write([]byte(token)); err != nil {
			return fmt.Errorf("writing token: %w", err)
		}

		r.l.Info("wrote token to file", "token", token)
		return nil
	}

	if err := writeToken(secret.Auth.ClientToken); err != nil {
		return fmt.Errorf("writing login token: %w", err)
	}

	for {
		select {
		case err := <-r.w.DoneCh():
			if err != nil {
				return fmt.Errorf("renewer done: %w", err)
			}

			return nil
		case renewal := <-r.w.RenewCh():
			if err := writeToken(renewal.Secret.Auth.ClientToken); err != nil {
				return fmt.Errorf("writing renewed token: %w", err)
			}
		}
	}
}

func (r *Renewer) Stop() error {
	r.w.Stop()
	return nil
}
