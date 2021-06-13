package goll

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/llir/llvm/ir"
)

type Program struct {
	m *ir.Module
}

func CompileDir(dir string) (*ir.Module, error) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

	prog := &Program{
		m: ir.NewModule(),
	}
	for _, p := range parsed {
		for _, file := range p.Files {
			err = prog.AddFile(file)
			if err != nil {
				return nil, err
			}
		}
	}

	return prog.m, nil
}

func CompileSrc(src string) (*ir.Module, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, err
	}

	prog := &Program{
		m: ir.NewModule(),
	}
	err = prog.AddFile(file)
	return prog.m, err
}

func (p *Program) AddFile(file *ast.File) error {
	return nil
}
