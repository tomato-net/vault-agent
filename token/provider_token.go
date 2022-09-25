package token

var _ Provider = (*ProviderToken)(nil)

type ProviderToken struct {
	token string
}

func NewProviderToken(token string) *ProviderToken {
	return &ProviderToken{
		token: token,
	}
}

func (t *ProviderToken) Token() (Token, error) {
	return Token{
		Token: t.token,
	}, nil
}
