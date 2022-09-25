package token

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/userpass"
	"github.com/tomato-net/vault-agent/config"
	"github.com/tomato-net/vault-agent/credentials"
)

var _ Provider = (*ProviderUserPass)(nil)

type ProviderUserPass struct {
	client *api.Client
	logger logr.Logger
	config config.Config
	creds  credentials.Reader
}

func NewProviderUserPass(client *api.Client, logger logr.Logger, config config.Config, creds credentials.Reader) *ProviderUserPass {
	return &ProviderUserPass{
		client: client,
		logger: logger.WithName("provider_userpass"),
		config: config,
		creds:  creds,
	}
}

func (l *ProviderUserPass) Token() (Token, error) {
	password, err := l.creds.Read()
	if err != nil {
		return Token{}, fmt.Errorf("credentials provider: %w", err)
	}

	userpassAuth, err := userpass.NewUserpassAuth(l.config.Username(), &userpass.Password{FromString: password.Data})
	if err != nil {
		return Token{}, fmt.Errorf("new userpass auth: %w", err)
	}

	// TODO: Use different error to stop retries to prevent lockout
	// TODO: Use context
	secret, err := l.client.Auth().Login(context.Background(), userpassAuth)
	if err != nil {
		return Token{}, fmt.Errorf("userpass login: %w", err)
	}

	return Token{
		Token: secret.Auth.ClientToken,
	}, nil
}
