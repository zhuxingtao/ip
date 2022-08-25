package ip

const (
	TypeSymbol  = "symbol"
	TypePlus    = "plus"
	TypeMinus   = "minus"
	TypeEof     = "eof"
	TypeFloat   = "float"
	TypeDefine  = "define"
	TypeSet     = "set!" //可以用来改变外部环境变量的值
	TypeIf      = "if"
	TypeInteger = "integer"
	TypeTrue    = "true"
	TypeFalse   = "false"
	TypeDo      = "do"
	TypeMacro   = "macro"
	TypeString  = "string"
	TypeComment = "comment"
	TypeLParen  = "lparen"
	TypeRParen  = "rparen"
	TypeLSquare = "l_square"
	TypeRSquare = "r_square"
	TypeQuote   = "quote"
	TypeArray   = "list"
)
