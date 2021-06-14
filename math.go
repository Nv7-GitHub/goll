package goll

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/llir/llvm/ir/constant"
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

	// Integer math
	_, ok := lhs.(*Int)
	if ok {
		l := lhs.(*Int)
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

	// Float math
	_, ok = lhs.(*Float)
	if ok {
		l := lhs.(*Float)
		r := rhs.(*Float)
		var lval value.Value
		var rval value.Value
		var kind *types.FloatType
		if l.IsShort() && r.IsShort() {
			lval = l.Short(p)
			rval = r.Short(p)
			kind = types.Float
		} else {
			lval = l.Long(p)
			rval = r.Long(p)
			kind = types.Double
		}

		switch expr.Op {
		case token.ADD:
			return NewFloat(kind, p.Block.NewFAdd(lval, rval)), nil

		case token.SUB:
			return NewFloat(kind, p.Block.NewFSub(lval, rval)), nil

		case token.MUL:
			return NewFloat(kind, p.Block.NewFMul(lval, rval)), nil

		case token.QUO:
			return NewFloat(kind, p.Block.NewFDiv(lval, rval)), nil

		case token.REM:
			return NewFloat(kind, p.Block.NewFRem(lval, rval)), nil

		default:
			return nil, fmt.Errorf("%s: unknown operation type %s", p.Fset.Position(expr.OpPos).String(), expr.Op.String())
		}
	}

	// String concatenation
	_, ok = lhs.(*String)
	if ok {
		l := lhs.(*String)
		r := rhs.(*String)

		stringType := stringTypeMap["string"]
		llen := p.Block.NewLoad(types.I64, p.Block.NewGetElementPtr(stringType, l.Value(), constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)))
		rlen := p.Block.NewLoad(types.I64, p.Block.NewGetElementPtr(stringType, r.Value(), constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)))
		outlen := p.Block.NewAdd(llen, rlen)

		out := p.Block.NewCall(getMalloc(p), p.Block.NewAdd(outlen, constant.NewInt(types.I64, 1))) // Add 1 for null terminator

		// Add null terminator
		end := p.Block.NewGetElementPtr(types.I8, out, outlen)
		p.Block.NewStore(constant.NewInt(types.I8, 0), end)

		// Memcpy first
		p.Block.NewCall(getMemcpy(p), out, l.Data(p), llen)

		// Memcpy second to different chunk
		outptr := p.Block.NewGetElementPtr(types.I8, out, llen)
		p.Block.NewCall(getMemcpy(p), outptr, r.Data(p), rlen)

		l.Cleanup(p)
		r.Cleanup(p)
		res := p.NewString(outlen, out)
		res.SetOwned(false)

		// Create string
		return res, nil
	}

	return nil, fmt.Errorf("%s: cannot perform operation on type %T", p.Pos(expr.X), lhs)
}
