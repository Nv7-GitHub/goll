package goll

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type Value interface {
	Cleanup(p *Program)
	SetOwned(bool)
	Value() value.Value
	Copy() Value
	SetValue(value.Value)
}

type Program struct {
	M    *ir.Module
	Fset *token.FileSet

	Fn    *ir.Func
	Block *ir.Block

	Vars map[string]Variable
}

func CompileDir(dir string) (*ir.Module, error) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

	prog := &Program{
		M:    ir.NewModule(),
		Fset: fset,
		Vars: make(map[string]Variable),
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

	return prog.M, nil
}

func CompileSrc(filename string, src string) (*ir.Module, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return nil, err
	}

	prog := &Program{
		M:    ir.NewModule(),
		Fset: fset,
		Vars: make(map[string]Variable),
	}
	err = prog.AddFile(file)
	if err != nil {
		return nil, err
	}

	prog.End()

	return prog.M, nil
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

func (p *Program) Pos(node ast.Node) string {
	return p.Fset.Position(node.Pos()).String()
}
