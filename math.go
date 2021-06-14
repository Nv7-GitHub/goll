package goll

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (p *Program) MathExpr(lhs, rhs Value, op token.Token, opPos token.Pos, lPos token.Pos) (Value, error) {
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

		switch op {
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

		case token.EQL:
			return NewBool(p.Block.NewICmp(enum.IPredEQ, lval, rval)), nil

		case token.LSS:
			return NewBool(p.Block.NewICmp(enum.IPredSLT, lval, rval)), nil

		case token.GTR:
			return NewBool(p.Block.NewICmp(enum.IPredSGT, lval, rval)), nil

		case token.NEQ:
			return NewBool(p.Block.NewICmp(enum.IPredNE, lval, rval)), nil

		case token.LEQ:
			return NewBool(p.Block.NewICmp(enum.IPredSLE, lval, rval)), nil

		case token.GEQ:
			return NewBool(p.Block.NewICmp(enum.IPredSGE, lval, rval)), nil

		default:
			return nil, fmt.Errorf("%s: unknown operation type %s", p.Fset.Position(opPos).String(), op.String())
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

		switch op {
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

		case token.EQL:
			return NewBool(p.Block.NewFCmp(enum.FPredOEQ, lval, rval)), nil

		case token.LSS:
			return NewBool(p.Block.NewFCmp(enum.FPredOLT, lval, rval)), nil

		case token.GTR:
			return NewBool(p.Block.NewFCmp(enum.FPredOGT, lval, rval)), nil

		case token.NEQ:
			return NewBool(p.Block.NewFCmp(enum.FPredONE, lval, rval)), nil

		case token.LEQ:
			return NewBool(p.Block.NewFCmp(enum.FPredOLE, lval, rval)), nil

		case token.GEQ:
			return NewBool(p.Block.NewFCmp(enum.FPredOGE, lval, rval)), nil

		default:
			return nil, fmt.Errorf("%s: unknown operation type %s", p.Fset.Position(opPos).String(), op.String())
		}
	}

	// String concatenation
	_, ok = lhs.(*String)
	if ok {
		l := lhs.(*String)
		r := rhs.(*String)

		var lval value.Value
		var rval value.Value
		if op == token.EQL || op == token.LSS || op == token.GTR || op == token.NEQ || op == token.LEQ || op == token.GEQ {
			lval = l.Data(p)
			rval = r.Data(p)
		}
		switch op {
		case token.EQL:
			cmped := p.Block.NewCall(getStrcmp(p), lval, rval)
			return NewBool(p.Block.NewICmp(enum.IPredEQ, cmped, constant.NewInt(types.I32, 0))), nil

		case token.LSS:
			cmped := p.Block.NewCall(getStrcmp(p), lval, rval)
			return NewBool(p.Block.NewICmp(enum.IPredSLT, cmped, constant.NewInt(types.I32, 0))), nil

		case token.GTR:
			cmped := p.Block.NewCall(getStrcmp(p), lval, rval)
			return NewBool(p.Block.NewICmp(enum.IPredSGT, cmped, constant.NewInt(types.I32, 0))), nil

		case token.NEQ:
			cmped := p.Block.NewCall(getStrcmp(p), lval, rval)
			return NewBool(p.Block.NewICmp(enum.IPredNE, cmped, constant.NewInt(types.I32, 0))), nil

		case token.LEQ:
			cmped := p.Block.NewCall(getStrcmp(p), lval, rval)
			return NewBool(p.Block.NewICmp(enum.IPredSLE, cmped, constant.NewInt(types.I32, 0))), nil

		case token.GEQ:
			cmped := p.Block.NewCall(getStrcmp(p), lval, rval)
			return NewBool(p.Block.NewICmp(enum.IPredSGE, cmped, constant.NewInt(types.I32, 0))), nil

		default:
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
	}

	return nil, fmt.Errorf("%s: cannot perform operation on type %T", p.Fset.Position(lPos).String(), lhs)
}

func (p *Program) CompileBinaryExpr(expr *ast.BinaryExpr) (Value, error) {
	lhs, err := p.CompileExpr(expr.X)
	if err != nil {
		return nil, err
	}
	rhs, err := p.CompileExpr(expr.Y)
	if err != nil {
		return nil, err
	}

	return p.MathExpr(lhs, rhs, expr.Op, expr.OpPos, expr.X.Pos())
}
