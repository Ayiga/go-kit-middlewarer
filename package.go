package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
)

type Package struct {
	dir      string
	name     string
	defs     map[*ast.Ident]types.Object
	typesPkg *types.Package

	imports    []Import
	types      []Type
	interfaces []Interface
	structs    []Struct
	files      []*File
	variables  []Variable
	// consts     []*Constants
}

// check type-checks the package.  The package must be OK to proceed.
func (pkg *Package) check(fs *token.FileSet, astFiles []*ast.File) {
	pkg.defs = make(map[*ast.Ident]types.Object)
	config := types.Config{Importer: importer.Default(), FakeImportC: true}
	info := &types.Info{
		Defs: pkg.defs,
	}
	typesPkg, err := config.Check(pkg.dir, fs, astFiles, info)
	if err != nil {
		log.Fatalf("checking package: %s", err)
	}
	pkg.typesPkg = typesPkg
}

func (pkg *Package) Summarize() {
	fmt.Println("Summary")
	fmt.Printf("%s:\n", pkg.name)

	for _, f := range pkg.files {
		fmt.Printf("-  %s:\n", f.fileName)
		fmt.Printf("\t- imports\n")

		for _, i := range f.imports {
			fmt.Printf("\t\t-  %s\n", i.ImportSpec())
		}

		fmt.Printf("\t- interfaces\n")
		for _, i := range f.interfaces {
			fmt.Printf("\t\t- %s\n", i.name)

			for _, t := range i.types {
				fmt.Printf("\t\t\t-  %s\n", t.String())
			}

			for _, m := range i.methods {
				fmt.Printf("\t\t\t-  %s\n", m.InterfaceMethodSpec())
			}
		}

		fmt.Printf("\t- structs\n")
		for _, i := range f.structs {
			fmt.Printf("\t\t-  %s\n", i.name)
		}

		fmt.Printf("\t- types\n")
		for _, t := range f.types {
			fmt.Printf("\t\t-  %s\n", t.String())
		}

		fmt.Printf("\t- variables\n")
		for _, v := range f.variables {
			fmt.Printf("\t\t-  %s\n", v.name)
		}
	}
}
