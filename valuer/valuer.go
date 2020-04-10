package valuer

import "strconv"

var typeMap = map[Type]string{
	NumberType:  "number",
	StringType:  "string",
	BooleanType: "bool",
	NilType:     "nil",
}

// Type represents type of Valuer.
type Type int

const (
	NumberType  Type = iota + 1 // number
	StringType                  // string
	BooleanType                 // bool
	NilType                     // nil
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
