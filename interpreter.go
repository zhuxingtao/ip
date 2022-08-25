package ip

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	UndefObj Undef
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
	store  map[string]SExpression
}

func NewEmptyEnv() *Env {
	env := &Env{
		parent: nil,
		store:  make(map[string]SExpression),
	}
	return env
}

func NewEnv(params []string, args []SExpression, outer *Env) *Env {
	env := &Env{
		parent: outer,
		store:  make(map[string]SExpression),
	}
	for i := 0; i < len(params); i++ {
		env.store[params[i]] = args[i]
	}
	return env
}

func StandardEnv() *Env {
	env := NewEmptyEnv()
	env.store["cons"] = NewFuncNode(procCons)
	env.store["map"] = NewFuncNode(procMap)
	env.store["expt"] = NewFuncNode(procQmi)
	env.store["append"] = NewFuncNode(procAppend)
	env.store["call/cc"] = NewFuncNode(CallCC)
	env.store["pow"] = env.store["expt"]
	env.store["pi"] = Node{
		Type:  NodeTypeNumber,
		Value: 3.1415926535,
	}
	env.store["false"] = Node{
		Type:  NodeTypeBoolean,
		Value: false,
	}
	env.store["true"] = Node{
		Type:  NodeTypeBoolean,
		Value: true,
	}
	env.store["null?"] = NewFuncNode(procIsNull)
	env.store["pair?"] = NewFuncNode(procIsPair)
	env.store["atom?"] = NewFuncNode(procIsAtom)
	env.store["="] = NewFuncNode(procIsEqual)
	env.store["length"] = NewFuncNode(procLength)
	env.store["equal?"] = env.store["="]
	env.store["<="] = NewFuncNode(procLTE)
	env.store[">="] = NewFuncNode(procGTE)
	env.store["<"] = NewFuncNode(procLess)
	env.store[">"] = NewFuncNode(procGreater)
	env.store["+"] = NewFuncNode(procPlus)
	env.store["-"] = NewFuncNode(procMinus)
	env.store["*"] = NewFuncNode(procMul)
	env.store["/"] = NewFuncNode(procDivide)
	env.store["//"] = NewFuncNode(procDiv)
	env.store["car"] = NewFuncNode(procCar)
	env.store["first"] = NewFuncNode(procCar)
	env.store["cdr"] = NewFuncNode(procCdr)
	env.store["rest"] = NewFuncNode(procCdr)
	env.store["boolean?"] = NewFuncNode(procArgIsBoolean)
	env.store["type?"] = NewFuncNode(procGetVarType)
	env.store["list"] = NewFuncNode(procList)
	env.store["time"] = NewFuncNode(procTime)
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

func (i *Procedure) Call(args []SExpression) (SExpression, error) {
	res, err := Eval(i.body, NewEnv(i.params, args, i.env))
	if err != nil {
		return UndefObj, err
	}
	return res, nil
}

func makeArgs(l []Node, env *Env) ([]SExpression, error) {
	args := make([]SExpression, 0, len(l))
	for _, arg := range l {
		v, err := Eval(arg, env)
		if err != nil {
			return args, err
		}
		args = append(args, v)
	}
	return args, nil
}

func toStrings(params Node) []string {
	if params.Type != NodeTypeList {
		panic("params must be a list")

	}
	l := params.Value.([]SExpression)
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

func Eval(node SExpression, env *Env) (SExpression, error) {
	n, ok := node.(Node)
	if !ok {
		return node, nil
	}
	if n.Type == NodeTypeNumber {
		return n, nil
	}
	if n.Type == NodeTypeString {
		return n, nil
	}
	if n.Type == NodeTypeBoolean {
		return n, nil
	}
	if n.Type == NodeTypeSymbol {
		str := n.Value.(string)
		return env.Find(str).store[str], nil
	}
	if n.Type == NodeTypeQuote {
		exp := n.Value.(Node)
		return makeList(exp), nil
	}

	nnodes := n.Value.([]SExpression)
	if len(nnodes) == 0 { // 空表达式
		return n, nil
	}
	nodes := make([]Node, 0, len(nnodes))
	for i := 0; i < len(nnodes); i++ {
		nodes = append(nodes, nnodes[i].(Node))
	}
	if nodes[0].Type == NodeTypeSymbol {
		if nodes[0].Value.(string)[0] == '\'' {
			return nodes[0].Value.(string)[1:], nil
		}
		switch nodes[0].Value {
		case "quote", "'":
			return makeList(nodes[1]), nil
		case "if":
			res, err := evalIf(nodes, env)
			return res, err
		case "define":
			err := evalDefine(nodes, env)
			if err != nil {
				return nil, err
			}
			return nil, nil
		case "set!":
			var err error
			car, cdr := nodes[1], nodes[2]
			if car.Type == NodeTypeSymbol {
				str := car.Value.(string)
				outerEnv := env.Find(str)
				outerEnv.store[str], err = Eval(cdr, env)
				return nil, err
			} else {
				return nil, errors.New("set! needs a symbol\n")
			}
		case "let":
			res, err := evalLet(nodes[1:], env)
			return res, err
		case "lambda":
			params, body := nodes[1], nodes[2]
			p := newProcedure(params, body, env)
			return Node{Type: NodeTypeFunc, Value: p}, nil
		case "display":
			res, err := Eval(nodes[1], env)
			return res, err
		case "list":
			args, err := makeArgs(nodes[1:], env)
			if err != nil {
				return nil, err
			}
			res := make([]SExpression, 0, len(args))
			for _, a := range args {
				v := a.(Node)
				res = append(res, v)
			}
			return Node{
				Type:  NodeTypeList,
				Value: res,
			}, nil
		case "begin":
			var res SExpression
			var err error
			for _, n := range nodes[1:] {
				res, err = Eval(n, env)
				if err != nil {
					return nil, err
				}
			}
			return res, nil
		}
	}
	if nodes[0].Type == NodeTypeNumber || nodes[0].Type == NodeTypeString {
		res := make([]SExpression, 0, len(nodes))
		for i := 0; i < len(nodes); i++ {
			v, err := Eval(nodes[i], env)
			if err != nil {
				return nil, err
			}
			res = append(res, v)
		}
		return Node{
			Type:  NodeTypeList,
			Value: res,
		}, nil
	}
	car, err := Eval(nodes[0], env)
	if err != nil {
		return nil, err
	}
	args, err := makeArgs(nodes[1:], env)
	if err != nil {
		return nil, err
	}
	//	fmt.Printf("car=%v args=%v\n", car, args)
	if f, ok := isBuiltinFunc(car); ok {
		res := f(args)
		return res, nil
	}
	// user define func
	if n, ok := car.(Node); ok {
		proc := n.Value.(Procedure)
		if len(args) < len(proc.params) {
			e := fmt.Sprintf("the argument count of  function is wrong needParams=%v body=%v\n", proc.params, proc.body)
			return nil, errors.New(e)
		}
		res, err := proc.Call(args)
		return res, err
	}
	return nil, nil
}

func isFalse(x SExpression) bool {
	b, ok := x.(Node)
	if ok {
		if b.Type == NodeTypeBoolean {
			return b.Value == false
		}
		if b.Type == NodeTypeList {
			return len(b.Value.([]SExpression)) == 0
		}
		return false
	} else {
		b, ok := x.([]SExpression)
		if ok {
			if len(b) == 0 {
				return true
			}
			return false
		}
		return false
	}
}

func isBuiltinFunc(x SExpression) (func([]SExpression) SExpression, bool) {
	n, ok := x.(Node)
	if !ok {
		return nil, false
	}
	ok = n.Type == NodeTypeBuiltinFunc
	if !ok {
		return nil, false
	}
	return n.Value.(func([]SExpression) SExpression), ok
}

func isList(x SExpression) bool {
	n, ok := x.(Node)
	return ok && n.Type == NodeTypeList
}

func makeList(n Node) SExpression {
	var ret strings.Builder
	res := make([]SExpression, 0)
	if isList(n) {
		v := n.Value.([]SExpression)
		if len(v) == 0 {
			return Node{Type: NodeTypeList, Value: make([]SExpression, 0)}
		}
		for _, item := range v {
			node := item.(Node)
			if node.Type == NodeTypeSymbol {
				res = append(res, Node{Type: NodeTypeString, Value: node.Value})
			} else if node.Type == NodeTypeList {
				tmp := makeList(node)
				if vv, ok := tmp.([]SExpression); ok {
					res = append(res, vv...)
				} else {
					res = append(res, tmp)
				}
			} else if node.Type == NodeTypeNumber {
				v := node.Value.(float64)
				vv := strconv.FormatFloat(v, 'g', 10, 64)
				res = append(res, Node{Type: NodeTypeString, Value: vv})
			} else if node.Type == NodeTypeString {
				res = append(res, node)
			}
		}
	} else if n.Type == NodeTypeSymbol || n.Type == NodeTypeString {
		res = append(res, n)
	}
	if n.Type == NodeTypeSymbol {
		ret.WriteString(n.Value.(string))
		return Node{
			Type:  NodeTypeString,
			Value: ret.String(),
		}
	}
	if n.Type == NodeTypeNumber {
		tmp := n.Value.(float64)
		s := strconv.FormatFloat(tmp, 'g', -1, 64)
		ret.WriteString(s)
		return Node{
			Type:  NodeTypeString,
			Value: ret.String(),
		}
	}
	fmt.Printf("\nres=%+v\n", res)
	retList := make([]SExpression, 0)
	for i := 0; i < len(res); i++ {
		v := res[i].(Node)
		if v.Type == NodeTypeString {
			retList = append(retList, v)
		}
		if v.Type == NodeTypeNumber {
			retList = append(retList, v)
		}
		if v.Type == NodeTypeList {
			retList = append(retList, v)
		}
	}
	return Node{Type: NodeTypeList, Value: retList}
}

type CallCCObj struct {
	val     SExpression
	argName string
}

type Undef struct {
}

type NilObj struct {
}

func (i Undef) String() string {
	return "#undefine obj"
}

func evalLet(args []Node, env *Env) (SExpression, error) {
	if len(args) < 2 {
		return Undef{}, errors.New("let args < 2")
	}
	bindings, ok := args[0].Value.([]SExpression)
	if !ok {
		return Undef{}, errors.New("")
	}
	newEnv := &Env{store: make(map[string]SExpression), parent: env}
	for _, exp := range bindings {
		binding, ok := exp.(Node)
		if !ok {
			return UndefObj, errors.New("let: syntax error (not a valid binding)")
		}
		s := binding.Value.([]SExpression)
		val, err := Eval(s[1], newEnv)
		if err != nil {
			return Undef{}, nil
		}
		newEnv.store[s[0].(Node).Value.(string)] = val
	}
	var ret SExpression
	var err error
	for _, exp := range args[1:] {
		ret, err = Eval(exp, newEnv)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil

}

func evalDefine(nodes []Node, env *Env) error {
	car, cdr := nodes[1], nodes[2]
	var err error
	if car.Type == NodeTypeSymbol {
		str := car.Value.(string)
		env.store[str], err = Eval(cdr, env)
		return err
	} else if car.Type == NodeTypeList {
		v := car.Value.([]SExpression)
		function := v[0].(Node).Value.(string)
		tmp := make([]string, 0)
		for _, a := range v[1:] {
			item := a.(Node).Value.(string)
			tmp = append(tmp, item)
		}
		p := Procedure{params: tmp, body: cdr, env: env}
		env.store[function] = p
		return nil
	}
	return errors.New("define parse error")
}

func evalIf(exps []Node, env *Env) (SExpression, error) {
	test, conSeq, alt := exps[1], exps[2], exps[3]
	r, err := Eval(test, env)
	if err != nil {
		return nil, err
	}
	if isFalse := isFalse(r); isFalse {
		return Eval(alt, env)
	} else {
		return Eval(conSeq, env)
	}
}
