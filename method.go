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
			m.params = append(m.params, createParam(f, names, "input", file))
			names = append(names, m.params[len(m.params)-1].names...)
		}
	}

	if fun.Results != nil {
		for _, f := range fun.Results.List {
			m.results = append(m.results, createParam(f, names, "output", file))
			names = append(names, m.results[len(m.results)-1].names...)
		}
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
