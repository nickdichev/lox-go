package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ziyoung/lox-go/ast"
	"github.com/ziyoung/lox-go/lexer"
	"github.com/ziyoung/lox-go/token"
)

type item struct {
	tok token.Token
	lit string
}

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

// func (p *Parser) peekToken() token.Token {
// 	p.peekCount++
// 	tok, lit := p.l.NextToken()
// 	p.items[p.peekCount] = item{tok, lit}
// 	return tok
// }

func (p *Parser) parseExpression() ast.Expr {
	return p.parseEquality()
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
	case token.True:
	case token.False:
		expr = &ast.BoolLiteral{
			Value: tok == token.True,
		}
	case token.Nil:
		expr = &ast.NullLiteral{}
	case token.String:
		expr = &ast.StringLiteral{
			Value: lit,
		}
	case token.Number:
		v, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			v = 0
			// FIXME: skip error now.
			fmt.Println(err)
		}
		expr = &ast.NumberLiteral{
			Value: v,
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
	s := fmt.Sprintf("%s :%s", p.l.Pos(), msg)
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
