package token

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/tomato-net/vault-agent/config"
)

func NewClient(config config.Config) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: config.VaultAddress(),
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("vault new client: %w", err)
	}

	return client, nil
}
