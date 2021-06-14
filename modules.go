package goll

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/llir/llvm/ir/value"
)

type Module map[string]func(p *Program, args ...ast.Expr) (Value, error)

var modules = map[string]Module{
	"external": map[string]func(p *Program, args ...ast.Expr) (Value, error){
		"Printf": func(p *Program, rawArgs ...ast.Expr) (Value, error) {
			fn := getPrintf(p)
			args := make([]value.Value, len(rawArgs))
			for i, arg := range rawArgs {
				argVal, err := p.CompileExpr(arg)
				if err != nil {
					return nil, err
				}
				args[i] = argVal.Data(p)
			}
			p.Block.NewCall(fn, args...)
			return nil, nil
		},
	},
}

func (p *Program) CompileGenDecl(decl *ast.GenDecl) error {
	var err error
	for _, val := range decl.Specs {
		err = p.CompileSpec(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Program) CompileSpec(spec ast.Spec) error {
	switch s := spec.(type) {
	case *ast.ImportSpec:
		name := strings.ReplaceAll(s.Path.Value, "\"", "")
		mod, exists := modules[name]
		if !exists {
			return fmt.Errorf("%s: no such module %s", p.Pos(spec), name)
		}
		p.Modules[name] = mod
		return nil

	default:
		return fmt.Errorf("%s: unknown spec type %T", p.Pos(spec), spec)
	}
}
