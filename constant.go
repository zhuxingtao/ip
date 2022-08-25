package ip

const (
	TypeSymbol  = iota //= "symbol"
	TypeEof            //= "eof"
	TypeFloat          //= "float"
	TypeDefine         //= "define"
	TypeSet            //= "set!"
	TypeIf             //= "if"
	TypeInteger        //= "integer"
	TypeTrue           //= "true"
	TypeFalse          //= "false"
	TypeString         //= "string"
	TypeComment        //= "comment"
	TypeLParen         //= "("
	TypeRParen         //= ")"
	TypeQuote          //= "quote"
)
