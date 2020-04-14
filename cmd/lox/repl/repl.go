package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ziyoung/lox-go/interpreter"
	"github.com/ziyoung/lox-go/lexer"
	"github.com/ziyoung/lox-go/parser"
)

const prompt = ">> "

// Start creates a REPL for Lox.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		l := lexer.New(line)
		p := parser.New(l)
		statements := p.Parse()
		interpreter.Interpret(statements)
	}
}
