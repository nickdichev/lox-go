package token

import "testing"

func TestIdentIsKeyword(t *testing.T) {
	tests := []struct {
		literal string
		tok     Token
	}{
		{"abc", Identifier},
		{"and", And},
		{"class", Class},
		{"else", Else},
		{"false", False},
		{"fun", Fun},
		{"for", For},
		{"if", If},
		{"nil", Nil},
		{"or", Or},
		{"print", Print},
		{"return", Return},
		{"super", Super},
		{"this", This},
		{"true", True},
		{"var", Var},
		{"while", While},
	}

	for _, test := range tests {
		if tok := Lookup(test.literal); tok != test.tok {
			t.Errorf("token of %s expected to be %s. got %s", test.literal, test.tok, tok)
		}
	}
}
