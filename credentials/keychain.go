package credentials

import (
	"fmt"

	"github.com/keybase/go-keychain"
	"github.com/tomato-net/vault-agent/config"
)

var _ Reader = (*KeychainAccessor)(nil)
var _ Writer = (*KeychainAccessor)(nil)

type KeychainAccessor struct {
	config config.Config
}

func NewKeychainAccessor(config config.Config) *KeychainAccessor {
	return &KeychainAccessor{config: config}
}

func (a *KeychainAccessor) Read() (PasswordCredential, error) {
	account := a.config.Username()
	if account == "" {
		return PasswordCredential{}, fmt.Errorf("no username in config")
	}

	query := keychainItem(account)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		return PasswordCredential{}, fmt.Errorf("read: %w", err)
	}

	if len(results) != 1 {
		return PasswordCredential{}, fmt.Errorf("no results found for %s", account)
	}

	password := string(results[0].Data)
	return PasswordCredential{Data: password}, nil
}

func (a *KeychainAccessor) Has() (bool, error) {
	cred, err := a.Read()
	if err != nil {
		return false, fmt.Errorf("read: %w", err)
	}

	return cred.Data != "", nil
}

func (a *KeychainAccessor) Write(credential PasswordCredential) error {
	account := a.config.Username()
	if account == "" {
		return fmt.Errorf("no username in config")
	}

	item := keychainItem(account)
	item.SetMatchLimit(keychain.MatchLimitOne)
	item.SetData([]byte(credential.Data))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	err := keychain.AddItem(item)
	if err != nil {
		return fmt.Errorf("add item: %w", err)
	}

	return nil
}

func (a *KeychainAccessor) Delete() error {
	account := a.config.Username()
	if account == "" {
		return fmt.Errorf("no username in config")
	}

	item := keychainItem(account)
	if err := keychain.DeleteItem(item); err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	return nil
}

func keychainItem(account string) keychain.Item {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService("vault-agent.tomato-net.github.com")
	item.SetLabel("vault agent credentials")
	item.SetAccount(account)
	return item
}
