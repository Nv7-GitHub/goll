package goll

import (
	"fmt"
	"go/ast"
)

func (p *Program) AddDecl(decl ast.Decl) error {
	switch decl.(type) {
	default:
		return fmt.Errorf("%s: unknown declaration type: %T", p.fset.Position(decl.Pos()).String(), decl)
	}
}
