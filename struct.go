package main

import (
	"go/ast"
)

// Value represents a declared constant.
type Struct struct {
	name  string // the name of the constant.
	types []Type

	pkg  *Package
	file File
}

func createStruct(name string, iface *ast.StructType, file File) Struct {
	stru := Struct{
		name:  name,
		types: nil,
	}

	return stru
}
