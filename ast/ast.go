package ast

import (
	"fmt"
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

func (*Literal) node() {}

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

func (*BlockStmt) node()    {}
func (*ClassStmt) node()    {}
func (*ExprStmt) node()     {}
func (*FunctionStmt) node() {}
func (*IfStmt) node()       {}
func (*PrintStmt) node()    {}
func (*ReturnStmt) node()   {}
func (*VarStmt) node()      {}
func (*WhileStmt) node()    {}

// Ident represents an identifier.
type Ident struct {
	Name string
}

func (ident *Ident) String() string { return ident.Name }

type Literal struct {
	Token token.Token
	Value string
}

func (*Literal) expr() {}

func (lit *Literal) String() string {
	switch lit.Token {
	case token.Nil:
		return "null"
	case token.True:
		return "true"
	case token.False:
		return "false"
	case token.Number, token.String:
		return lit.Value
	}
	panic("unknown Literal.")
}

type (
	AssignExpr struct {
		Left  string
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
		Name   *Ident
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
		Name   *Ident
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
		Name string
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
	return fmt.Sprintf("%s(%s)", e.Callee, strings.Join(args, ", "))
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
	return e.Name
}

type (
	BlockStmt struct {
		Statements []Stmt
	}
	ClassStmt struct {
		Name       *Ident
		SuperClass VariableExpr
		Methods    []FunctionStmt
	}
	ExprStmt struct {
		Expression Expr
	}
	FunctionStmt struct {
		Name   string
		Params []*Ident
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
		Name        *Ident
		Initializer Expr
	}
	WhileStmt struct {
		Condition Expr
		Body      Stmt
	}
)

func (*BlockStmt) stmt()    {}
func (*ClassStmt) stmt()    {}
func (*ExprStmt) stmt()     {}
func (*FunctionStmt) stmt() {}
func (*IfStmt) stmt()       {}
func (*PrintStmt) stmt()    {}
func (*ReturnStmt) stmt()   {}
func (*VarStmt) stmt()      {}
func (*WhileStmt) stmt()    {}

func (s *BlockStmt) String() string {
	var sb strings.Builder
	sb.WriteString("{ ")
	for _, stmt := range s.Statements {
		sb.WriteString(stmt.String())
	}
	sb.WriteString(" }")
	return sb.String()
}

func (s *ClassStmt) String() string { return "TODO" }

func (s *ExprStmt) String() string {
	return s.Expression.String() + ";"
}

func (s *FunctionStmt) String() string {
	var sb strings.Builder
	sb.WriteString("fun ")
	sb.WriteString(s.Name)
	sb.WriteString("(")
	params := make([]string, len(s.Params))
	for i, p := range s.Params {
		params[i] = p.Name
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") { ")
	for _, stmt := range s.Body {
		sb.WriteString(stmt.String())
	}
	sb.WriteString(" }")
	return sb.String()
}

func (s *IfStmt) String() string {
	var sb strings.Builder
	sb.WriteString("if (")
	sb.WriteString(s.Condition.String())
	sb.WriteString(") ")
	sb.WriteString(s.ThenBranch.String())
	if s.ElseBranch != nil {
		sb.WriteString(" else ")
		sb.WriteString(s.ElseBranch.String())
	}
	return sb.String()
}

func (s *PrintStmt) String() string {
	var sb strings.Builder
	sb.WriteString("print ")
	sb.WriteString(s.Expression.String())
	sb.WriteRune(';')
	return sb.String()
}

func (s *ReturnStmt) String() string {
	str := "return"
	if s.Value != nil {
		str += " " + s.Value.String()
	}
	return str + ";"
}

func (s *VarStmt) String() string {
	var sb strings.Builder
	sb.WriteString("var ")
	sb.WriteString(s.Name.String())
	sb.WriteString(" = ")
	sb.WriteString(s.Initializer.String())
	sb.WriteRune(';')
	return sb.String()
}

func (s *WhileStmt) String() string {
	var sb strings.Builder
	sb.WriteString("while (")
	sb.WriteString(s.Condition.String())
	sb.WriteString(") ")
	sb.WriteString(s.Body.String())
	return sb.String()
}
