package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type Method struct {
	name    string
	params  []Param
	results []Param

	hasContextParam   bool
	contextParamName  string
	hasErrResult      bool
	errorResultName   string
	moreThanOneResult bool

	pkg        *Package
	file       File
	imports    []Import
	types      []Type
	interfaces []Interface
}

func createMethod(field *ast.Field, reservedNames []string, file File) Method {
	name := field.Names[0].Name
	names := append([]string{}, reservedNames...)
	names = append(names, name)
	fun, ok := field.Type.(*ast.FuncType)
	if !ok {
		return Method{
			name: name,
		}
	}

	m := Method{
		name:    name,
		params:  make([]Param, 0, fun.Params.NumFields()),
		results: make([]Param, 0, fun.Results.NumFields()),
	}

	if fun.Params != nil {
		for _, f := range fun.Params.List {
			param := createParam(f, names, "input", file)
			paramNames := param.names

			if param.typ.String() == "context.Context" {
				m.hasContextParam = true
				m.contextParamName = paramNames[0]
			}

			for _, imp := range file.pkg.imports {
				if strings.HasPrefix(param.typ.String(), fmt.Sprintf("%s.", imp.name)) {
					imp.isParam = true
				}
			}

			m.params = append(m.params, param)
			names = append(names, paramNames...)
		}
	}

	if fun.Results != nil {
		numResult := 0
		for _, f := range fun.Results.List {
			param := createParam(f, names, "output", file)
			paramNames := param.names

			if param.typ.String() == "error" {
				m.hasErrResult = true
				m.errorResultName = "err"
				if len(paramNames) > 0 {
					m.errorResultName = paramNames[0]
				}
			}

			for _, imp := range file.pkg.imports {
				if strings.HasPrefix(param.typ.String(), fmt.Sprintf("%s.", imp.name)) {
					imp.isParam = true
				}
			}

			if len(paramNames) > 0 {
				numResult += len(paramNames)
			} else {
				numResult++
			}

			m.results = append(m.results, param)
			names = append(names, paramNames...)
		}

		m.moreThanOneResult = numResult > 1
	}

	return m
}

func (m Method) usedNames() []string {
	var result []string
	result = append(result, m.name)
	for _, p := range m.params {
		result = append(result, p.names...)
	}

	for _, r := range m.results {
		result = append(result, r.names...)
	}
	return result
}

func (m Method) InterfaceMethodSpec() string {
	params := make([]string, 0, len(m.params))
	results := make([]string, 0, len(m.results))
	resultNameCount := 0

	for _, p := range m.params {
		params = append(params, p.ParamSpec())
	}

	for _, r := range m.results {
		results = append(results, r.ParamSpec())
		resultNameCount += len(r.names)
	}

	if len(m.results) > 1 || resultNameCount > 0 {
		return fmt.Sprintf("%s(%s) (%s)", m.name, strings.Join(params, ", "), strings.Join(results, ", "))
	}
	return fmt.Sprintf("%s(%s) %s", m.name, strings.Join(params, ", "), strings.Join(results, ", "))
}

func (m Method) methodArguments() string {
	var result []string
	for _, p := range m.params {
		result = append(result, p.ParamSpec())
	}

	return strings.Join(result, ", ")
}

func (m Method) methodArgumentNames() string {
	var result []string
	for _, p := range m.params {
		result = append(result, p.names...)
	}

	return strings.Join(result, ", ")
}

func (m Method) methodResults() string {
	var result []string
	for _, p := range m.results {
		result = append(result, p.ParamSpec())
	}

	return strings.Join(result, ", ")
}

func (m Method) methodResultNames() string {
	var result []string
	for _, p := range m.results {
		result = append(result, p.names...)
	}

	return strings.Join(result, ", ")
}
