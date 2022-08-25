package ip

import (
	"errors"
	"log"
	"strconv"
)

type INode interface {
}

type Node struct {
	Type  string
	Value interface{}
}

const (
	NodeTypeArray  = "array"
	NodeTypeList   = "list"
	NodeTypeNumber = "number"
	NodeTypeSymbol = "symbol"
	NodeTypeString = "string"
)

type Parser struct {
	Lexer *Lexer
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		Lexer: lexer,
	}
}

func (i *Parser) Parse() ([]INode, error) {
	tokens := i.Lexer.Tokenize()
	//fmt.Printf("tokens:%v\n", tokens)
	idx, length := 0, len(tokens)
	nodes := make([]INode, 0)
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

func parserExpr(tokens []Token, idx int) (INode, int, error) {
	if len(tokens) == 0 {
		return nil, 0, errors.New("error parsing expression")
	}
	nextIdx := idx
	var err error
	var expr INode
	switch tokens[idx].Type {
	//case TypeDefine:
	case TypeArray:
		//	(list 0 1)
		idx++
		l := make([]INode, 0)
		for tokens[idx].Type != TypeRParen {
			expr, nextIdx, _ = parserExpr(tokens, idx)
			l = append(l, expr)
			idx = nextIdx
		}
		expr = Node{
			Type:  NodeTypeArray,
			Value: l,
		}
	case TypeLParen:
		l := make([]INode, 0)
		idx++
		if tokens[idx].Type == TypeRParen {
			v := Node{
				Type:  NodeTypeArray,
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

	case TypeLSquare:
	case TypeMacro:
	case TypeInteger, TypeFloat:
		i, err := strconv.ParseFloat(tokens[idx].Value, 64)
		if err != nil {
			return nil, 0, err
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
	case TypeString:
		expr = Node{
			Type:  NodeTypeString,
			Value: tokens[idx].Value,
		}
		idx++
	case TypeTrue, TypeFalse, TypeIf, TypeSymbol, TypeDefine, TypeSet:
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
