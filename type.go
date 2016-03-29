package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type Type struct {
	typ            ast.Expr
	requiredImport *Import // required import to use, can be nil?
	pkg            *Package

	isBasePackage bool
}

func createType(typ ast.Expr, pkg *Package) Type {
	t := Type{
		typ: typ,
		pkg: pkg,
	}

	if !strings.Contains(t.String(), ".") &&
		strings.ToLower(t.String()) != t.String() { // exported.
		t.isBasePackage = true
	}

	return t
}

func (t Type) TypeName() string {
	return ""
}

func (t Type) Equal(o Type) bool {
	return t.String() == o.String()
}

func (t Type) String() string {
	if t.isBasePackage {
		return resolveFieldTypes(t.typ, t.pkg.name)
	}
	return resolveFieldTypes(t.typ, "")
}

func resolveFieldTypes(t ast.Expr, pkgName string) string {
	switch t1 := t.(type) {
	case *ast.StructType:
		return "struct{}"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.SelectorExpr:
		// we have to override this crap...
		return fmt.Sprintf("%s.%s", t1.X, resolveFieldTypes(t1.Sel, ""))
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", resolveFieldTypes(t1.X, pkgName))
	case *ast.Ident:
		if pkgName != "" {
			return fmt.Sprintf("%s.%s", pkgName, t1)
		}
		return fmt.Sprintf("%s", t1)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", resolveFieldTypes(t1.Key, pkgName), resolveFieldTypes(t1.Value, pkgName))
	case *ast.ArrayType:
		l := ""
		if t1.Len != nil {
			// we have an array, not a slice.. pity...
			l = fmt.Sprintf("%s", t1.Len)
		}
		return fmt.Sprintf("[%s]%s", l, resolveFieldTypes(t1.Elt, pkgName))
	default:
		return fmt.Sprintf("UKNOWN: +V", t)
	}
}
