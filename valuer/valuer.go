package valuer

import (
	"strconv"

	"github.com/ziyoung/lox-go/ast"
)

var typeMap = map[Type]string{
	NumberType:   "number",
	StringType:   "string",
	BooleanType:  "bool",
	NilType:      "nil",
	FunctionType: "function",
	ReturnType:   "return",
}

// Type represents type of Valuer.
type Type int

const (
	NumberType   Type = iota + 1 // number
	StringType                   // string
	BooleanType                  // bool
	NilType                      // nil
	FunctionType                 // function
	ReturnType                   // return
)

func (typ Type) String() string {
	if s, ok := typeMap[typ]; ok {
		return s
	}
	return "unknown"
}

type Valuer interface {
	Type() Type
	String() string
}

type Number struct {
	Value float64
}

// Type returns its Type.
func (*Number) Type() Type { return NumberType }

func (num *Number) String() string {
	return strconv.FormatFloat(num.Value, 'f', -1, 64)
}

type String struct {
	Value string
}

// Type returns its Type.
func (*String) Type() Type { return StringType }

func (s *String) String() string { return s.Value }

type Boolean struct {
	Value bool
}

// Type returns its Type.
func (*Boolean) Type() Type { return BooleanType }

func (b *Boolean) String() string { return strconv.FormatBool(b.Value) }

type Nil struct{}

// Type returns its Type.
func (*Nil) Type() Type { return NilType }

func (*Nil) String() string { return "nil" }

type Function struct {
	Name    string
	Params  []*ast.Ident
	Body    []ast.Stmt
	Closure *Environment
}

// Type returns its Type.
func (*Function) Type() Type { return FunctionType }

func (fn *Function) String() string {
	return "<fn " + fn.Name + ">"
}

// Arity returns size of params.
func (fn *Function) Arity() int {
	return len(fn.Params)
}

type ReturnValue struct {
	Value Valuer
}

// Type returns its Type.
func (*ReturnValue) Type() Type { return ReturnType }

func (rt *ReturnValue) String() string {
	return rt.Value.String()
}
