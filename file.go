package main

import (
	"go/ast"
)

// File holds a single parsed file and associated data.
type File struct {
	pkg      *Package  // Packge to which this file belongs.
	file     *ast.File // Parsed AST.
	path     string
	fileName string
	// These fields are reset for each type being generated.

	imports    []*Import
	types      []Type
	interfaces []Interface
	structs    []Struct
	variables  []Variable
	// consts     []*Constants
}

// genImportsAndTypes, false says we've consumed the entry down, and not to inform us of sub-entries
func (f *File) genImportsAndTypes(node ast.Node) bool {
	switch t := node.(type) {
	case *ast.ImportSpec:
		// filter out context.Context, the reason for this is that we'd like to
		// automatically pass context into the trasnports when they occur.
		if t.Path.Value == "\"context\"" {
			return false
		}

		imp := createImport(t)
		f.imports = append(f.imports, imp)
		f.pkg.imports = append(f.pkg.imports, imp)
		return false
	case *ast.TypeSpec:
		switch s := t.Type.(type) {
		case *ast.InterfaceType:
			// filter out unexported
			if !t.Name.IsExported() {
				return false
			}

			i := createInterface(t.Name.Name, s, []string{}, *f)
			f.interfaces = append(f.interfaces, i)
			f.pkg.interfaces = append(f.pkg.interfaces, i)

			// add it to the list of declared types as well, since it
			// will be identified as a type to import, potentially
			f.pkg.types = append(f.pkg.types, createType(t.Name, f.pkg))
			f.types = append(f.types, createType(t.Name, f.pkg))
			return false
		case *ast.StructType:
			// filter out unexported
			if !t.Name.IsExported() {
				return false
			}
			stru := createStruct(t.Name.Name, s, *f)
			f.pkg.structs = append(f.pkg.structs, stru)
			f.structs = append(f.structs, stru)

			// add it to the list of declared types as well, since it
			// will be identified as a type to import, potentially
			f.pkg.types = append(f.pkg.types, createType(t.Name, f.pkg))
			f.types = append(f.types, createType(t.Name, f.pkg))
			return false
		case *ast.Ident:
			// filter out Unexported
			if !t.Name.IsExported() {
				return false
			}
			// new type, based off of a primitive, likely
			newType := createType(t.Name, f.pkg)
			f.pkg.types = append(f.pkg.types, newType)
			f.types = append(f.types, newType)
			return false
		default:
			log.Printf("Node Type: %#v\n", s)
		}
	case *ast.GenDecl:
		// Gen Decl are grouped declarations, things like
		// var(), const(), import()
		// if we don't handle this, it will be invoked at
		// the top level, and we'll get the underlying nodes
		// anyway...
	case *ast.ValueSpec:
		var typ *Type
		if t.Type != nil {
			t1 := createType(t.Type, f.pkg)
			typ = &t1
		}

		for i := 0; i < len(t.Names); i++ {
			// filter out unexported
			if !t.Names[i].IsExported() {
				continue
			}
			name := t.Names[i].Name
			// value := t.Values[i]
			v := createVariable(name, "", typ)
			f.pkg.variables = append(f.pkg.variables, v)
			f.variables = append(f.variables, v)
		}
		return false
	case *ast.FuncDecl, nil, *ast.Ident, *ast.CommentGroup, *ast.Comment:
		return false
	case *ast.File: // weird
	default:
		log.Printf("Node Type: %#v\n", t)
	}

	return true
}
