package ip

func BaseFunc(env *Env) {
	v := []string{
		"(define range (lambda (a b) (if (= a b) (quote ()) (cons a (range (+ a 1) b)))))",
		"(define (mod a b) (- a (* b ( // a b))))",
	}
	for _, i := range v {
		l := NewLexer(i)
		parser := NewParser(l)
		astNodes, _ := parser.Parse()
		for _, n := range astNodes {
			tmp := n.(Node)
			Eval(tmp, env)
		}
	}
}
