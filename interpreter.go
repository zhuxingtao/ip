package ip

import (
	"fmt"
	"strconv"
)

type Interpreter struct {
	Text         string
	Position     int
	CurrentToken *Token
}

func NewInterpreter(text string) *Interpreter {
	return &Interpreter{
		Text: text,
	}
}

func (i *Interpreter) Error(s string) error {
	return nil
}

type Env struct {
	parent *Env
	store  map[string]interface{}
}

func NewEmptyEnv() *Env {
	env := &Env{
		parent: nil,
		store:  make(map[string]interface{}),
	}
	return env
}

func NewEnv(params []string, args []interface{}, outer *Env) *Env {
	env := &Env{
		parent: outer,
		store:  make(map[string]interface{}),
	}
	for i := 0; i < len(params); i++ {
		env.store[params[i]] = args[i]
	}
	return env
}

func StandardEnv() *Env {
	env := NewEmptyEnv()
	env.store["cons"] = func(args []interface{}) interface{} {
		res := make([]interface{}, 1)
		res[0] = args[0]
		vAstNode := args[1].([]interface{})
		if len(vAstNode) == 0 {
			return res
		} else {
			res = append(res, vAstNode...)
		}
		return res
	}
	env.store["map"] = func(args []interface{}) interface{} {
		function := args[0].(Procedure)
		array := args[1].([]interface{})
		res := make([]interface{}, 0, len(array))
		for _, item := range array {
			v := function.Call([]interface{}{item})
			res = append(res, v)
		}
		return res
	}
	env.store["expt"] = func(args []interface{}) interface{} {
		n, ok1 := args[0].(float64)
		m, ok2 := args[1].(float64)
		if !ok1 || !ok2 {
			panic("+ needs numbers")
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
		return ksm(int64(n), int64(m))
	}
	env.store["call/cc"] = CallCC
	env.store["pow"] = env.store["expt"]
	env.store["pi"] = 3.1415926535
	env.store["false"] = false
	env.store["true"] = true
	env.store["="] = func(args []interface{}) interface{} {
		//fmt.Printf("a = %v %T b= %v %T \n", args[0], args[0], args[1], args[1])
		if (args[0] != nil && args[1] == nil) || (args[0] == nil && args[1] != nil) {
			return false
		}
		n, ok1 := args[0].(float64)
		m, ok2 := args[1].(float64)
		if (ok1 && !ok2) || (ok2 && !ok1) {
			panic("= needs numbers")
		}
		if !ok1 && !ok2 {
			a, ok1 := args[0].(bool)
			b, ok2 := args[1].(bool)
			if !ok1 || !ok2 {
				//fmt.Printf("object can't compare \n")
				a, ok1 := args[0].(string)
				b, ok2 := args[1].(string)
				if !ok1 || !ok2 {
					panic("object can't compare")
				}
				return a == b
			}
			return a == b
		}
		return n == m
		//n, ok1 := args[0].(int64)
		//m, ok2 := args[1].(int64)

		//if (!ok2 && ok1) || (!ok1 && ok2) {
		//	return false
		//}
		//if ok1 && ok2 {
		//	return m == n
		//}
		//a, ok1 := args[0].(float64)
		//b, ok2 := args[1].(float64)
		//if !ok1 || !ok2 {
		//	return false
		//}
		//return a == b
	}
	env.store["equal?"] = env.store["="]
	env.store["<="] = func(args []interface{}) interface{} {
		var a, b float64
		n, ok1 := args[0].(float64)
		if ok1 {
			a = n
		} else {
			panic("> needs numbers")
		}
		m, ok2 := args[1].(float64)
		if ok2 {
			b = m
		} else {
			panic("> needs numbers")
		}
		return a <= b
	}
	env.store[">="] = func(args []interface{}) interface{} {
		var a, b float64
		n, ok1 := args[0].(float64)
		if ok1 {
			a = n
		} else {
			panic(">= needs numbers")
		}
		m, ok2 := args[1].(float64)
		if ok2 {
			b = m
		} else {
			panic(">= needs numbers")
		}
		return a >= b
	}
	env.store["<"] = func(args []interface{}) interface{} {
		var a, b float64
		n, ok1 := args[0].(float64)
		if ok1 {
			a = n
		} else {
			panic("> needs numbers")
		}
		m, ok2 := args[1].(float64)
		if ok2 {
			b = m
		} else {
			panic("> needs numbers")
		}
		return a < b
	}
	env.store[">"] = func(args []interface{}) interface{} {
		var a, b float64
		n, ok1 := args[0].(float64)
		if ok1 {
			a = n
		} else {
			panic("> needs numbers")
		}
		m, ok2 := args[1].(float64)
		if ok2 {
			b = m
		} else {
			panic("> needs numbers")
		}
		return a > b
	}
	env.store["+"] = func(args []interface{}) interface{} {
		var s float64 = 0
		for _, arg := range args {
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
		return s
	}
	env.store["-"] = func(args []interface{}) interface{} {
		s, ok1 := args[0].(float64)
		if !ok1 {
			b, ok1 := args[0].(bool)
			if ok1 {
				if b {
					s = 1
				} else {
					s = 0
				}
			} else {
				fmt.Printf("- ffneeds numbers")
				return nil
			}
		}
		for _, arg := range args[1:] {
			m, ok1 := arg.(float64)
			if !ok1 {
				b, ok1 := arg.(bool)
				if ok1 {
					if b {
						s -= 1
					}
					continue
				} else {
					fmt.Printf("- needs numbers")
					return nil
				}
			} else {
				s -= m
			}
		}
		return s
	}
	env.store["*"] = func(args []interface{}) interface{} {
		s, ok1 := args[0].(float64)
		if !ok1 {
			panic("* args[0] needs numbers")
		}

		for _, arg := range args[1:] {
			m, ok2 := arg.(float64)
			if !ok2 {
				fmt.Printf("args[1]=%v %T\n", args[1], args[1])
				panic("* needs numbers")
			}
			s *= m
		}
		return s
	}
	env.store["/"] = func(args []interface{}) interface{} {
		s, ok1 := args[0].(float64)
		if !ok1 {
			panic("+ needs numbers")
		}
		for _, arg := range args[1:] {
			n, ok1 := arg.(float64)
			if n == 0 {
				return "divisor can't be zero"
			}
			if !ok1 {
				panic("+ needs numbers")
			}
			s /= n
		}
		return s
	}
	env.store["//"] = func(args []interface{}) interface{} {
		s, ok1 := args[0].(float64)
		if !ok1 {
			panic("+ needs numbers")
		}
		for _, arg := range args[1:] {
			n, ok1 := arg.(float64)
			if n == 0 {
				return "divisor can't be zero"
			}
			if !ok1 {
				panic("+ needs numbers")
			}
			s /= n
		}
		b := int64(s)
		return float64(b)
	}
	env.store["car"] = func(args []interface{}) interface{} {
		if len(args) == 0 {
			return nil
		}
		v := args[0].([]interface{})
		if len(v) == 0 {
			return nil
		}
		return v[0]
	}
	env.store["first"] = env.store["car"]
	env.store["cdr"] = func(args []interface{}) interface{} {
		v := args[0].([]interface{})
		if len(v) > 0 {
			return v[1:]
		}
		return []interface{}{}
	}
	env.store["rest"] = env.store["cdr"]
	env.store["boolean?"] = func(args []interface{}) interface{} {
		if len(args) == 0 {
			return false
		}
		v := args[0]
		_, ok := v.(bool)
		return ok
	}
	return env
}

func (i *Env) Find(v string) *Env {
	if i.store[v] != nil || i.parent == nil {
		return i
	}
	return i.parent.Find(v)
}

type Procedure struct {
	params []string
	body   Node
	env    *Env
}

func newProcedure(params Node, body Node, env *Env) Procedure {
	strList := toStrings(params)
	return Procedure{strList, body, env}
}

func (i *Procedure) Call(args []interface{}) interface{} {
	return Eval(i.body, NewEnv(i.params, args, i.env))
}

func makeArgs(l []Node, env *Env) []interface{} {
	args := make([]interface{}, 0, len(l))
	for _, arg := range l {
		args = append(args, Eval(arg, env))
	}
	return args
}

func toStrings(params Node) []string {
	if params.Type != NodeTypeList {
		panic("params must be a list")

	}
	l := params.Value.([]INode)
	strList := make([]string, 0, len(l))
	for i := 0; i < len(l); i++ {
		tmp := l[i].(Node)
		if tmp.Type == NodeTypeSymbol {
			str := tmp.Value.(string)
			strList = append(strList, str)
		} else {
			panic("params need symbols")
		}
	}
	return strList
}

var k int

func Eval(n Node, env *Env) interface{} {
	if n.Type == NodeTypeNumber {
		return n.Value.(float64)
	}
	if n.Type == NodeTypeString {
		return n.Value.(string)
	}
	if n.Type == NodeTypeSymbol {
		str := n.Value.(string)
		return env.Find(str).store[str]
	}
	if n.Type == NodeTypeArray {
		v := n.Value.([]INode)
		res := make([]interface{}, 0, len(v))
		for _, tmp := range v {
			tmpNode := tmp.(Node)
			item := Eval(tmpNode, env)
			res = append(res, item)
		}
		return res
	}
	nnodes := n.Value.([]INode)
	if len(nnodes) == 0 { // 空表达式
		return []interface{}{}
	}
	nodes := make([]Node, 0, len(nnodes))
	for i := 0; i < len(nnodes); i++ {
		nodes = append(nodes, nnodes[i].(Node))
	}
	if nodes[0].Type == NodeTypeSymbol {
		switch nodes[0].Value {
		case "quote":
			return makeList(nodes[1])
		case "if":
			test, conseq, alt := nodes[1], nodes[2], nodes[3]
			r := Eval(test, env)
			//fmt.Printf("r=%v %T\n", r, r)
			if isFalse := isFalse(r); isFalse {
				return Eval(alt, env)
			} else {
				return Eval(conseq, env)
			}
		case "define":
			car, cdr := nodes[1], nodes[2]
			if car.Type == NodeTypeSymbol {
				str := car.Value.(string)
				env.store[str] = Eval(cdr, env)
				return nil
			} else if car.Type == NodeTypeList {
				v := car.Value.([]INode)
				function := v[0].(Node).Value.(string)
				tmp := make([]string, 0)
				for _, a := range v[1:] {
					item := a.(Node).Value.(string)
					tmp = append(tmp, item)
				}
				p := Procedure{params: tmp, body: cdr, env: env}
				env.store[function] = p
				return nil
			} else {
				panic("define parse error")
			}
		case "set!":
			car, cdr := nodes[1], nodes[2]
			if car.Type == NodeTypeSymbol {
				str := car.Value.(string)
				outerEnv := env.Find(str)
				outerEnv.store[str] = Eval(cdr, env)
				return nil
			} else {
				fmt.Printf("set! needs a symbol\n")
				return nil
			}
		case "lambda":
			//fmt.Printf("nodes=%v\n\n", nodes)
			params, body := nodes[1], nodes[2]
			p := newProcedure(params, body, env)
			//	fmt.Printf("p=%+v", p)
			return p
		case "list":
			res := make([]interface{}, 0, len(nodes)-1)
			for i := 1; i < len(nodes); i++ {
				res = append(res, Eval(nodes[i], env))
			}
			return res
		case "display":
			res := Eval(nodes[1], env)
			fmt.Printf("%v\n", res)
		}
	}
	if nodes[0].Type == NodeTypeNumber || nodes[0].Type == NodeTypeString {
		res := make([]interface{}, 0, len(nodes))
		for i := 0; i < len(nodes); i++ {
			res = append(res, Eval(nodes[i], env))
		}
		return res
	}
	if nodes[0].Type == NodeTypeArray {
		v := n.Value.([]INode)
		tmp := v[0]
		tmpNode := tmp.(Node)
		item := Eval(tmpNode, env)
		return item
	}
	car := Eval(nodes[0], env)
	args := makeArgs(nodes[1:], env)
	if f, ok := isFunc(car); ok {
		res := func() interface{} {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("err=%v", r)
				}
			}()
			return f(args)
			//	fmt.Printf("car = %v res = %v\n", car, res)
		}()
		return res
	}
	if proc, ok := car.(Procedure); ok {
		if len(args) < len(proc.params) {
			fmt.Printf("the argument count of  function is wrong\n")
			return nil
		}
		res := proc.Call(args)
		return res
	}
	return ""
}

func isFalse(x interface{}) bool {
	b, ok := x.(bool)
	if ok {
		return !b
	} else {
		b, ok := x.([]interface{})
		if ok {
			if len(b) == 0 {
				return true
			}
			return false
		}
		return false
	}
}

func isFunc(x interface{}) (func([]interface{}) interface{}, bool) {
	prim, ok := x.(func([]interface{}) interface{})
	return prim, ok
}

func isList(x interface{}) ([]interface{}, bool) {
	l, ok := x.([]interface{})
	return l, ok
}

func makeList(n Node) interface{} {
	res := make([]interface{}, 0)
	if n.Type == NodeTypeList {
		v := n.Value.([]INode)
		if len(v) == 0 {
			res = append(res, "()")
		}
		for _, item := range v {
			node := item.(Node)
			if node.Type == NodeTypeSymbol {
				res = append(res, node.Value.(string))
			} else if node.Type == NodeTypeList {
				tmp := makeList(node)
				if vv, ok := tmp.([]interface{}); ok {
					res = append(res, vv...)
				} else {
					res = append(res, tmp)
				}
			} else if node.Type == NodeTypeNumber {
				v := node.Value.(float64)
				vv := strconv.FormatFloat(v, 'g', 10, 64)
				res = append(res, vv)
			}
		}
	} else if n.Type == NodeTypeSymbol {
		res = append(res, n.Value.(string))
	}
	if n.Type == NodeTypeSymbol {
		return res[0]
	}
	return res
}

type CallCCObj struct {
	val     interface{}
	argName string
}

func CallCC(args []interface{}) (v interface{}) {
	f := args[0].(Procedure)
	defer func(f Procedure) {
		if r := recover(); r != nil {
			//fmt.Printf("calcc r= %v a=%v\n", r, args[0].(Procedure))
			if o, ok := r.(CallCCObj); ok {
				if o.argName == f.params[0] {
					v = o.val
				} else {
					panic(o)
				}
			}
		}
	}(f)
	throw := func(val []interface{}) interface{} {
		o := CallCCObj{val: val[0], argName: f.params[0]}
		panic(o)
	}
	v = f.Call([]interface{}{throw})
	return v
}
