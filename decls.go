package goll

import (
	"fmt"
	"go/ast"
)

func (p *Program) CompileDecl(decl ast.Decl) error {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return p.CompileFuncDecl(d)

	default:
		return fmt.Errorf("%s: unknown declaration type: %T", p.Pos(decl), decl)
	}
}
