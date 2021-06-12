package goll

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/llir/llvm/ir"
)

func CompileDir(dir string) (*ir.Module, error) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

	m := ir.NewModule()
	for _, p := range parsed {
		for _, file := range p.Files {
			err = AddFile(m, file)
			if err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func CompileSrc(src string) (*ir.Module, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, err
	}

	m := ir.NewModule()
	err = AddFile(m, file)
	return m, err
}

func AddFile(mod *ir.Module, file *ast.File) error {
	return nil
}
