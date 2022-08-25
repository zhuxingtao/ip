package ip

type TokenType string
type Token struct {
	Type  TokenType
	Value string
}

func newToken(t TokenType, literal string) Token {
	return Token{Type: t, Value: literal}
}
