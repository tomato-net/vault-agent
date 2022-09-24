package token

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/ldap"
	"github.com/tomato-net/vault-agent/config"
	"github.com/tomato-net/vault-agent/credentials"
)

var _ Provider = (*ProviderLDAP)(nil)

type ProviderLDAP struct {
	client *api.Client
	logger logr.Logger
	config config.Config
	creds  credentials.Reader
}

func NewProviderLDAP(client *api.Client, logger logr.Logger, config config.Config, creds credentials.Reader) *ProviderLDAP {
	return &ProviderLDAP{
		client: client,
		logger: logger.WithName("provider_ldap"),
		config: config,
		creds:  creds,
	}
}

func (l *ProviderLDAP) Token() (Token, error) {
	password, err := l.creds.Read()
	if err != nil {
		return Token{}, fmt.Errorf("credentials provider: %w", err)
	}

	ldapAuth, err := ldap.NewLDAPAuth(l.config.Username(), &ldap.Password{FromString: password.Data})
	if err != nil {
		return Token{}, fmt.Errorf("new ldap auth: %w", err)
	}

	// TODO: Use different error to stop retries to prevent lockout
	// TODO: Use context
	secret, err := l.client.Auth().Login(context.Background(), ldapAuth)
	if err != nil {
		return Token{}, fmt.Errorf("ldap login: %w", err)
	}

	return Token{
		Token: secret.Auth.ClientToken,
	}, nil
}
