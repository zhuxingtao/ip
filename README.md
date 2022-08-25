#  ip

ip is a scheme interpreter  implemented by golang, 
inspired by

    - https://interpreterbook.com/
    - http://norvig.com/lispy.html

### Example
 -  go run cmd/main.go qsort.scm 
 -  go run cmd/main.go t.scm
 -  go run cmd/main.go --repl=1

### Features:
 - Interactive REPL shell
 - Type: String, Number, Quote, LambdaProcess, List, Bool ...
 - Call with current continuation
 - Hand-Written lexer, parser  
