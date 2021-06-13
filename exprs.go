package goll

import (
	"fmt"
	"go/ast"
)

func (p *Program) CompileExpr(expr ast.Expr) (Value, error) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return p.CompileBasicLit(e)

	case *ast.BinaryExpr:
		return p.CompileBinaryExpr(e)

	case *ast.Ident:
		return p.CompileIdent(e)

	default:
		return nil, fmt.Errorf("%s: unknown expression type: %T", p.Pos(expr), expr)
	}
}
