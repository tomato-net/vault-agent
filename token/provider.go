package token

type Provider interface {
	Token() (Token, error)
}
