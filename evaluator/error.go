package evaluator

import "github.com/ziyoung/lox-go/token"

type runtimeError struct {
	s     string
	token token.Token
}

func (r *runtimeError) Error() string { return r.s }
