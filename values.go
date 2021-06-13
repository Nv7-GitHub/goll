package goll

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Int struct {
	kind *types.IntType
	val  value.Value
}

func (i *Int) Cleanup(_ *Program)          {}
func (i *Int) SetOwned(_ bool)             {}
func (i *Int) Value() value.Value          { return i.val }
func (i *Int) SetValue(newVal value.Value) { i.val = newVal }
func (i *Int) IsShort() bool               { return i.kind.Equal(types.I32) }
func (i *Int) Copy() Value {
	cp := *i
	return &cp
}
func (i *Int) Short(p *Program) value.Value {
	if i.kind.Equal(types.I32) {
		return i.val
	}
	return p.Block.NewTrunc(i.val, types.I32)
}
func (i *Int) Long(p *Program) value.Value {
	if i.kind.Equal(types.I64) {
		return i.val
	}
	return p.Block.NewSExt(i.val, types.I64)
}

func NewIntConst(kind *types.IntType, val int64) *Int {
	return &Int{
		kind: kind,
		val:  constant.NewInt(kind, val),
	}
}

func NewInt(kind *types.IntType, val value.Value) *Int {
	return &Int{
		kind: kind,
		val:  val,
	}
}
