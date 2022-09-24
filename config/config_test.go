package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("gets string", func(t *testing.T) {
		_getter = func(k string) interface{} {
			return k
		}

		got := get[string]("my_key")
		assert.Equal(t, "my_key", got)
	})

	t.Run("returns default if cant cast", func(t *testing.T) {
		_getter = func(k string) interface{} {
			return map[string]string{k: "value"}
		}

		got := get[string]("hello there")
		assert.Equal(t, "", got)
	})
}

func TestConfig_VaultAddress(t *testing.T) {
	_getter = func(k string) interface{} {
		if k == string(KeyVaultAddress) {
			return "https://test.vault.address"
		}

		return ""
	}

	t.Run("returns vault address", func(t *testing.T) {
		subject := &config{}
		got := subject.VaultAddress()
		assert.Equal(t, "https://test.vault.address", got)
	})
}

func TestConfig_Username(t *testing.T) {
	_getter = func(k string) interface{} {
		if k == string(KeyUsername) {
			return "test-username"
		}

		return ""
	}

	t.Run("returns vault address", func(t *testing.T) {
		subject := &config{}
		got := subject.Username()
		assert.Equal(t, "test-username", got)
	})
}

func TestConfig_TokenFile(t *testing.T) {
	_getter = func(k string) interface{} {
		if k == string(KeyTokenFile) {
			return "~/.test-file"
		}

		return ""
	}

	t.Run("returns vault address", func(t *testing.T) {
		subject := &config{}
		got := subject.TokenFile()
		assert.Equal(t, "~/.test-file", got)
	})
}
