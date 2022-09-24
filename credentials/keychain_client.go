package credentials

import "github.com/keybase/go-keychain"

type NewKeychainItem func() keychain.Item

type AddKeychainItem func(keychain.Item) error

type DeleteKeychainItem func(keychain.Item) error

type QueryKeychainItem func(keychain.Item) ([]keychain.QueryResult, error)

var _keychainNewItem = keychain.NewItem
var _keychainAddItem = keychain.AddItem
var _keychainDeleteItem = keychain.DeleteItem
var _keychainQueryItem = keychain.QueryItem
