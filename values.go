package goll

import (
	"fmt"

	"github.com/llir/irutil"
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
func (i *Int) Data(_ *Program) value.Value { return i.val }
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

type String struct {
	freeable bool
	val      value.Value
}

func (s *String) Cleanup(p *Program) {
	if s.freeable {
		p.Block.NewCall(getFree(p), p.Block.NewExtractValue(s.val, 1))
	}
}
func (s *String) SetOwned(owned bool) { s.freeable = !owned } // True to disable freeable
func (s *String) Value() value.Value  { return s.val }
func (s *String) Data(p *Program) value.Value {
	value := p.Block.NewGetElementPtr(stringTypeMap["string"], s.val, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	return p.Block.NewLoad(types.I8Ptr, value)
}
func (s *String) SetValue(newVal value.Value) { s.val = newVal }
func (s *String) Copy() Value {
	cp := *s
	return &cp
}

func (p *Program) NewStringFromGo(val string) *String {
	stringType := stringTypeMap["string"]
	v := p.Block.NewAlloca(stringType)
	length := p.Block.NewGetElementPtr(stringType, v, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	p.Block.NewStore(constant.NewInt(types.I64, int64(len(val)+1)), length)

	vl := p.M.NewGlobalDef(fmt.Sprintf(".str.%d", p.TmpUsed), irutil.NewCString(val))
	p.TmpUsed++

	ptr := p.Block.NewGetElementPtr(vl.ContentType, vl, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))

	value := p.Block.NewGetElementPtr(stringType, v, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	p.Block.NewStore(ptr, value)

	return &String{
		freeable: false,
		val:      v,
	}
}
