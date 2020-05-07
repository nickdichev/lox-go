package errors

import (
	"github.com/ziyoung/lox-go/token"
)

type RuntimeError struct {
	s     string
	token token.Token
}

func (r *RuntimeError) Error() string {
	return r.s
}

// Error throws runtime error.
func Error(token token.Token, s string) {
	panic(RuntimeError{token: token, s: s})
}
