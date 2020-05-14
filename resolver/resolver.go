package resolver

import (
	"github.com/ziyoung/lox-go/ast"
	"github.com/ziyoung/lox-go/errors"
	"github.com/ziyoung/lox-go/token"
)

type functionType int

const (
	FunctionNone functionType = iota // not return
	Function
	Initializer
	Method
)

type classType int

const (
	ClassNone classType = iota
	Class
)

var (
	scopes          = NewScopes()
	curFunctionType = FunctionNone
	curClassType    = ClassNone
)

func Resolve(node ast.Node) {
	switch n := node.(type) {
	default:
		panic("Resolve failed: unknown ast type.")
	case *ast.VariableExpr:
		resolveVariableExpr(n)
	case *ast.AssignExpr:
		resolveAssignExpr(n)
	case *ast.BinaryExpr:
		resolveBinaryExpr(n)
	case *ast.UnaryExpr:
		resolveUnaryExpr(n)
	case *ast.LogicalExpr:
		resolveLogicalExpr(n)
	case *ast.GroupingExpr:
		resolveGroupExpr(n)
	case *ast.CallExpr:
		resolveCallExpr(n)
	case *ast.GetExpr:
		resolveGetExpr(n)
	case *ast.SetExpr:
		resolveSetExpr(n)
	case *ast.ThisExpr:
		resolveThisExpr(n)
	case *ast.Literal:
		// do nothing.
	case *ast.BlockStmt:
		resolveBlockStmt(n)
	case *ast.VarStmt:
		resolveVarStmt(n)
	case *ast.FunctionStmt:
		resolveFunctionStmt(n)
	case *ast.ExprStmt:
		resolveExprStmt(n)
	case *ast.IfStmt:
		resolveIfStmt(n)
	case *ast.WhileStmt:
		resolveWhileStmt(n)
	case *ast.PrintStmt:
		resolvePrintStmt(n)
	case *ast.ReturnStmt:
		resolveReturnStmt(n)
	case *ast.ClassStmt:
		resolveClassStmt(n)
	}
}

func resolveVariableExpr(expr *ast.VariableExpr) {
	if exist, init := scopes.check(expr.Name); exist && !init {
		errors.Error(token.Identifier, "Cannot read local variable in its own initializer.")
		return
	}
	resolveLocal(expr, expr.Name)
}

func resolveLocal(expr ast.Expr, name string) {
	switch n := expr.(type) {
	case *ast.VariableExpr:
		// if variable doesn't exist in scopes, we regard it as a glabol variable.
		for i := len(scopes) - 1; i >= 0; i-- {
			if _, ok := scopes[i][name]; ok {
				n.Distance = len(scopes) - 1 - i
			}
		}
	case *ast.ThisExpr:
		exist := false
		for i := len(scopes) - 1; i >= 0; i-- {
			if _, ok := scopes[i][name]; ok {
				exist = true
				break
			}
		}
		if !exist {
			errors.Error(token.This, "Cannot use 'this' outside of a class.")
		}
	}
}

func resolveAssignExpr(expr *ast.AssignExpr) {
	Resolve(expr.Value)
	resolveLocal(expr.Left, expr.Left.Name)
}

func resolveBinaryExpr(expr *ast.BinaryExpr) {
	Resolve(expr.Left)
	Resolve(expr.Right)
}

func resolveUnaryExpr(expr *ast.UnaryExpr) {
	Resolve(expr.Right)
}

func resolveLogicalExpr(expr *ast.LogicalExpr) {
	Resolve(expr.Left)
	Resolve(expr.Right)
}

func resolveGroupExpr(expr *ast.GroupingExpr) {
	Resolve(expr.Expression)
}

func resolveCallExpr(expr *ast.CallExpr) {
	Resolve(expr.Callee)

	for _, arg := range expr.Arguments {
		Resolve(arg)
	}
}

func resolveGetExpr(expr *ast.GetExpr) {
	Resolve(expr.Object)
}

func resolveSetExpr(expr *ast.SetExpr) {
	Resolve(expr.Object)
	Resolve(expr.Value)
}

func resolveThisExpr(expr *ast.ThisExpr) {
	if curClassType == ClassNone {
		errors.Error(token.This, "Cannot use 'this' outside of a class.")
		return
	}
	resolveLocal(expr, "this")
}

func resolveBlockStmt(block *ast.BlockStmt) {
	scopes.begin()
	resolveBlock(block.Statements)
	scopes.end()
}

func resolveBlock(statements []ast.Stmt) {
	for _, stmt := range statements {
		Resolve(stmt)
	}
}

func resolveVarStmt(stmt *ast.VarStmt) {
	name := stmt.Name.Name
	scopes.declare(name)
	if stmt.Initializer != nil {
		Resolve(stmt.Initializer)
	}
	scopes.define(name)
}

func resolveFunctionStmt(stmt *ast.FunctionStmt) {
	scopes.declare(stmt.Name)
	scopes.define(stmt.Name)
	resolveFunction(stmt, Function)
}

func resolveFunction(function *ast.FunctionStmt, typ functionType) {
	enclosingFunction := curFunctionType
	curFunctionType = typ
	defer func() {
		curFunctionType = enclosingFunction
	}()

	scopes.begin()
	for _, param := range function.Params {
		scopes.declare(param.Name)
		scopes.define(param.Name)
	}
	resolveBlock(function.Body)
	scopes.end()
}

func resolveExprStmt(stmt *ast.ExprStmt) {
	Resolve(stmt.Expression)
}

func resolveIfStmt(stmt *ast.IfStmt) {
	Resolve(stmt.Condition)
	Resolve(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		Resolve(stmt.ElseBranch)
	}
}

func resolveWhileStmt(stmt *ast.WhileStmt) {
	Resolve(stmt.Condition)
	Resolve(stmt.Body)
}

func resolvePrintStmt(stmt *ast.PrintStmt) {
	Resolve(stmt.Expression)
}

func resolveReturnStmt(stmt *ast.ReturnStmt) {
	if curFunctionType == FunctionNone {
		errors.Error(token.Return, "Cannot return from top-level code.")
		return
	}
	if stmt.Value != nil {
		if curFunctionType == Initializer {
			errors.Error(token.Return, "Cannot return a value from an initializer.")
			return
		}
		Resolve(stmt.Value)
	}
}

func resolveClassStmt(stmt *ast.ClassStmt) {
	scopes.declare(stmt.Name)
	scopes.define(stmt.Name)

	enclosingClass := curClassType
	curClassType = Class
	defer func() {
		curClassType = enclosingClass
	}()

	scopes.begin()
	scopes.declare("this")
	scopes.define("this")
	for _, method := range stmt.Methods {
		typ := Method
		if method.IsInitializer {
			typ = Initializer
		}
		resolveFunction(method, typ)
	}
	scopes.end()
}
