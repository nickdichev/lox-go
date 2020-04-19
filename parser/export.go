package parser

import (
	"github.com/ziyoung/lox-go/ast"
	"github.com/ziyoung/lox-go/lexer"
)

// helper functions for testing.

// ParseExpr parses expression.
func ParseExpr(input string) (expr ast.Expr, err error) {
	l := lexer.New(input)
	p := New(l)
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
	return p.parseExpression(), nil
}

// ParseStmts parses stamtements.
func ParseStmts(input string) ([]ast.Stmt, error) {
	l := lexer.New(input)
	p := New(l)
	return p.Parse()
}
