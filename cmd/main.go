package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"ip"
	"log"
	"os"
)

func eval(str io.Reader, env *ip.Env) {

}
func repl(str io.Reader, env *ip.Env) {
	ip.BaseFunc(env)
	prompt := "zxt>"
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s", prompt)
		line, _ := in.ReadString('\n')
		l := ip.NewLexer(line)
		parser := ip.NewParser(l)
		astNodes, _ := parser.Parse()
		for _, n := range astNodes {
			tmp := n.(ip.Node)
			//fmt.Printf("tmp = %+v\n", tmp)
			v := ip.Eval(tmp, env)
			if v != nil {
				fmt.Printf("%v\n", v)
			}
		}

	}
}

func main() {
	//s := "(define take ;;this is comment ; (lambda (n seq) (if (<= n 0) (quote ()) (cons (car seq) (take (- n 1) (cdr seq))))))"
	//l := ip.NewLexer(s)
	//parser := ip.NewParser(l)
	//astNodes, _ := parser.Parse()
	//fmt.Printf("len = %v\n", len(astNodes))
	//var globalEnv = ip.StandardEnv()
	//for _, n := range astNodes {
	//	v := n.(ip.Node)
	//	if v.Type == ip.NodeTypeList {
	//		val := ip.Eval(v, globalEnv)
	//		fmt.Printf(" %+v\n", val)
	//	}
	//}
	isRepl := flag.Bool("repl", false, "Run as an interactive repl")
	flag.Parse()
	args := flag.Args()
	//default to repl if no files given
	if *isRepl || len(args) == 0 {
		repl(os.Stdin, ip.StandardEnv())
	} else {
		filepath := args[0]
		file, err := os.Open(filepath)
		if err != nil {
			log.Fatal("Error opening file to read!")
		}
		defer file.Close()
		env := ip.StandardEnv()
		source, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal("Error when source code to read!")
		}
		l := ip.NewLexer(string(source))
		parser := ip.NewParser(l)
		astNodes, _ := parser.Parse()
		for _, n := range astNodes {
			tmp := n.(ip.Node)
			//fmt.Printf("tmp = %+v\n", tmp)
			v := ip.Eval(tmp, env)
			if v != nil {
				fmt.Printf("%v\n", v)
			}
		}
	}
}
