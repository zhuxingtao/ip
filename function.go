package ip

import (
	"fmt"
	"time"
)

func procCons(args []SExpression) SExpression {
	value := make([]SExpression, 0)
	res := Node{
		Type:  NodeTypeList,
		Value: value,
	}
	if args[0].(Node).Type != NodeTypeList {
		value = append(value, args[0])
	} else {
		value = args[0].(Node).Value.([]SExpression)
	}
	vAstNode := args[1].(Node).Value.([]SExpression)
	if len(vAstNode) != 0 {
		value = append(value, vAstNode...)
	}
	res.Value = value
	return res
}

func procMap(args []SExpression) SExpression {
	function := args[0].(Node).Value.(Procedure)
	array := args[1].(Node)
	res := make([]SExpression, 0, len(array.Value.([]SExpression)))
	for _, item := range array.Value.([]SExpression) {
		v, _ := function.Call([]SExpression{item})
		res = append(res, v)
	}
	return Node{
		Type:  NodeTypeList,
		Value: res,
	}
}

func procAppend(args []SExpression) SExpression {
	a := args[0].(Node)
	b := args[1].(Node)
	va := a.Value.([]SExpression)
	vb := b.Value.([]SExpression)
	res := make([]SExpression, 0, len(va)+len(vb))
	ret := Node{
		Type: NodeTypeList,
	}
	res = append(res, va...)
	res = append(res, vb...)
	ret.Value = res
	return ret
}

func CallCC(args []SExpression) (v SExpression) {
	f := args[0].(Procedure)
	defer func(f Procedure) {
		if r := recover(); r != nil {
			if o, ok := r.(CallCCObj); ok {
				if o.argName == f.params[0] {
					v = o.val
				} else {
					panic(o)
				}
			}
		}
	}(f)
	throw := func(val []SExpression) SExpression {
		o := CallCCObj{val: val[0], argName: f.params[0]}
		panic(o)
	}
	v, _ = f.Call([]SExpression{throw})
	return v
}

func procArgIsBoolean(args []SExpression) SExpression {
	res := Node{
		Type:  NodeTypeBoolean,
		Value: false,
	}
	if len(args) == 0 {
		return res
	}
	v := args[0]
	_, ok := v.(bool)
	res.Value = ok
	return res
}

func procCar(args []SExpression) SExpression {
	if len(args) == 0 {
		return nil
	}
	n := args[0].(Node)
	if n.Type == NodeTypeList {
		v := n.Value
		a := v.([]SExpression)
		if len(a) == 0 {
			return nil
		}
		return a[0]
	}
	return nil
}

func procCdr(args []SExpression) SExpression {
	n := args[0].(Node)
	if n.Type == NodeTypeList {
		v := n.Value.([]SExpression)
		if len(v) > 0 {
			return Node{
				Type:  NodeTypeList,
				Value: v[1:],
			}
		}
	}
	return Node{
		Type:  NodeTypeList,
		Value: make([]SExpression, 0),
	}
}

func procIsAtom(args []SExpression) SExpression {
	res := Node{
		Type:  NodeTypeBoolean,
		Value: false,
	}
	l := args[0].(Node)
	if l.Type == NodeTypeNumber || l.Type == NodeTypeString {
		res.Value = true
	}
	return res
}

func procIsPair(args []SExpression) SExpression {
	res := false
	if args[0].(Node).Type == NodeTypeList {
		v := args[0].(Node).Value.([]SExpression)
		if len(v) == 0 {
			res = false
		} else {
			res = true
		}
	}

	return Node{
		Type:  NodeTypeBoolean,
		Value: res,
	}
}

func procIsNull(args []SExpression) SExpression {
	res := Node{
		Type:  NodeTypeBoolean,
		Value: false,
	}
	switch args[0].(type) {
	case Node:
		n := args[0].(Node)
		if n.Type == NodeTypeList {
			res.Value = len(n.Value.([]SExpression)) == 0
		}
	case []SExpression:
		v := args[0].([]SExpression)
		if len(v) != 0 {
			res.Value = false
		} else {
			res.Value = true
		}
	}
	return res
}

func procIsEqual(args []SExpression) SExpression {
	res := Node{
		Type: NodeTypeBoolean,
	}
	switch args[0].(type) {
	case bool:
		x := args[0].(bool)
		if v, ok := args[1].(bool); ok {
			res.Value = v == x
		}
		return res
	case float64:
		x := args[0].(float64)
		if v, ok := args[1].(float64); ok {
			res.Value = v == x
		}
		return res
	case string:
		x := args[0].(string)
		if v, ok := args[1].(string); ok {
			res.Value = v == x
			return res
		}
		return res
	}
	x, y := args[0].(Node).Value, args[1].(Node).Value
	if (args[0] != nil && args[1] == nil) || (args[0] == nil && args[1] != nil) {
		return res
	}
	n, ok1 := x.(float64)
	m, ok2 := y.(float64)
	if (ok1 && !ok2) || (ok2 && !ok1) {
		panic("= needs numbers")
	}
	if !ok1 && !ok2 {
		a, ok1 := x.(bool)
		b, ok2 := y.(bool)
		if !ok1 || !ok2 {
			a, ok1 := x.(string)
			b, ok2 := y.(string)
			if !ok1 || !ok2 {
				panic("object can't compare")
			}
			res.Value = a == b
			return res
		}
		res.Value = a == b
		return res
	}
	res.Value = n == m
	return res
}

func procLength(args []SExpression) SExpression {
	ret := 0
	n := args[0].(Node)
	switch n.Type {
	case NodeTypeList:
		ret = len(n.Value.([]SExpression))
	case NodeTypeString:
		ret = len(n.Value.(string))
	}
	return Node{
		Type:  NodeTypeNumber,
		Value: float64(ret),
	}
}

func procMul(args []SExpression) SExpression {
	x := args[0].(Node)
	s, ok1 := x.Value.(float64)
	if !ok1 {
		fmt.Printf("* args[0] needs numbers\n")
		return nil
	}
	for _, arg := range args[1:] {
		y := arg.(Node)
		m, ok2 := y.Value.(float64)
		if !ok2 {
			fmt.Printf("* needs numbers\n")
			return nil
		}
		s *= m
	}
	return Node{
		Type:  NodeTypeNumber,
		Value: s,
	}
}

func procDivide(args []SExpression) SExpression {
	x := args[0].(Node)
	s, ok1 := x.Value.(float64)
	if !ok1 {
		fmt.Printf("/ needs numbers args[0]=%v %T\n", args[0], args[0])
		return nil
	}
	for _, arg := range args[1:] {
		y := arg.(Node)
		n, ok1 := y.Value.(float64)
		if n == 0 {
			return "divisor can't be zero"
		}
		if !ok1 {
			fmt.Printf("/ needs numbers")
			return nil
		}
		s /= n
	}
	return Node{
		Type:  NodeTypeNumber,
		Value: s,
	}
}

func procDiv(args []SExpression) SExpression {
	x := args[0].(Node)
	s, ok1 := x.Value.(float64)
	if !ok1 {
		panic("+ needs numbers")
	}
	for _, arg := range args[1:] {
		y := arg.(Node)
		n, ok1 := y.Value.(float64)
		if n == 0 {
			return "divisor can't be zero"
		}
		if !ok1 {
			panic("+ needs numbers")
		}
		s /= n
	}
	b := int64(s)
	return Node{
		Type:  NodeTypeNumber,
		Value: float64(b),
	}
}

func procLTE(args []SExpression) SExpression {
	x := args[0].(Node)
	y := args[1].(Node)
	var a, b float64
	n, ok1 := x.Value.(float64)
	if ok1 {
		a = n
	} else {
		panic("<= needs numbers")
	}
	m, ok2 := y.Value.(float64)
	if ok2 {
		b = m
	} else {
		panic("<= needs numbers")
	}
	return Node{
		Type:  NodeTypeBoolean,
		Value: a <= b,
	}
}

func procGTE(args []SExpression) SExpression {
	x := args[0].(Node).Value
	y := args[1].(Node).Value
	var a, b float64
	n, ok1 := x.(float64)
	if ok1 {
		a = n
	} else {
		panic(">= needs numbers")
	}
	m, ok2 := y.(float64)
	if ok2 {
		b = m
	} else {
		panic(">= needs numbers")
	}
	return Node{
		Type:  NodeTypeBoolean,
		Value: a >= b,
	}
}

func procLess(args []SExpression) SExpression {
	x := args[0].(Node).Value
	y := args[1].(Node).Value
	var a, b float64
	n, ok1 := x.(float64)
	if ok1 {
		a = n
	} else {
		fmt.Printf("> needs numbers\n")
		return nil
	}
	m, ok2 := y.(float64)
	if ok2 {
		b = m
	} else {
		fmt.Printf("> needs numbers")
		return nil
	}
	return Node{
		Type:  NodeTypeBoolean,
		Value: a < b,
	}
}

func procGreater(args []SExpression) SExpression {
	x := args[0].(Node).Value
	y := args[1].(Node).Value
	var a, b float64
	n, ok1 := x.(float64)
	if ok1 {
		a = n
	} else {
		fmt.Printf("> needs numbers\n")
		return nil
	}
	m, ok2 := y.(float64)
	if ok2 {
		b = m
	} else {
		fmt.Printf("> needs numbers")
		return nil
	}
	return Node{
		Type:  NodeTypeBoolean,
		Value: a > b,
	}
}

func procPlus(args []SExpression) SExpression {
	var s float64 = 0
	for _, a := range args {
		arg := a.(Node).Value
		n, ok1 := arg.(float64)
		if !ok1 {
			m, ok2 := arg.(bool)
			if ok2 {
				if m == true {
					s += 1
				}
				continue
			} else {
				fmt.Printf("type=%T %v\n", arg, arg)
				panic("+ needs numbers")
			}
		}
		s += n
	}
	ret := Node{
		Type:  NodeTypeNumber,
		Value: s,
	}
	return ret
}

func procMinus(args []SExpression) SExpression {
	a := args[0].(Node).Value
	s, ok1 := a.(float64)
	if !ok1 {
		b, ok1 := a.(bool)
		if ok1 {
			if b {
				s = 1
			} else {
				s = 0
			}
		} else {
			fmt.Printf("- args[0] needs numbers\n")
			return nil
		}
	}
	for _, a := range args[1:] {
		arg := a.(Node).Value
		m, ok1 := arg.(float64)
		if !ok1 {
			b, ok1 := arg.(bool)
			if ok1 {
				if b {
					s -= 1
				}
				continue
			} else {
				fmt.Printf("- needs numbers\n")
				return nil
			}
		} else {
			s -= m
		}
	}
	ret := Node{
		Type:  NodeTypeNumber,
		Value: s,
	}
	return ret
}

func procQmi(args []SExpression) SExpression {
	x := args[0].(Node).Value
	y := args[1].(Node).Value
	n, ok1 := x.(float64)
	m, ok2 := y.(float64)
	if !ok1 || !ok2 {
		panic("qmi needs numbers")
	}
	ksm := func(a, b int64) float64 {
		var res int64 = 1
		for b > 0 {
			if b&1 == 1 {
				res *= a
			}
			a = a * a
			b = b >> 1
		}
		return float64(res)
	}
	res := ksm(int64(n), int64(m))
	return Node{Type: NodeTypeNumber, Value: res}
}

func procGetVarType(args []SExpression) SExpression {
	x := args[0].(Node)
	return Node{
		Type:  NodeTypeString,
		Value: x.Type,
	}
}

func procList(args []SExpression) SExpression {
	res := make([]SExpression, 0, len(args))
	for _, a := range args {
		v := a.(Node)
		res = append(res, v)
	}
	return Node{
		Type:  NodeTypeList,
		Value: res,
	}
}

func procTime(args []SExpression) SExpression {
	t := time.Now()
	return Node{
		Type:  NodeTypeTime,
		Value: t,
	}
}

func NewFuncNode(f func(args []SExpression) SExpression) Node {
	return Node{
		Type:  NodeTypeBuiltinFunc,
		Value: f,
	}
}
