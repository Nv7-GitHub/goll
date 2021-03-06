package goll

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func getFree(p *Program) *ir.Func {
	name := "free"
	_, exists := p.CFuncs[name]
	if !exists {
		p.CFuncs[name] = p.M.NewFunc(name, types.Void, ir.NewParam("src", types.I8Ptr))
	}
	return p.CFuncs[name]
}

func getMalloc(p *Program) *ir.Func {
	name := "malloc"
	_, exists := p.CFuncs[name]
	if !exists {
		p.CFuncs[name] = p.M.NewFunc(name, types.I8Ptr, ir.NewParam("len", types.I64))
	}
	return p.CFuncs[name]
}

func getMemcpy(p *Program) *ir.Func {
	name := "memcpy"
	_, exists := p.CFuncs[name]
	if !exists {
		p.CFuncs[name] = p.M.NewFunc(name, types.I8Ptr, ir.NewParam("dst", types.I8Ptr), ir.NewParam("src", types.I8Ptr), ir.NewParam("cnt", types.I64))
	}
	return p.CFuncs[name]
}

func getPrintf(p *Program) *ir.Func {
	name := "printf"
	_, exists := p.CFuncs[name]
	if !exists {
		printf := p.M.NewFunc(name, types.I32, ir.NewParam("format", types.I8Ptr))
		printf.Sig.Variadic = true
		p.CFuncs[name] = printf
	}
	return p.CFuncs[name]
}

func getStrcmp(p *Program) *ir.Func {
	name := "strcmp"
	_, exists := p.CFuncs[name]
	if !exists {
		p.CFuncs[name] = p.M.NewFunc(name, types.I32, ir.NewParam("a", types.I8Ptr), ir.NewParam("b", types.I8Ptr))
	}
	return p.CFuncs[name]
}
