package main

import (
	"go/ast"
	"strings"
)

// Value represents a declared constant.
type Interface struct {
	name    string // the name of the constant.
	methods []Method
	types   []Type

	pkg  *Package
	file File
}

func createInterface(name string, iface *ast.InterfaceType, reservedNames []string, file File) Interface {
	names := append([]string{}, reservedNames...)
	names = append(names, name)

	interf := Interface{
		name:    name,
		methods: make([]Method, 0, iface.Methods.NumFields()),
		types:   nil,
	}
	for _, f := range iface.Methods.List {
		if len(f.Names) > 0 {
			// This is a method
			interf.methods = append(interf.methods, createMethod(f, names, file))
		} else {
			// this is an interface.
			// var suggestedName = ""

			n := resolveFieldTypes(f.Type)
			potentialNamePieces := strings.Split(n, ".")
			if len(potentialNamePieces) > 0 {
				// suggestedName = strings.ToLower(potentialNamePieces[len(potentialNamePieces)-1])
			}

			// p := createParam(f, reservedNames, suggestedName, file)

			// if len(p.names) > 0 {
			// 	reservedNames = append(reservedNames, p.names[0])
			// 	// add it to the reserved name
			// }

			interf.types = append(interf.types, createType(f.Type, file.pkg))
		}
	}
	return interf
}
