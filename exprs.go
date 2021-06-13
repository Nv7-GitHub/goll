package goll

import (
	"fmt"
	"go/ast"
)

func (p *Program) CompileExpr(expr ast.Expr) (Value, error) {
	switch expr.(type) {
	default:
		return nil, fmt.Errorf("%s: unknown expression type: %T", p.Fset.Position(expr.Pos()).String(), expr)
	}
}
