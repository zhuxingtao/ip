package ip

import (
	"io"
	"io/ioutil"
	"log"
	"unicode"
)

type Lexer struct {
	Input        string
	Position     int
	ReadPosition int
	Char         byte
}

func NewLexer(input string) *Lexer {
	return &Lexer{Input: input}
}

// read the next char
func (i *Lexer) advance() {
	i.Char = i.peek()
	i.Position = i.ReadPosition
	i.ReadPosition += 1
}

// get the next char
func (i *Lexer) peek() byte {
	if i.ReadPosition >= len(i.Input) {
		return 0
	}
	return i.Input[i.ReadPosition]
}

func (i *Lexer) skipWhiteSpace() {
	for unicode.IsSpace(rune(i.Char)) || i.Char == '\n' {
		i.advance()
	}

}

func (i *Lexer) getFloat(start int) Token {
	//advance to skip .
	i.advance()
	for unicode.IsDigit(rune(i.peek())) {
		i.advance()
	}
	return newToken(TypeFloat, i.Input[start:i.ReadPosition])
}

func (i *Lexer) getInteger() Token {
	old := i.Position
	if i.Char == '-' {
		//skip first char if minus
		i.advance()
	}
	//peek and not advance since advance is called at the end of scanToken, and this could cause us to jump and skip a step
	for unicode.IsDigit(rune(i.peek())) {
		i.advance()
	}
	if i.peek() == '.' {
		return i.getFloat(old)
	}
	return newToken(TypeInteger, i.Input[old:i.ReadPosition])
}

func (i *Lexer) getSymbol() Token {
	old := i.Position
	for !unicode.IsSpace(rune(i.peek())) && i.peek() != 0 && i.peek() != ')' && i.peek() != ']' && i.peek() != '(' {
		i.advance()
	}
	//use position because when l.Char is at a space, l.ReadPosition will be one ahead
	val := i.Input[old:i.ReadPosition]
	var token Token
	switch val {
	case "set!":
		token = newToken(TypeSet, val)
	case "define":
		token = newToken(TypeDefine, val)
	case "if":
		token = newToken(TypeIf, val)
	case "#t":
		token = newToken(TypeTrue, "true")
	case "#f", "nil":
		token = newToken(TypeFalse, "false")
		//	case "list":
		//		token = newToken(TypeArray, val)
	//will add others later
	default:
		token = newToken(TypeSymbol, val)
	}
	return token
}

//function to get entire string or symbol token
func (i *Lexer) getUntil(until byte, token TokenType, after bool) Token {
	old := i.Position
	//get until assumes we eat the last token, which is why we don't use peek
	for i.Char != until && i.Char != 0 {
		i.advance()
	}
	if after && i.Char != 0 {
		i.advance()
	}
	return newToken(token, i.Input[old:i.Position])
}

func (i *Lexer) scanToken() *Token {
	//skips white space and new lines
	i.skipWhiteSpace()
	var token Token
	switch i.Char {
	case '(':
		token = newToken(TypeLParen, "(")
	case ')':
		token = newToken(TypeRParen, ")")
	case '\'':
		token = newToken(TypeQuote, "'")
	case '-':
		if unicode.IsDigit(rune(i.peek())) {
			token = i.getInteger()
		} else {
			token = i.getSymbol()
		}

	case ';':
		if i.peek() == ';' {
			//current char is ; and next char is ; so advance twice
			i.advance()
			i.advance()
			token = i.getUntil(';', TypeComment, false)
		} else {
			token = i.getUntil('\n', TypeComment, false)
		}

	case '"':
		//skip the first "
		i.advance()
		token = i.getUntil('"', TypeString, false)
	case 0:
		token = newToken(TypeEof, "EOF")
	default:
		if unicode.IsDigit(rune(i.Char)) {
			token = i.getInteger()
		} else {
			//more things potentially here
			token = i.getSymbol()
		}
	}
	i.advance()
	return &token
}

func (i *Lexer) Tokenize() []Token {
	var tokens []Token
	//set the first character
	i.advance()
	for i.Position < len(i.Input) {
		next := i.scanToken()
		if next.Type != TypeComment {
			tokens = append(tokens, *next)
		}
	}
	return tokens
}

//Takes as input the source code as a string and returns a list of tokens
func Read(reader io.Reader) []Token {
	source := loadReader(reader)
	l := NewLexer(source)
	tokens := l.Tokenize()
	return tokens
}

func loadReader(reader io.Reader) string {
	//todo: ReadAll puts everything in memory, very inefficient for large files
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("Error trying to read source file: ", err)
	}
	return string(bs)
}
