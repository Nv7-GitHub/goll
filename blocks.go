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

func (p *Program) CompileForStmt(stm *ast.ForStmt) error {
	if stm.Init != nil {
		err := p.CompileStmt(stm.Init)
		if err != nil {
			return err
		}
	}

	condblock := p.Fn.NewBlock(fmt.Sprintf("forcond.%d", p.TmpUsed))
	p.TmpUsed++
	p.Block.NewBr(condblock)
	p.Block = condblock

	cond, err := p.CompileExpr(stm.Cond)
	if err != nil {
		return err
	}

	body := p.Fn.NewBlock(fmt.Sprintf("forbody.%d", p.TmpUsed))
	p.TmpUsed++

	end := p.Fn.NewBlock(fmt.Sprintf("forend.%d", p.TmpUsed))
	p.TmpUsed++

	p.Block.NewCondBr(cond.Value(), body, end)

	p.Block = body

	for _, stm := range stm.Body.List {
		err = p.CompileStmt(stm)
		if err != nil {
			return err
		}
	}

	if stm.Post != nil {
		err = p.CompileStmt(stm.Post)
		if err != nil {
			return err
		}
	}

	p.Block.NewBr(condblock)

	p.Block = end
	return nil
}
