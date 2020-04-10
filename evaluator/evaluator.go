package evaluator

import (
	"strconv"

	"github.com/ziyoung/lox-go/ast"
	"github.com/ziyoung/lox-go/token"
	"github.com/ziyoung/lox-go/valuer"
)

var (
	True  = &valuer.Boolean{Value: true}
	False = &valuer.Boolean{Value: false}
	Nil   = &valuer.Nil{}
)

func Eval(node ast.Node) valuer.Valuer {
	switch n := node.(type) {
	case *ast.Literal:
		return evalLiteral(n)
	case *ast.BinaryExpr:
		return evalBinaryExpr(n)
	case *ast.UnaryExpr:
		return evalUnaryExpr(n)
	case *ast.GroupingExpr:
		return Eval(n.Expression)
	}

	panic("unknown ast type.")
}

func evalLiteral(lit *ast.Literal) valuer.Valuer {
	switch lit.Token {
	case token.True:
		return True
	case token.False:
		return False
	case token.String:
		return &valuer.String{Value: lit.Value}
	case token.Number:
		v, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			panic(err)
		}
		return &valuer.Number{Value: v}
	case token.Nil:
		return &valuer.Nil{}
	}

	panic("unexpected literal.")
}

func evalBinaryExpr(expr *ast.BinaryExpr) valuer.Valuer {
	left := Eval(expr.Left)
	right := Eval(expr.Right)

	switch op := expr.Operator; op {
	case token.EqualEqual:
		t := isEqual(left, right)
		return toBooleanValuer(t)
	case token.BangEqual:
		t := !isEqual(left, right)
		return toBooleanValuer(t)
	case token.Greater:
		a, b := checkNumberOperands(op, left, right)
		t := a > b
		return toBooleanValuer(t)
	case token.GreaterEqual:
		a, b := checkNumberOperands(op, left, right)
		t := a >= b
		return toBooleanValuer(t)
	case token.Less:
		a, b := checkNumberOperands(op, left, right)
		t := a < b
		return toBooleanValuer(t)
	case token.LessEqual:
		a, b := checkNumberOperands(op, left, right)
		t := a <= b
		return toBooleanValuer(t)
	case token.Minus:
		a, b := checkNumberOperands(op, left, right)
		v := a - b
		return &valuer.Number{Value: v}
	case token.Plus:
		return doPlusOperation(left, right)
	case token.Slash:
		a, b := checkNumberOperands(op, left, right)
		if b == float64(0) {
			panic(runtimeError{
				token: op,
				s:     "Divisor can't be 0.",
			})
		}
		v := a / b
		return &valuer.Number{Value: v}
	case token.Star:
		a, b := checkNumberOperands(op, left, right)
		v := a * b
		return &valuer.Number{Value: v}
	}

	panic("unexpected binary expression.")
}

func evalUnaryExpr(expr *ast.UnaryExpr) valuer.Valuer {
	right := Eval(expr.Right)
	switch op := expr.Operator; op {
	case token.Bang:
		t := !isTruthy(right)
		return toBooleanValuer(t)
	case token.Minus:
		v := checkNumberOperand(op, right)
		return &valuer.Number{Value: -v}
	}

	panic("unexpected unary expression.")
}

func checkNumberOperand(operator token.Token, right valuer.Valuer) float64 {
	if a, ok := right.(*valuer.Number); ok {
		return a.Value
	}
	panic(runtimeError{
		token: operator,
		s:     "Operand must be a number.",
	})
}

func checkNumberOperands(operator token.Token, left, right valuer.Valuer) (float64, float64) {
	a, ok := left.(*valuer.Number)
	b, ok1 := right.(*valuer.Number)
	if !(ok && ok1) {
		panic(runtimeError{
			token: operator,
			s:     "Operands must be numbers.",
		})
	}
	return a.Value, b.Value
}

func doPlusOperation(left, right valuer.Valuer) valuer.Valuer {
	switch l := left.(type) {
	case *valuer.Number, *valuer.String:
		switch r := right.(type) {
		case *valuer.Number:
			if n, ok := l.(*valuer.Number); ok {
				return &valuer.Number{Value: n.Value + r.Value}
			}
			s, _ := l.(*valuer.String)
			return &valuer.String{
				Value: s.Value + r.String(),
			}
		case *valuer.String:
			if n, ok := l.(*valuer.Number); ok {
				return &valuer.String{
					Value: n.String() + r.Value,
				}
			}
			s, _ := l.(*valuer.String)
			return &valuer.String{Value: s.Value + r.Value}
		}
	}

	panic(runtimeError{
		token: token.Plus,
		s:     "Operands must be numbers or strings.",
	})
}

func isEqual(a, b valuer.Valuer) bool {
	_, ok := a.(*valuer.Boolean)
	_, ok1 := b.(*valuer.Boolean)
	if ok || ok1 {
		return isTruthy(a) == isTruthy(b)
	}

	switch a1 := a.(type) {
	case *valuer.Number:
		if b1, ok := b.(*valuer.Number); ok {
			return a1.Value == b1.Value
		}
	case *valuer.Nil:
		if _, ok := b.(*valuer.Nil); ok {
			return true
		}
	case *valuer.String:
		if b1, ok := b.(*valuer.String); ok {
			return a1.Value == b1.Value
		}
	}
	return false
}

func isTruthy(value valuer.Valuer) bool {
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case *valuer.Boolean:
		return v.Value
	case *valuer.Number:
		return v.Value != float64(0)
	case *valuer.Nil:
		return false
	case *valuer.String:
		return v.Value != ""
	}
	return false
}

func toBooleanValuer(t bool) *valuer.Boolean {
	if t {
		return True
	}
	return False
}
