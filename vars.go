package goll

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	Storage value.Value
	Value   Value
}

var assignMathMap = map[token.Token]token.Token{
	token.ADD_ASSIGN: token.ADD,
	token.SUB_ASSIGN: token.SUB,
	token.MUL_ASSIGN: token.MUL,
	token.QUO_ASSIGN: token.QUO,
	token.REM_ASSIGN: token.REM,
}

func (p *Program) CompileAssignStmt(stm *ast.AssignStmt) error {
	rhs, err := p.CompileExpr(stm.Rhs[0])
	if err != nil {
		return err
	}

	st, err := p.GetStorable(stm.Lhs[0], rhs, stm.Tok == token.DEFINE)
	if err != nil {
		return err
	}

	op, exists := assignMathMap[stm.Tok]
	if exists {
		lhs, err := p.CompileExpr(stm.Lhs[0])
		if err != nil {
			return err
		}
		rhs, err = p.BinaryExpr(lhs, rhs, op, stm.TokPos, stm.Lhs[0].Pos())
		if err != nil {
			return err
		}
	}

	p.Block.NewStore(rhs.Value(), st)
	return nil
}

func (p *Program) GetStorable(expr ast.Expr, val Value, redefine bool) (value.Value, error) {
	switch e := expr.(type) {
	case *ast.Ident:
		v, exists := p.Vars[e.Name]
		if exists && !redefine {
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
		if stm.Name == "true" {
			return NewBoolConst(true), nil
		}
		if stm.Name == "false" {
			return NewBoolConst(false), nil
		}
		return nil, fmt.Errorf("%s: no such variable %s", p.Pos(stm), stm.Name)
	}
	newV := v.Value.Copy()
	newV.SetValue(p.Block.NewLoad(newV.Value().Type(), v.Storage))
	newV.SetOwned(true)
	return newV, nil
}
