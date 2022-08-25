package ip

type TokenType int
type Token struct {
	Type  TokenType
	Value string
}

func newToken(t TokenType, literal string) Token {
	return Token{Type: t, Value: literal}
}
