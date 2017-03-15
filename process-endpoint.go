package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

func processEndpoint(g *Generator, f *File) {
	gopath := os.Getenv("GOPATH")
	var buf bytes.Buffer

	tmpl, err := template.New(
		"endpoint",
	).Funcs(
		templateFuncs,
	).ParseFiles(
		filepath.Join(
			gopath,
			"src",
			"github.com",
			"ayiga",
			"go-kit-middlewarer",
			"tmpl",
			"endpoint.tmpl",
		),
	)
	if err != nil {
		log.Fatalf("Template Parse Error: %s", err)
	}

	convertedPath := filepath.ToSlash(f.pkg.dir)

	endpointPackage := createImportWithPath(path.Join(convertedPath, "endpoint"))
	basePackage := createImportWithPath(convertedPath)

	for _, interf := range f.interfaces {
		err := tmpl.ExecuteTemplate(&buf, "endpoint.tmpl", createTemplateBase(basePackage, endpointPackage, interf, f.imports))
		if err != nil {
			log.Fatalf("Template execution failed: %s\n", err)
		}
	}

	filename := "defs_gen.go"

	file := openFile(filepath.Join(".", "endpoint"), filename)
	defer file.Close()

	fmt.Fprint(file, string(formatBuffer(buf, filename)))
}

func init() {
	registerProcess("endpoint", processEndpoint)
}
