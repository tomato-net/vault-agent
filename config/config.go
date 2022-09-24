package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	KeyVaultAddress key = "vault_address"
	KeyUsername         = "username"
	KeyPassword         = "password"
	KeyTokenFile        = "token_file"
)

type key string

func (k key) String() string {
	return string(k)
}

type Getter func(string) interface{}

var _getter Getter = viper.Get

type Config interface {
	VaultAddress() string
	Username() string
	Password() string
	TokenFile() string
}

type config struct{}

func New() (Config, error) {
	viper.SetConfigName(".varc")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.vault-agent")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("VAULT_AGENT")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading in config: %w", err)
	}

	return &config{}, nil
}

func (c *config) VaultAddress() string {
	return get[string](KeyVaultAddress)
}

func (c *config) Username() string {
	return get[string](KeyUsername)
}

func (c *config) Password() string {
	return get[string](KeyPassword)
}

func (c *config) TokenFile() string {
	return get[string](KeyTokenFile)
}

func get[T any](k key) T {
	untyped := _getter(k.String())
	if typed, ok := untyped.(T); ok {
		return typed
	}

	return *new(T)
}
