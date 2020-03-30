package token

// Token is a lexical token of lox programing language.
type Token int

const (
	Illegal Token = iota
	EOF
	Comment

	// single-character

	LeftParen  // (
	RightParen // )
	LeftBrace  // {
	RightBrace // }
	Comma      // ,
	Dot        // .
	Minus      // -
	Plus       // +
	Semicolon  // ;
	Slash      // /
	Star       // *

	Bang         // !
	BangEqual    // !=
	Equal        // =
	EqualEqual   // ==
	Greater      // >
	GreaterEqual // >=
	Less         // <
	LessEqual    // <=

	Identifier // abc
	String     // "abc"
	Number     // 123

	keywordBegin

	And    // and
	Class  // class
	Else   // else
	False  // false
	Fun    // fun
	For    // for
	If     // if
	Nil    // nil
	Or     // or
	Print  // print
	Return // return
	Super  // super
	This   // this
	True   // true
	Var    // var
	While  // while

	keywordEnd
)

var tokens = [...]string{
	Illegal:      "illegal",
	EOF:          "EOF",
	Comment:      "comment",
	LeftParen:    "(",
	RightParen:   ")",
	LeftBrace:    "{",
	RightBrace:   "}",
	Comma:        ",",
	Dot:          ".",
	Minus:        "-",
	Plus:         "+",
	Semicolon:    ";",
	Slash:        "/",
	Star:         "*",
	Bang:         "!",
	BangEqual:    "!=",
	Equal:        "=",
	EqualEqual:   "==",
	Greater:      ">",
	GreaterEqual: ">=",
	Less:         "<",
	LessEqual:    "<=",
	Identifier:   "identifier",
	String:       "string",
	Number:       "number",
	And:          "and",
	Class:        "class",
	Else:         "else",
	False:        "false",
	Fun:          "fun",
	For:          "for",
	If:           "if",
	Nil:          "nil",
	Or:           "or",
	Print:        "print",
	Return:       "return",
	Super:        "super",
	This:         "this",
	True:         "true",
	Var:          "var",
	While:        "while",
}

var keywords = map[string]Token{}

func init() {
	i, j := int(keywordBegin)+1, int(keywordEnd)
	for ; i < j; i++ {
		keywords[tokens[i]] = Token(i)
	}
}

func (tok Token) String() string {
	i := int(tok)
	if i > len(tokens)-1 {
		return ""
	}
	return tokens[i]
}

// Lookup returns the token type associated with a given string.
func Lookup(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Identifier
}
