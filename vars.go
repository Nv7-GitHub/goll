package goll

import (
	"fmt"
	"go/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	Value   Value
	Storage *ir.InstAlloca
}

func (p *Program) CompileAssignStmt(stm *ast.AssignStmt) error {
	rhs, err := p.CompileExpr(stm.Rhs[0])
	if err != nil {
		return err
	}
	st, err := p.GetStorable(stm.Lhs[0], rhs)
	if err != nil {
		return err
	}

	p.Block.NewStore(rhs.Value(), st)
	return nil
}

func (p *Program) GetStorable(expr ast.Expr, val Value) (value.Value, error) {
	switch e := expr.(type) {
	case *ast.Ident:
		v, exists := p.Vars[e.Name]
		if exists {
			v.Value.Cleanup(p)
			return v.Storage, nil
		}
		return p.Block.NewAlloca(val.Value().Type()), nil

	default:
		return nil, fmt.Errorf("%s: cannot store to type %T", p.Fset.Position(expr.Pos()).String(), expr)
	}
}
