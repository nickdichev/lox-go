package main

import (
	"fmt"
	"os"

	"github.com/ziyoung/lox-go/cmd/lox/repl"
)

func main() {
	fmt.Fprintln(os.Stdout, "Lox programing language.")
	fmt.Fprintln(os.Stdout, "Feel free to type commands.")
	repl.Start(os.Stdin, os.Stdout)
}