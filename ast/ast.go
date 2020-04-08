package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ziyoung/lox-go/token"
)

// Node represents a node in Lox.
// All node types should implement the Node interface.
type Node interface {
	node()
	String() string
}

// Expr represents an expression that can be evaluated to a value.
type Expr interface {
	Node
	expr()
}

// Stmt represent a statement in Lox.
type Stmt interface {
	Node
	stmt()
}

func (*Ident) node() {}

func (*NumberLiteral) node() {}
func (*StringLiteral) node() {}
func (*BoolLiteral) node()   {}
func (*NullLiteral) node()   {}

func (*AssignExpr) node()   {}
func (*BinaryExpr) node()   {}
func (*CallExpr) node()     {}
func (*GetExpr) node()      {}
func (*GroupingExpr) node() {}
func (*LogicalExpr) node()  {}
func (*SetExpr) node()      {}
func (*SuperExpr) node()    {}
func (*ThisExpr) node()     {}
func (*UnaryExpr) node()    {}
func (*VariableExpr) node() {}

func (*BlockStmt) node()      {}
func (*ClassStmt) node()      {}
func (*ExpressionStmt) node() {}
func (*FunctionStmt) node()   {}
func (*IfStmt) node()         {}
func (*PrintStmt) node()      {}
func (*ReturnStmt) node()     {}
func (*VarStmt) node()        {}
func (*WhileStmt) node()      {}

// Ident represents an identifier.
type Ident struct {
	Name string
}

func (ident *Ident) String() string { return ident.Name }

type (
	NumberLiteral struct {
		Value float64
	}
	StringLiteral struct {
		Value string
	}
	BoolLiteral struct {
		Value bool
	}
	NullLiteral struct{}
)

func (*NumberLiteral) expr() {}
func (*StringLiteral) expr() {}
func (*BoolLiteral) expr()   {}
func (*NullLiteral) expr()   {}

func (lit *NumberLiteral) String() string {
	return strconv.FormatFloat(lit.Value, 'f', -1, 64)
}

func (lit *StringLiteral) String() string {
	return lit.Value
}

func (lit *BoolLiteral) String() string {
	return strconv.FormatBool(lit.Value)
}

func (lit *NullLiteral) String() string {
	return "null"
}

type (
	AssignExpr struct {
		Left  Ident
		Value Expr
	}
	BinaryExpr struct {
		Left     Expr
		Operator token.Token
		Right    Expr
	}
	CallExpr struct {
		Callee    Expr
		Arguments []Expr
	}
	GetExpr struct {
		Object Expr
		Name   Ident
	}
	GroupingExpr struct {
		Expression Expr
	}
	LogicalExpr struct {
		Left     Expr
		Operator token.Token
		Right    Expr
	}
	SetExpr struct {
		Object Expr
		Name   Ident
		Value  Expr
	}
	SuperExpr struct {
		// Method  Ident
		Keyword token.Token
		Method  token.Token
	}
	ThisExpr struct {
		Keyword token.Token
	}
	UnaryExpr struct {
		Operator token.Token
		Right    Expr
	}
	VariableExpr struct {
		Name token.Token
	}
)

func (*AssignExpr) expr()   {}
func (*BinaryExpr) expr()   {}
func (*CallExpr) expr()     {}
func (*GetExpr) expr()      {}
func (*GroupingExpr) expr() {}
func (*LogicalExpr) expr()  {}
func (*SetExpr) expr()      {}
func (*SuperExpr) expr()    {}
func (*ThisExpr) expr()     {}
func (*UnaryExpr) expr()    {}
func (*VariableExpr) expr() {}

func (e *AssignExpr) String() string {
	return fmt.Sprintf("%s = %s", e.Left, e.Value)
}

func (e *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Operator, e.Right)
}

func (e *CallExpr) String() string {
	args := make([]string, len(e.Arguments))
	for i, arg := range e.Arguments {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", e.Callee, strings.Join(args, " ,"))
}

func (e *GetExpr) String() string {
	return "TODO"
}

func (e *GroupingExpr) String() string {
	return fmt.Sprintf("(%s)", e.Expression)
}

func (e *LogicalExpr) String() string {
	return fmt.Sprintf("%s %s %s", e.Left, e.Operator, e.Right)
}

func (e *SetExpr) String() string {
	return "TODO"
}

func (e *SuperExpr) String() string {
	return "TODO"
}

func (e *ThisExpr) String() string {
	return "TODO"
}

func (e *UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator, e.Right)
}

func (e *VariableExpr) String() string {
	return e.Name.String()
}

type (
	BlockStmt struct {
		Statements []Stmt
	}
	ClassStmt struct {
		Name       Ident
		SuperClass VariableExpr
		Methods    []FunctionStmt
	}
	ExpressionStmt struct {
		Expression Expr
	}
	FunctionStmt struct {
		Name   Ident
		Params []Ident
		Body   []Stmt
	}
	IfStmt struct {
		Condition  Expr
		ThenBranch Stmt
		ElseBranch Stmt
	}
	PrintStmt struct {
		Expression Expr
	}
	ReturnStmt struct {
		Keyword token.Token
		Value   Expr
	}
	VarStmt struct {
		Name        Ident
		Initializer Expr
	}
	WhileStmt struct {
		Condition Expr
		Body      Stmt
	}
)

func (*BlockStmt) stmt()      {}
func (*ClassStmt) stmt()      {}
func (*ExpressionStmt) stmt() {}
func (*FunctionStmt) stmt()   {}
func (*IfStmt) stmt()         {}
func (*PrintStmt) stmt()      {}
func (*ReturnStmt) stmt()     {}
func (*VarStmt) stmt()        {}
func (*WhileStmt) stmt()      {}

func (s *BlockStmt) String() string {
	var sb strings.Builder
	sb.WriteRune('{')
	for _, stmt := range s.Statements {
		sb.WriteString(stmt.String())
	}
	sb.WriteRune('}')
	return sb.String()
}

func (s *ClassStmt) String() string      { return "TODO" }
func (s *ExpressionStmt) String() string { return "TODO" }
func (s *FunctionStmt) String() string   { return "TODO" }
func (s *IfStmt) String() string         { return "TODO" }
func (s *PrintStmt) String() string      { return "TODO" }
func (s *ReturnStmt) String() string     { return "TODO" }
func (s *VarStmt) String() string        { return "TODO" }
func (s *WhileStmt) String() string      { return "TODO" }
