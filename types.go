package goll

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/llir/llvm/ir/types"
)

var stringTypeMap = map[string]types.Type{
	"int":     types.I32,
	"int32":   types.I32,
	"int64":   types.I64,
	"float":   types.Float,
	"float32": types.Float,
	"float64": types.Double,
	"string":  types.NewStruct(types.I64, types.I8Ptr),
}

func (p *Program) ConvertTypeString(t string) types.Type {
	return stringTypeMap[t]
}

func EvaluateString(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s[1:len(s)-1], "\\n", "\n"), "\\", "")
}

func (p *Program) CompileBasicLit(lit *ast.BasicLit) (Value, error) {
	switch lit.Kind {
	case token.INT:
		v, err := strconv.ParseInt(lit.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		return NewIntConst(stringTypeMap["int"].(*types.IntType), v), nil

	case token.FLOAT:
		v, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			return nil, err
		}
		return NewFloatConst(stringTypeMap["float"].(*types.FloatType), v), nil

	case token.STRING:
		return p.NewStringFromGo(EvaluateString(lit.Value)), nil

	default:
		return nil, fmt.Errorf("%s: unknown literal type %s", p.Pos(lit), lit.Kind.String())
	}
}

func (p *Program) GetValFromType(n ast.Node, kind types.Type) (Value, error) {
	switch t := kind.(type) {
	case *types.IntType:
		return NewIntConst(t, 0), nil
	}

	return nil, fmt.Errorf("%s: unknown value for type %s", p.Pos(n), kind.String())
}
