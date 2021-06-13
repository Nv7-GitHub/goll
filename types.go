package goll

import "github.com/llir/llvm/ir/types"

var stringTypeMap = map[string]types.Type{
	"int": types.I32,
}

func (p *Program) ConvertTypeString(t string) types.Type {
	return stringTypeMap[t]
}
