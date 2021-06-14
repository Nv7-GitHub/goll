package goll

import (
	"fmt"
	"go/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (p *Program) CompileFuncDecl(decl *ast.FuncDecl) error {
	p.CleanupFunc()

	params := make([]*ir.Param, 0)
	for _, par := range decl.Type.Params.List {
		for _, name := range par.Names {
			kind := p.ConvertTypeString(par.Type.(*ast.Ident).Name)
			param := ir.NewParam(name.Name, kind)

			val, err := p.GetValFromType(par, kind)
			if err != nil {
				return err
			}

			params = append(params, param)
			p.Vars[name.Name] = Variable{
				Value:   val,
				Storage: param,
			}
		}
	}
	var retType types.Type = types.Void
	if decl.Type.Results != nil {
		retType = p.ConvertTypeString(decl.Type.Results.List[0].Type.(*ast.Ident).Name)
	}

	p.Fn = p.M.NewFunc(decl.Name.Name, retType, params...)
	p.Funcs[decl.Name.Name] = p.Fn

	p.Block = p.Fn.NewBlock("entry")

	// Params
	for _, par := range decl.Type.Params.List {
		for _, name := range par.Names {
			kind := p.ConvertTypeString(par.Type.(*ast.Ident).Name)
			val, err := p.GetValFromType(par, kind)
			if err != nil {
				return err
			}
			vr := p.Block.NewAlloca(kind)
			p.Block.NewStore(val.Value(), vr)
			p.Vars[name.Name] = Variable{
				Value:   val,
				Storage: vr,
			}
		}
	}

	var err error
	for _, stmt := range decl.Body.List {
		err = p.CompileStmt(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Program) CompileReturnStmt(stm *ast.ReturnStmt) error {
	if len(stm.Results) < 1 {
		p.Block.NewRet(nil)
		return nil
	}

	v, err := p.CompileExpr(stm.Results[0])
	if err != nil {
		return err
	}

	p.Block.NewRet(v.Value())
	return nil
}

func (p *Program) CompileCallExpr(expr *ast.CallExpr) (Value, error) {
	_, ok := expr.Fun.(*ast.SelectorExpr)
	if ok {
		return nil, fmt.Errorf("%s: selector functions aren't implemented", p.Pos(expr))
	}

	name := expr.Fun.(*ast.Ident).Name
	fn, exists := p.Funcs[name]
	if !exists {
		return nil, fmt.Errorf("%s: unknown function %s", p.Pos(expr), name)
	}

	args := make([]value.Value, len(expr.Args))
	for i, arg := range expr.Args {
		v, err := p.CompileExpr(arg)
		if err != nil {
			return nil, err
		}
		args[i] = v.Value()
	}

	res := p.Block.NewCall(fn, args...)
	v, err := p.GetValFromType(expr, res.Type())
	if err != nil {
		return nil, err
	}
	v.SetValue(res)
	return v, nil
}
