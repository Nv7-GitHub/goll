package goll

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func (p *Program) CleanupFunc() {
	if p.Fn != nil && p.Block.Term == nil {
		p.Block.NewRet(nil)
	}
	for _, vr := range p.Vars {
		vr.Value.Cleanup(p)
	}
	p.Vars = make(map[string]Variable)
}

func (p *Program) End() {
	p.CleanupFunc()
}

func (s *String) Cleanup(p *Program) {
	if s.freeable {
		p.Block.NewCall(getFree(p), p.Block.NewLoad(types.I8Ptr, p.Block.NewGetElementPtr(stringTypeMap["string"], s.val, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))))
	}
}
