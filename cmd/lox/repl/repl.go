package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ziyoung/lox-go/interpreter"
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
		expr, err := parser.ParseExpr(line)
		if err != nil {
			fmt.Fprintln(out, err.Error())
			continue
		}
		v := interpreter.Eval(expr)
		if v != nil {
			fmt.Fprintf(out, "\033[1;30m%s\033[0m %s\n", v.Type(), v)
		}
	}
}
