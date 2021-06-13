package goll

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (p *Program) CompileBinaryExpr(expr *ast.BinaryExpr) (Value, error) {
	lhs, err := p.CompileExpr(expr.X)
	if err != nil {
		return nil, err
	}
	rhs, err := p.CompileExpr(expr.Y)
	if err != nil {
		return nil, err
	}

	l, ok := lhs.(*Int)
	if ok {
		r := rhs.(*Int)
		var lval value.Value
		var rval value.Value
		var kind *types.IntType
		if l.IsShort() && r.IsShort() {
			lval = l.Short(p)
			rval = r.Short(p)
			kind = types.I32
		} else {
			lval = l.Long(p)
			rval = r.Long(p)
			kind = types.I64
		}

		switch expr.Op {
		case token.ADD:
			return NewInt(kind, p.Block.NewAdd(lval, rval)), nil

		case token.SUB:
			return NewInt(kind, p.Block.NewSub(lval, rval)), nil

		case token.MUL:
			return NewInt(kind, p.Block.NewMul(lval, rval)), nil

		case token.QUO:
			return NewInt(kind, p.Block.NewSDiv(lval, rval)), nil

		case token.REM:
			return NewInt(kind, p.Block.NewSRem(lval, rval)), nil

		default:
			return nil, fmt.Errorf("%s: unknown operation type %s", p.Fset.Position(expr.OpPos).String(), expr.Op.String())
		}
	}

	return nil, fmt.Errorf("%s: cannot perform operation on type %T", p.Pos(expr.X), lhs)
}
