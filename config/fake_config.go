package config

type FakeConfig struct {
	vaultAddress string
	username     string
	tokenFile    string
}

func NewFake(vaultAddress, username, tokenFile string) *FakeConfig {
	return &FakeConfig{
		vaultAddress: vaultAddress,
		username:     username,
		tokenFile:    tokenFile,
	}
}

func (f *FakeConfig) VaultAddress() string {
	return f.vaultAddress
}

func (f *FakeConfig) Username() string {
	return f.username
}

func (f *FakeConfig) TokenFile() string {
	return f.tokenFile
}
