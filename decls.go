package goll

import (
	"fmt"
	"go/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (p *Program) CompileDecl(decl ast.Decl) error {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return p.CompileFuncDecl(d)

	default:
		return fmt.Errorf("%s: unknown declaration type: %T", p.fset.Position(decl.Pos()).String(), decl)
	}
}

func (p *Program) CompileFuncDecl(decl *ast.FuncDecl) error {
	p.CleanupFunc()

	params := make([]*ir.Param, 0)
	for _, par := range decl.Type.Params.List {
		for _, name := range par.Names {
			params = append(params, ir.NewParam(name.Name, p.ConvertTypeString(par.Type.(*ast.Ident).Name)))
		}
	}
	var retType types.Type = types.Void
	if decl.Type.Results != nil {
		retType = p.ConvertTypeString(decl.Type.Results.List[0].Type.(*ast.Ident).Name)
	}

	p.fn = p.m.NewFunc(decl.Name.Name, retType, params...)
	p.block = p.fn.NewBlock("entry")

	return nil
}
