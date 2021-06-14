package goll

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/llir/llvm/ir/types"
)

var stringTypeMap = map[string]types.Type{
	"int":    types.I32,
	"int32":  types.I32,
	"int64":  types.I64,
	"string": types.NewStruct(types.I64, types.I8Ptr),
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

	case token.STRING:
		return p.NewStringFromGo(lit.Value), nil

	default:
		return nil, fmt.Errorf("%s: unknown literal type %s", p.Pos(lit), lit.Kind.String())
	}
}

func (p *Program) GetValFromType(n ast.Node, kind types.Type) (Value, error) {
	switch t := kind.(type) {
	case *types.IntType:
		return NewIntConst(t, 0), nil

	default:
		return nil, fmt.Errorf("%s: unknown value for type %s", p.Pos(n), kind.String())
	}
}
