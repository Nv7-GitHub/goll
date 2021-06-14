package goll

import (
	"fmt"
	"go/ast"

	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	Storage value.Value
	Value   Value
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
			p.Vars[e.Name] = Variable{
				Value:   val,
				Storage: v.Storage,
			}
			return v.Storage, nil
		}

		storage := p.Block.NewAlloca(val.Value().Type())
		p.Vars[e.Name] = Variable{
			Storage: storage,
			Value:   val,
		}
		return storage, nil

	default:
		return nil, fmt.Errorf("%s: cannot store to type %T", p.Pos(expr), expr)
	}
}

func (p *Program) CompileIdent(stm *ast.Ident) (Value, error) {
	v, exists := p.Vars[stm.Name]
	if !exists {
		return nil, fmt.Errorf("%s: no such variable %s", p.Pos(stm), stm.Name)
	}
	newV := v.Value.Copy()
	newV.SetValue(p.Block.NewLoad(newV.Value().Type(), v.Storage))
	return newV, nil
}
