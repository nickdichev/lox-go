package resolver

import (
	"github.com/ziyoung/lox-go/ast"
	"github.com/ziyoung/lox-go/errors"
	"github.com/ziyoung/lox-go/token"
)

var scopes = NewScopes()

func Resolve(node ast.Node) {
	switch n := node.(type) {
	default:
		panic("Resolve fail: unknown ast type.")
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
	}
	resolveLocal(expr, expr.Name)
}

func resolveLocal(expr ast.Expr, name string) {
	switch n := expr.(type) {
	case *ast.VariableExpr:
		for i := len(scopes) - 1; i >= 0; i-- {
			if _, ok := scopes[i][name]; ok {
				n.Distance = len(scopes) - 1 - i
			}
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
	resolveFunction(stmt)
}

func resolveFunction(function *ast.FunctionStmt) {
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
	if stmt.Value != nil {
		Resolve(stmt.Value)
	}
}

func resolveClassStmt(stmt *ast.ClassStmt) {
	scopes.declare(stmt.Name)
	scopes.define(stmt.Name)

	scopes.begin()
	for _, method := range stmt.Methods {
		resolveFunction(method)
	}
	scopes.end()
}
