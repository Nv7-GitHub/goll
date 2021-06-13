package goll

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/llir/llvm/ir"
)

type Program struct {
	m    *ir.Module
	fset *token.FileSet

	fn    *ir.Func
	block *ir.Block
}

func CompileDir(dir string) (*ir.Module, error) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

	prog := &Program{
		m:    ir.NewModule(),
		fset: fset,
	}
	for _, p := range parsed {
		for _, file := range p.Files {
			err = prog.AddFile(file)
			if err != nil {
				return nil, err
			}
		}
	}
	prog.End()

	return prog.m, nil
}

func CompileSrc(filename string, src string) (*ir.Module, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return nil, err
	}

	prog := &Program{
		m:    ir.NewModule(),
		fset: fset,
	}
	err = prog.AddFile(file)
	if err != nil {
		return nil, err
	}

	prog.End()

	return prog.m, nil
}

func (p *Program) AddFile(file *ast.File) error {
	for _, decl := range file.Decls {
		err := p.CompileDecl(decl)
		if err != nil {
			return err
		}
	}
	return nil
}
