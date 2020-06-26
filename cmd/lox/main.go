package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ziyoung/lox-go/cmd/lox/repl"
	"github.com/ziyoung/lox-go/interpreter"
	"github.com/ziyoung/lox-go/lexer"
	"github.com/ziyoung/lox-go/parser"
)

func main() {
	if len(os.Args) >= 2 {
		name := os.Args[1]
		b, err := ioutil.ReadFile(name)
		if err != nil {
			panic(err)
		}
		l := lexer.New(string(b))
		p := parser.New(l)
		if statements, err := p.Parse(); err == nil && len(statements) != 0 {
			interpreter.Interpret(statements)
		}
		return
	}

	fmt.Fprintln(os.Stdout, "Lox programing language.")
	fmt.Fprintln(os.Stdout, "Feel free to type commands.")
	fmt.Fprintln(os.Stdout, "Type \"exit\" to exit.")
	repl.Start(os.Stdin, os.Stdout)
}
