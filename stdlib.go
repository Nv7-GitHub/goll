package goll

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func getFree(p *Program) *ir.Func {
	_, exists := p.CFuncs["free"]
	if !exists {
		p.CFuncs["free"] = p.M.NewFunc("free", types.Void, ir.NewParam("src", types.I8Ptr))
	}
	return p.CFuncs["free"]
}
