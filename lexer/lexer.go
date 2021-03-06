package lexer

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/ziyoung/lox-go/token"
)

var eof = rune(-1)

var (
	// identifer error
	errUnterminated = errors.New("unterminated string")
	errEspace       = errors.New("invalid escape char")
	errInvalidChar  = errors.New("invalid unicode char")

	// number error
	errLessPower = errors.New("power is required")
)

// Lexer represents a lexical scanner for Lox programing language.
type Lexer struct {
	s      *scanner.Scanner
	ch     rune
	tokBuf *strings.Builder
}

func (l *Lexer) consume() {
	if l.isAtEnd() {
		return
	}
	ch := l.s.Next()
	if ch == scanner.EOF {
		l.ch = eof
		return
	}
	l.ch = ch
}

func (l *Lexer) peek() rune {
	ch := l.s.Peek()
	return ch
}

func (l *Lexer) skip() {
	for unicode.IsSpace(l.ch) {
		l.consume()
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.ch == eof
}

func (l *Lexer) match(ch rune) bool {
	l.consume()
	if l.isAtEnd() || l.ch != ch {
		return false
	}
	l.consume()
	return true
}

func (l *Lexer) error(msg string) {
	// FIXME: position is wrong
	fmt.Fprintf(os.Stderr, "%s %s\n", l.Pos().String(), msg)
}

func (l *Lexer) readIdentifier() string {
	l.tokBuf.Reset()
	for isAlphaNumeric(l.ch) {
		l.tokBuf.WriteRune(l.ch)
		l.consume()
	}
	return l.tokBuf.String()
}

func (l *Lexer) readString() (string, error) {
	l.tokBuf.Reset()
	l.consume()
	if l.ch == '"' {
		l.consume()
		return "", nil
	}

	for l.ch != '"' {
		if l.isAtEnd() {
			l.error(errUnterminated.Error())
			return "", errUnterminated
		} else if l.ch == '\\' {
			peekCh := l.peek()
			if peekCh == eof {
				l.error(errEspace.Error())
				return "", errEspace
			}
			l.consume()
			switch peekCh {
			case '"':
				l.tokBuf.WriteRune('"')
			case 'u':
				code := make([]rune, 4)
				for i := range code {
					l.consume()
					if !unicode.Is(unicode.Hex_Digit, l.ch) {
						l.error(errInvalidChar.Error())
						return "", errInvalidChar
					}
					code[i] = l.ch
				}
				l.tokBuf.WriteRune(charCode2Rune(string(code)))
			}
		} else {
			l.tokBuf.WriteRune(l.ch)
		}
		l.consume()
	}
	// end ".
	l.consume()
	return l.tokBuf.String(), nil
}

func (l *Lexer) readNumber() (string, error) {

	l.tokBuf.Reset()
	for unicode.IsNumber(l.ch) {
		l.tokBuf.WriteRune(l.ch)
		l.consume()
	}

	if l.ch == '.' {
		if !unicode.IsNumber(l.peek()) {
			return l.tokBuf.String(), nil
		}
		l.tokBuf.WriteRune(l.ch)
		l.consume()
		for unicode.IsNumber(l.ch) {
			l.tokBuf.WriteRune(l.ch)
			l.consume()
		}
	}

	if l.ch == 'E' || l.ch == 'e' {
		seenPower := false
		l.tokBuf.WriteRune(l.ch)
		l.consume()
		if l.ch == '+' || l.ch == '-' {
			l.tokBuf.WriteRune(l.ch)
			l.consume()
		}
		for unicode.IsNumber(l.ch) {
			seenPower = true
			l.tokBuf.WriteRune(l.ch)
			l.consume()
		}
		if !seenPower {
			l.error(errLessPower.Error())
			return "", errLessPower
		}
	}

	return l.tokBuf.String(), nil
}

// NextToken reads and returns token and literal.
// It returns token.Illegal for invalid string or number.
// It return token.EOF at the end of input string.
func (l *Lexer) NextToken() (tok token.Token, literal string) {
	l.skip()

	switch l.ch {
	case '(':
		tok = token.LeftParen
		literal = "("
	case ')':
		tok = token.RightParen
		literal = ")"
	case '{':
		tok = token.LeftBrace
		literal = "{"
	case '}':
		tok = token.RightBrace
		literal = "}"
	case ',':
		tok = token.Comma
		literal = ","
	case '.':
		tok = token.Dot
		literal = "."
	case '-':
		tok = token.Minus
		literal = "-"
	case '+':
		tok = token.Plus
		literal = "+"
	case ';':
		tok = token.Semicolon
		literal = ";"
	case '/':
		tok = token.Slash
		literal = "/"
	case '*':
		tok = token.Star
		literal = "*"
	case '!':
		if l.match('=') {
			tok = token.BangEqual
			literal = "!="
		} else {
			tok = token.Bang
			literal = "!"
		}
		return
	case '=':
		if l.match('=') {
			tok = token.EqualEqual
			literal = "=="
		} else {
			tok = token.Equal
			literal = "="
		}
		return
	case '>':
		if l.match('=') {
			tok = token.GreaterEqual
			literal = ">="
		} else {
			tok = token.Greater
			literal = ">"
		}
		return
	case '<':
		if l.match('=') {
			tok = token.LessEqual
			literal = "<="
		} else {
			tok = token.Less
			literal = "<"
		}
		return
	case '"':
		liter, err := l.readString()
		if err != nil {
			return token.Illegal, liter
		}
		tok = token.String
		literal = liter
		return
	case eof:
		tok = token.EOF
		return
	default:
		if unicode.IsLetter(l.ch) {
			literal = l.readIdentifier()
			tok = token.Lookup(literal)
			return
		} else if unicode.IsNumber(l.ch) {
			liter, err := l.readNumber()
			if err != nil {
				return token.Illegal, ""
			}
			tok = token.Number
			literal = liter
			return
		}

		tok = token.Illegal
		literal = ""
	}

	l.consume()
	return
}

// Pos returns current position of lexer.
func (l *Lexer) Pos() scanner.Position {
	return l.s.Pos()
}

func charCode2Rune(code string) rune {
	v, err := strconv.ParseInt(code, 16, 32)
	if err != nil {
		return unicode.ReplacementChar
	}
	return rune(v)
}

func isAlphaNumeric(ch rune) bool { return unicode.IsLetter(ch) || unicode.IsNumber(ch) || ch == '_' }

// New return an instance of Lexer.
func New(input string) *Lexer {
	s := &scanner.Scanner{}
	s.Init(strings.NewReader(input))
	l := &Lexer{
		s:      s,
		tokBuf: &strings.Builder{},
	}
	l.consume()
	return l
}
