package goll

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/llir/llvm/ir/types"
)

var stringTypeMap = map[string]types.Type{
	"int":   types.I32,
	"int32": types.I32,
	"int64": types.I64,
}

func (p *Program) ConvertTypeString(t string) types.Type {
	return stringTypeMap[t]
}

func (p *Program) CompileBasicLit(lit *ast.BasicLit) (Value, error) {
	switch lit.Kind {
	case token.INT:
		v, err := strconv.ParseInt(lit.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		return NewIntConst(stringTypeMap["int"].(*types.IntType), v), nil

	default:
		return nil, fmt.Errorf("%s: unknown literal type %s", p.Pos(lit), lit.Kind.String())
	}
}
