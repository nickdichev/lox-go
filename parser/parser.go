package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/ziyoung/lox-go/ast"
	"github.com/ziyoung/lox-go/lexer"
	"github.com/ziyoung/lox-go/token"
)

// Parser represents the lox parser.
type Parser struct {
	l *lexer.Lexer

	tok token.Token
	lit string

	trace  bool
	indent int
}

func (p *Parser) nextToken() token.Token {
	if p.isAtEnd() {
		return token.EOF
	}
	tok, lit := p.l.NextToken()
	p.tok = tok
	p.lit = lit
	return tok
}

// Parse returns all statements of input.
func (p *Parser) Parse() []ast.Stmt {
	statements := make([]ast.Stmt, 0)
	for !p.isAtEnd() {
		stmt := p.parseDeclaration()
		statements = append(statements, stmt)
	}
	return statements
}

func (p *Parser) parseDeclaration() (stmt ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(parseError); ok {
				stmt = nil
				p.synchronize()
			} else {
				panic(r)
			}
		}
	}()
	if p.match(token.Var) {
		stmt = p.parseVarDeclaration()
		return
	}
	stmt = p.parseStatement()
	return
}

func (p *Parser) parseVarDeclaration() *ast.VarStmt {
	name := p.lit
	p.expect(token.Identifier, "Expect variable name.")
	var stmt = &ast.VarStmt{
		Name: &ast.Ident{
			Name: name,
		},
	}
	var initializer ast.Expr
	if p.match(token.Equal) {
		initializer = p.parseExpression()
	}
	p.expect(token.Semicolon, "Expect ';' after variable declaration.")
	stmt.Initializer = initializer
	return stmt
}

func (p *Parser) parseStatement() ast.Stmt {
	if p.match(token.Print) {
		return p.parsePrintStatement()
	}
	if p.match(token.LeftBrace) {
		return p.parseBlockStatement()
	}
	return p.parseExprStatement()
}

func (p *Parser) parsePrintStatement() ast.Stmt {
	expr := p.parseExpression()
	p.expect(token.Semicolon, "Expect ';' after value.")
	return &ast.PrintStmt{
		Expression: expr,
	}
}

func (p *Parser) parseBlockStatement() ast.Stmt {
	statements := make([]ast.Stmt, 0)
	for !(p.check(token.RightBrace) || p.isAtEnd()) {
		statements = append(statements, p.parseDeclaration())
	}
	p.expect(token.RightBrace, "Expect '}' after block.")
	return &ast.BlockStmt{
		Statements: statements,
	}
}

func (p *Parser) parseExprStatement() ast.Stmt {
	expr := p.parseExpression()
	p.expect(token.Semicolon, "Expect ';' after expression.")
	return &ast.ExprStmt{
		Expression: expr,
	}
}

func (p *Parser) parseExpression() ast.Expr {
	return p.parseAssignment()
}

func (p *Parser) parseAssignment() ast.Expr {
	expr := p.parseEquality()
	if p.match(token.Equal) {
		v := p.parseAssignment()
		if variable, ok := expr.(*ast.VariableExpr); ok {
			return &ast.AssignExpr{
				Left:  variable.Name,
				Value: v,
			}
		}
		p.error("Invalid assignment target.")
	}
	return expr
}

func (p *Parser) parseEquality() ast.Expr {
	expr := p.parseComparison()
	operator := p.tok
	for p.match(token.EqualEqual, token.BangEqual) {
		right := p.parseComparison()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
		operator = p.tok
	}
	return expr
}

func (p *Parser) parseComparison() ast.Expr {
	expr := p.parseAddition()
	operator := p.tok
	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		right := p.parseAddition()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
		operator = p.tok
	}
	return expr
}

func (p *Parser) parseAddition() ast.Expr {
	expr := p.parseMultiplacation()
	operator := p.tok
	for p.match(token.Plus, token.Minus) {
		right := p.parseMultiplacation()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
		operator = p.tok
	}
	return expr
}

func (p *Parser) parseMultiplacation() ast.Expr {
	expr := p.parseUnary()
	operator := p.tok
	for p.match(token.Slash, token.Star) {
		right := p.parseUnary()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
		operator = p.tok
	}
	return expr
}

func (p *Parser) parseUnary() ast.Expr {
	operator := p.tok
	if p.match(token.Bang, token.Minus) {
		right := p.parseUnary()
		return &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (expr ast.Expr) {
	tok, lit := p.tok, p.lit
	switch tok {
	case token.True, token.False, token.Nil, token.String, token.Number:
		expr = &ast.Literal{
			Token: tok,
			Value: lit,
		}
	case token.Identifier:
		expr = &ast.VariableExpr{
			Name: lit,
		}
	case token.LeftParen:
		p.nextToken()
		inner := p.parseExpression()
		p.expect(token.RightParen, "Expect ) after expression.")
		expr = &ast.GroupingExpr{
			Expression: inner,
		}
		return
	default:
		p.error("Expected expression.")
	}
	p.nextToken()
	return expr
}

func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		switch p.tok {
		case token.Semicolon:
			p.nextToken()
			return
		case token.Class, token.Fun, token.Var, token.If, token.While, token.Print, token.Return:
			return
		}
		p.nextToken()
	}
}

func (p *Parser) match(tokens ...token.Token) bool {
	for _, tok := range tokens {
		if p.check(tok) {
			p.nextToken()
			return true
		}
	}
	return false
}

func (p *Parser) expect(tok token.Token, msg string) {
	if p.check(tok) {
		p.nextToken()
		return
	}
	p.error(msg)
}

func (p *Parser) error(msg string) {
	s := fmt.Sprintf("%s %s", p.l.Pos(), msg)
	fmt.Fprintln(os.Stderr, s)
	panic(parseError{s})
}

func (p *Parser) check(tok token.Token) bool {
	if p.isAtEnd() {
		return false
	}
	return p.tok == tok
}

func (p *Parser) isAtEnd() bool {
	return p.tok == token.EOF
}

// New returns a parser instance.
func New(l *lexer.Lexer) *Parser {
	parser := &Parser{
		l: l,
	}
	parser.nextToken()
	return parser
}

// trace
func identLevel(count int) string {
	count--
	if count < 0 {
		count = 0
	}
	return strings.Repeat("\t", count)
}

func trace(p *Parser, msg string) (*Parser, string) {
	fmt.Printf("%sBEGIN %s\n", identLevel(p.indent), msg)
	p.indent++
	return p, msg
}

func unTrace(p *Parser, msg string) *Parser {
	count := p.indent
	if count < 0 {
		count = 0
	}
	fmt.Printf("%sEND %s\n", identLevel(p.indent), msg)
	p.indent--
	return p
}

// ParseExpr parses expression. If Parse function is finished, it should be deleted.
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
