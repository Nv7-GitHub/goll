package goll

import (
	"fmt"
	"go/ast"

	"github.com/llir/llvm/ir"
)

func (p *Program) CompileIfStmt(stm *ast.IfStmt) error {
	cond, err := p.CompileExpr(stm.Cond)
	if err != nil {
		return err
	}

	start := p.Block
	body := p.Fn.NewBlock(fmt.Sprintf("ifbody.%d", p.TmpUsed))
	p.TmpUsed++
	p.Block = body

	for _, stm := range stm.Body.List {
		err = p.CompileStmt(stm)
		if err != nil {
			return err
		}
	}

	var el *ir.Block
	if stm.Else != nil {
		el = p.Fn.NewBlock(fmt.Sprintf("ifelse.%d", p.TmpUsed))
		p.TmpUsed++
		p.Block = el

		for _, stm := range stm.Else.(*ast.BlockStmt).List {
			err = p.CompileStmt(stm)
			if err != nil {
				return err
			}
		}
	}

	end := p.Fn.NewBlock(fmt.Sprintf("ifend.%d", p.TmpUsed))
	p.TmpUsed++
	if stm.Else == nil {
		start.NewCondBr(cond.Value(), body, end)
	} else {
		start.NewCondBr(cond.Value(), body, el)
		el.NewBr(end)
	}

	body.NewBr(end)

	p.Block = end
	return nil
}
