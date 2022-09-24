package token

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/ldap"
	"github.com/tomato-net/vault-agent/config"
)

var _ Provider = (*ProviderLDAP)(nil)

type ProviderLDAP struct {
	client *api.Client
	logger logr.Logger
	config config.Config
}

func NewProviderLDAP(client *api.Client, logger logr.Logger, config config.Config) *ProviderLDAP {
	return &ProviderLDAP{
		client: client,
		logger: logger,
		config: config,
	}
}

func (l *ProviderLDAP) Token() (Token, error) {
	ldapAuth, err := ldap.NewLDAPAuth(l.config.Username(), &ldap.Password{FromString: l.config.Password()})
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
