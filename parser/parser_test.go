package parser

import (
	"testing"

	"github.com/ziyoung/lox-go/ast"
	"github.com/ziyoung/lox-go/lexer"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "123 - -456 - 789",
			expected: "((123 - (-456)) - 789)",
		},
		{
			input:    "1 >= 2 == !false",
			expected: "((1 >= 2) == (!false))",
		},
		{
			input:    "123 - 456 * 789 / 123",
			expected: "(123 - ((456 * 789) / 123))",
		},
	}
	for i, test := range tests {
		p := newParserFromInput(test.input)
		expr, err := parseExpression(p)
		if err != nil {
			t.Fatalf("test[%d] error occurs. error is %s", i, err.Error())
		}
		if expr.String() != test.expected {
			t.Fatalf("test[%d] expected expression is %q. got %q.", i, test.expected, expr.String())
		}
	}
}

func TestParseExpressionRecover(t *testing.T) {
	input := "123 + 456 -;123+456"
	expected := "(123 + 456)"
	p := newParserFromInput(input)
	expr, err := parseExpression(p)
	if err == nil {
		t.Fatalf("parser doesn't fail. get expression %s", expr)
	}

	p.synchronize()
	expr, err = parseExpression(p)
	if err != nil {
		t.Fatalf("error occurs. error is %s", err.Error())
	}
	if expr.String() != expected {
		t.Fatalf("expected expression is %q. got %q.", expected, expr.String())
	}
}

func TestParsePrintStament(t *testing.T) {
	input := `var a = 0;
a = a + 10;
var b = a = a + 100;
print a;`
	expected := []string{
		"var a = 0;",
		"a = (a + 10);",
		"var b = a = (a + 100);",
		"print a;",
	}
	p := newParserFromInput(input)
	statements, err := p.Parse()
	if err != nil {
		t.Fatalf("parse failed. error: %s", err.Error())
	}
	if len(statements) != len(expected) {
		t.Fatalf("length of staments should be %d. got %d", len(expected), len(statements))
	}
	for i, stmt := range statements {
		if stmt.String() != expected[i] {
			t.Errorf("test [%d]: expected text is %q. got %q", i, expected[i], stmt.String())
		}
	}
}

func newParserFromInput(input string) *Parser {
	l := lexer.New(input)
	return New(l)
}

func parseExpression(p *Parser) (expr ast.Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			if parseErr, ok := r.(parseError); ok {
				err = &parseErr
				expr = nil
			} else {
				panic(r)
			}
		}
	}()
	expr = p.parseExpression()
	return expr, nil
}
