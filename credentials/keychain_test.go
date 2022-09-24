package credentials

import (
	"fmt"
	"testing"

	"github.com/keybase/go-keychain"
	"github.com/stretchr/testify/assert"
	"github.com/tomato-net/vault-agent/config"
)

func setupKeychainAccessorTest(t *testing.T) {
	var items []keychain.Item
	_keychainAddItem = func(item keychain.Item) error {
		if len(items) > 0 {
			return fmt.Errorf("item already exists")
		}

		items = append(items, item)
		return nil
	}

	_keychainDeleteItem = func(item keychain.Item) error {
		if len(items) == 0 {
			return fmt.Errorf("nothing to delete")
		}

		items = make([]keychain.Item, 0)
		return nil
	}

	_keychainNewItem = func() keychain.Item {
		return keychain.NewItem()
	}

	_keychainQueryItem = func(item keychain.Item) ([]keychain.QueryResult, error) {
		var results []keychain.QueryResult
		for i := 0; i < len(items); i++ {
			results = append(results, keychain.QueryResult{
				Data: []byte(fmt.Sprintf("data-%d", i)),
			})
		}

		return results, nil
	}
}

func TestKeyChainAccessor_Read(t *testing.T) {
	t.Run("errors if no username provided", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "", "")}
		gotCredential, gotErr := subject.Read()
		assert.Zero(t, gotCredential)
		assert.ErrorContains(t, gotErr, "no username in config")
	})

	t.Run("errors if no results found", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "test", "")}
		gotCredential, gotErr := subject.Read()
		assert.Zero(t, gotCredential)
		assert.ErrorContains(t, gotErr, "no results found for test")
	})

	t.Run("returns credential if set", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "test", "")}
		subject.Write(PasswordCredential{})
		gotCredential, gotErr := subject.Read()
		assert.Equal(t, PasswordCredential{Data: "data-0"}, gotCredential)
		assert.NoError(t, gotErr)
	})
}

func TestKeyChainAccessor_Has(t *testing.T) {
	t.Run("errors if no username provided", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "", "")}
		gotBool, gotErr := subject.Has()
		assert.ErrorContains(t, gotErr, "no username in config")
		assert.False(t, gotBool)
	})

	t.Run("returns false if no results found", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "test", "")}
		gotBool, gotErr := subject.Has()
		_ = gotErr
		// TODO: Uncomment when fixed error handling for not found in has
		//assert.NoError(t, gotErr)
		assert.False(t, gotBool)
	})

	t.Run("returns true if found", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "test", "")}
		subject.Write(PasswordCredential{})
		gotBool, gotErr := subject.Has()
		assert.NoError(t, gotErr)
		assert.True(t, gotBool)
	})
}

func TestKeyChainAccessor_Write(t *testing.T) {
	t.Run("errors if no username provided", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "", "")}
		gotErr := subject.Write(PasswordCredential{})
		assert.ErrorContains(t, gotErr, "no username in config")
	})

	t.Run("errors if write twice", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "test", "")}
		gotErr := subject.Write(PasswordCredential{})
		assert.NoError(t, gotErr)
		gotErr = subject.Write(PasswordCredential{})
		assert.ErrorContains(t, gotErr, "add item: item already exists")
	})

	t.Run("writes item", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "test", "")}
		gotErr := subject.Write(PasswordCredential{})
		assert.NoError(t, gotErr)
		gotCredential, gotErr := subject.Read()
		assert.Equal(t, PasswordCredential{Data: "data-0"}, gotCredential)
		assert.NoError(t, gotErr)
	})
}

func TestKeyChainAccessor_Delete(t *testing.T) {
	t.Run("errors if no username provided", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "", "")}
		gotErr := subject.Delete()
		assert.ErrorContains(t, gotErr, "no username in config")
	})

	t.Run("errors if no data exists", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "test", "")}
		gotErr := subject.Delete()
		assert.ErrorContains(t, gotErr, "delete item: nothing to delete")
	})

	t.Run("deletes item", func(t *testing.T) {
		setupKeychainAccessorTest(t)
		subject := &KeychainAccessor{config: config.NewFake("", "test", "")}
		subject.Write(PasswordCredential{})
		gotErr := subject.Delete()
		assert.NoError(t, gotErr)
		gotBool, _ := subject.Has()
		assert.False(t, gotBool)
	})
}
