package goll

import (
	"fmt"
	"go/ast"
)

func (p *Program) CompileStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		return p.CompileAssignStmt(s)

	case *ast.ReturnStmt:
		return p.CompileReturnStmt(s)

	case *ast.ExprStmt:
		_, err := p.CompileExpr(s.X)
		return err

	case *ast.IfStmt:
		return p.CompileIfStmt(s)

	default:
		return fmt.Errorf("%s: unknown statement type: %T", p.Pos(stmt), stmt)
	}
}
