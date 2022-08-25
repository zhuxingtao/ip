package ip

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	Type  string
	Value SExpression
}

const (
	NodeTypeNil         = "nil"
	NodeTypeQuote       = "quote"
	NodeTypeList        = "list"
	NodeTypeNumber      = "number"
	NodeTypeSymbol      = "symbol"
	NodeTypeBoolean     = "boolean"
	NodeTypeString      = "string"
	NodeTypeFunc        = "user_function"
	NodeTypeBuiltinFunc = "builtin_function"
	NodeTypeTime        = "time"
)

func (i Node) Inspect() string {
	switch i.Type {
	case NodeTypeNumber:
		res := strconv.FormatFloat(i.Value.(float64), 'g', -1, 64)
		return res
	case NodeTypeString:
		return i.Value.(string)
	case NodeTypeBoolean:
		v := i.Value.(bool)
		if v {
			return "#t"
		}
		return "#f"
	case NodeTypeFunc:
		p := i.Value.(Procedure)
		return fmt.Sprintf("params=%v body=%v", p.params, p.body)
	case NodeTypeTime:
		t := i.Value.(time.Time)
		return fmt.Sprintf("%v", t.Format("2006-01-02 15:04:05"))
	case NodeTypeList:
		var out bytes.Buffer

		var elements []string
		for _, e := range i.Value.([]SExpression) {
			elements = append(elements, e.(Node).Inspect())
		}

		out.WriteString("(")
		out.WriteString(strings.Join(elements, " "))
		out.WriteString(")")

		return out.String()
	}
	return ""
}
