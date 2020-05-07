package resolver

import (
	"fmt"

	"github.com/ziyoung/lox-go/errors"
	"github.com/ziyoung/lox-go/token"
)

// Scopes represents variable scopes.
type Scopes []map[string]bool

func (s *Scopes) check(name string) (exist bool, init bool) {
	if !s.isEmpty() {
		scope := s.peek()
		if _, ok := scope[name]; ok {
			return true, ok
		}
	}
	return false, false
}

func (s *Scopes) begin() {
	scope := make(map[string]bool)
	*s = append(*s, scope)
}

func (s *Scopes) end() {
	s.pop()
}

func (s Scopes) peek() map[string]bool {
	if s.isEmpty() {
		panic("scope peek error: empty scopes")
	}
	return s[len(s)-1]
}

func (s *Scopes) pop() {
	if s.isEmpty() {
		panic("scope pop error: empty scopes")
	}
	*s = (*s)[:len(*s)-1]
}

func (s Scopes) isEmpty() bool {
	return len(s) == 0
}

func (s Scopes) declare(name string) {
	if s.isEmpty() {
		return
	}
	scope := s.peek()
	if _, ok := scope[name]; ok {
		errors.Error(token.Var, fmt.Sprintf("variable name %q has been already delcared in this scope.", name))
	}
	scope[name] = false
}

func (s Scopes) define(name string) {
	if s.isEmpty() {
		return
	}
	scope := s.peek()
	scope[name] = true
}

// NewScopes returns Scopes instance.
func NewScopes() Scopes {
	scopes := make([]map[string]bool, 0)
	return scopes
}
