package lexer

import (
	"testing"

	"github.com/ziyoung/lox-go/token"
)

func TestReadSimpleToken(t *testing.T) {
	input := `
() {}
, . - + ;
/ * !
= == !=
> >=
< <=`
	l := New(input)
	tests := []struct {
		expectTok     token.Token
		expectLiteral string
	}{
		{token.LeftParen, "("},
		{token.RightParen, ")"},
		{token.LeftBrace, "{"},
		{token.RightBrace, "}"},
		{token.Comma, ","},

		{token.Dot, "."},
		{token.Minus, "-"},
		{token.Plus, "+"},
		{token.Semicolon, ";"},
		{token.Slash, "/"},

		{token.Star, "*"},
		{token.Bang, "!"},
		{token.Equal, "="},
		{token.EqualEqual, "=="},
		{token.BangEqual, "!="},

		{token.Greater, ">"},
		{token.GreaterEqual, ">="},
		{token.Less, "<"},
		{token.LessEqual, "<="},
	}

	for i, test := range tests {
		tok, literal := l.NextToken()
		if test.expectTok != tok {
			t.Fatalf("test [%d]: expected token is %s. got %s", i, test.expectTok, tok)
		}

		if test.expectLiteral != literal {
			t.Fatalf("test [%d]: expected literal is %s. got %s", i, test.expectLiteral, literal)
		}
	}

	tok, _ := l.NextToken()
	if token.EOF != tok {
		t.Fatalf("expected token is EOF. got %s", tok)
	}
}

func TestReadString(t *testing.T) {
	input := `"abc xyz \u5b57符串 lox\u8Bed言"`
	expected := "abc xyz 字符串 lox语言"
	l := New(input)
	tok, literal := l.NextToken()

	if tok != token.String {
		t.Fatalf("expected token is string. got %s", tok)
	}

	if literal != expected {
		t.Fatalf("expected literal is %q. got %q", expected, literal)
	}
}

func TestReadInvalidString(t *testing.T) {
	tests := []string{
		`"abc`,
		`"\"`,
		`"\u"`,
		`"\udef"`,
	}

	for i, test := range tests {
		l := New(test)
		tok, literal := l.NextToken()

		if tok != token.Illegal {
			t.Fatalf("test [%d]: expected token is illegal. got %s", i, tok)
		}

		if literal != "" {
			t.Fatalf("test [%d]: expected literal is empty string. got %q", i, literal)
		}
	}
}

func TestReadIdentifier(t *testing.T) {
	input := `
abc 		xyz 		a123 		A_123			X_x_
and 		class 	else 		false 		fun
for 		if 			nil 		or 				print
return 	super 	this 		true			var
while
`
	tests := []struct {
		expectTok     token.Token
		expectLiteral string
	}{
		{token.Identifier, "abc"},
		{token.Identifier, "xyz"},
		{token.Identifier, "a123"},
		{token.Identifier, "A_123"},
		{token.Identifier, "X_x_"},

		{token.And, "and"},
		{token.Class, "class"},
		{token.Else, "else"},
		{token.False, "false"},
		{token.Fun, "fun"},

		{token.For, "for"},
		{token.If, "if"},
		{token.Nil, "nil"},
		{token.Or, "or"},
		{token.Print, "print"},

		{token.Return, "return"},
		{token.Super, "super"},
		{token.This, "this"},
		{token.True, "true"},
		{token.Var, "var"},

		{token.While, "while"},
	}
	l := New(input)

	for i, test := range tests {
		tok, literal := l.NextToken()

		if tok != test.expectTok {
			t.Fatalf("test [%d]: expected token is %s. got %s", i, test.expectTok, tok)
		}

		if literal != test.expectLiteral {
			t.Fatalf("test [%d]: expected literal is %s. got %s", i, test.expectLiteral, literal)
		}
	}
}

func TestReadNumber(t *testing.T) {
	input := `
0
01
0123
123
123.0

0.123
123.456
1E1
01E123
123E0

123E+1
123E-1
123.45e1
123.45e+1
123.45e-1`
	l := New(input)
	tests := []string{
		"0", "01", "0123", "123", "123.0",
		"0.123", "123.456", "1E1", "01E123", "123E0",
		"123E+1", "123E-1", "123.45e1", "123.45e+1", "123.45e-1",
	}

	for i, test := range tests {
		tok, literal := l.NextToken()
		if tok != token.Number {
			t.Fatalf("test [%d]: expected token is number. got %s", i, tok)
		}

		if literal != test {
			t.Fatalf("test [%d]: expected literal is %q. got %q", i, test, literal)
		}
	}
}

func TestReadInvalidNumber(t *testing.T) {
	tests := []string{
		"123E",
		"123E.",
		"123.45e-",
		"123.45e+",
	}

	for i, test := range tests {
		l := New(test)
		tok, literal := l.NextToken()

		if tok != token.Illegal {
			t.Fatalf("test [%d]: expected token is illegal. got %s", i, tok)
		}

		if literal != "" {
			t.Fatalf("test [%d]: expected literal is empty string. got %q", i, literal)
		}
	}
}
