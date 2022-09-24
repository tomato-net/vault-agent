package token

import (
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/hashicorp/vault/api"
	"github.com/tomato-net/vault-agent/config"
)

// Renewer handles the renewal of a connection to vault
type Renewer struct {
	client   *api.Client
	config   config.Config
	watcher  *api.LifetimeWatcher
	logger   logr.Logger
	provider Provider
}

func NewRenewer(client *api.Client, logger logr.Logger, provider Provider, config config.Config) *Renewer {
	return &Renewer{
		client:   client,
		config:   config,
		logger:   logger,
		provider: provider,
	}
}

func (r *Renewer) Start() error {
	token, err := r.provider.Token()
	if err != nil {
		return fmt.Errorf("provider token: %w", err)
	}

	r.client.SetToken(token.Token)
	// TODO: Is this necessary?
	secret, err := r.client.Auth().Token().Create(&api.TokenCreateRequest{TTL: "100"})
	if err != nil {
		return fmt.Errorf("create token: %w", err)
	}

	watcher, err := r.client.NewLifetimeWatcher(&api.LifetimeWatcherInput{Secret: secret})
	if err != nil {
		return fmt.Errorf("new watcher: %w", err)
	}

	r.watcher = watcher
	go r.watcher.Start()

	writeToken := func(token string) error {
		f, err := os.OpenFile(r.config.TokenFile(), os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return fmt.Errorf("opening token file: %w", err)
		}

		defer f.Close()

		if _, err := f.Write([]byte(token)); err != nil {
			return fmt.Errorf("writing token: %w", err)
		}

		r.logger.Info("wrote token to file", "token", token)
		return nil
	}

	if err := writeToken(secret.Auth.ClientToken); err != nil {
		return fmt.Errorf("writing login token: %w", err)
	}

	go time.AfterFunc(5*time.Second, func() {
		for {
			if _, err := r.client.Auth().Token().LookupSelf(); err != nil {
				r.logger.Error(err, "failed looking up self")
			} else {
				r.logger.Info("successfully looked up self")
			}

			time.Sleep(10 * time.Second)
		}
	})

	for {
		select {
		case err := <-r.watcher.DoneCh():
			if err != nil {
				return fmt.Errorf("renewer done: %w", err)
			}

			return nil
		case renewal := <-r.watcher.RenewCh():
			if err := writeToken(renewal.Secret.Auth.ClientToken); err != nil {
				return fmt.Errorf("writing renewed token: %w", err)
			}
		}
	}
}

func (r *Renewer) Stop() error {
	r.watcher.Stop()
	return nil
}