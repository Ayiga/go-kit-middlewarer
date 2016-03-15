package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type Param struct {
	names []string
	typ   Type
}

func createParam(field *ast.Field, reservedNames []string, suggestion string, file File) Param {
	p := Param{
		names: make([]string, 0, len(field.Names)),
		typ:   createType(field.Type, file.pkg),
	}

	for _, n := range field.Names {
		p.names = append(p.names, n.Name)
	}

	// no name specified, let's create one...
	if len(p.names) == 0 {
		n := determineLocalName(suggestion, reservedNames)
		p.names = []string{n}
	}

	return p
}

func (p Param) ParamSpec() string {
	if len(p.names) > 0 {
		return fmt.Sprintf("%s %s", strings.Join(p.names, ", "), p.typ)
	}
	return p.typ.String()
}
