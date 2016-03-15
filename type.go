package main

import (
	"fmt"
	"go/ast"
)

type Type struct {
	typ ast.Expr
	pkg *Package
}

func createType(typ ast.Expr, pkg *Package) Type {
	return Type{
		typ: typ,
		pkg: pkg,
	}
}

func (t Type) TypeName() string {
	return ""
}

func (t Type) String() string {
	return resolveFieldTypes(t.typ)
}

func resolveFieldTypes(t ast.Expr) string {
	switch t1 := t.(type) {
	case *ast.StructType:
		return "struct{}"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t1.X, resolveFieldTypes(t1.Sel))
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", resolveFieldTypes(t1.X))
	case *ast.Ident:
		return fmt.Sprintf("%s", t1)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", resolveFieldTypes(t1.Key), resolveFieldTypes(t1.Value))
	case *ast.ArrayType:
		l := ""
		if t1.Len != nil {
			// we have an array, not a slice.. pity...
			l = fmt.Sprintf("%s", t1.Len)
		}
		return fmt.Sprintf("[%s]%s", l, resolveFieldTypes(t1.Elt))
	default:
		return fmt.Sprintf("UKNOWN: +V", t)
	}
}
