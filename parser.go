package ip

import (
	"errors"
	"log"
	"strconv"
)

type SExpression interface {
}

type Parser struct {
	Lexer *Lexer
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		Lexer: lexer,
	}
}

func (i *Parser) Parse() ([]SExpression, error) {
	tokens := i.Lexer.Tokenize()
	//fmt.Printf("tokens:%v\n", tokens)
	idx, length := 0, len(tokens)
	nodes := make([]SExpression, 0)
	for idx < length && tokens[idx].Type != TypeEof {
		expr, nextIdx, err := parserExpr(tokens, idx)
		if err != nil {
			log.Fatal("Error parsing tokens:", err)
		}
		//fmt.Printf("expr=%v", expr)
		idx = nextIdx
		nodes = append(nodes, expr)
	}
	return nodes, nil
}

var (
	nilNode = Node{
		Type:  NodeTypeNil,
		Value: nil,
	}
)

func parserExpr(tokens []Token, idx int) (SExpression, int, error) {
	if len(tokens) == 0 {
		return nilNode, 0, errors.New("error parsing expression")
	}
	nextIdx := idx
	var err error
	var expr SExpression
	switch tokens[idx].Type {
	case TypeLParen:
		l := make([]SExpression, 0)
		idx++
		if tokens[idx].Type == TypeRParen {
			v := Node{
				Type:  NodeTypeList,
				Value: l,
			}
			return v, idx + 1, nil
		}
		for tokens[idx].Type != TypeRParen {
			expr, nextIdx, _ = parserExpr(tokens, idx)
			l = append(l, expr)
			idx = nextIdx
		}
		expr = Node{
			Type:  NodeTypeList,
			Value: l,
		}
		idx++

	case TypeInteger, TypeFloat:
		i, err := strconv.ParseFloat(tokens[idx].Value, 64)
		if err != nil {
			return nilNode, 0, err
		}
		idx++
		expr = Node{
			Type:  NodeTypeNumber,
			Value: i,
		}
	case TypeQuote:
		idx++
		nextExpr, nextIdx, errorL := parserExpr(tokens, idx)
		if errorL != nil {
			log.Fatal("Error parsing quote!")
		}
		expr = nextExpr
		idx = nextIdx
		expr = Node{
			Type:  NodeTypeQuote,
			Value: expr,
		}
	case TypeString:
		expr = Node{
			Type:  NodeTypeString,
			Value: tokens[idx].Value,
		}
		idx++
	case TypeTrue, TypeFalse:
		expr = Node{
			Type:  NodeTypeBoolean,
			Value: tokens[idx].Value,
		}
	case TypeIf, TypeSymbol, TypeDefine, TypeSet:
		expr = Node{
			Type:  NodeTypeSymbol,
			Value: tokens[idx].Value,
		}
		idx++
	default:
		log.Fatal("parse error", tokens[idx], "\n", tokens)
	}

	return expr, idx, err
}
