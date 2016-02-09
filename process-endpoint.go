package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
)

func processEndpoint(g *Generator, f *File) {
	gopath := os.Getenv("GOPATH")
	var buf bytes.Buffer

	tmpl, err := template.ParseFiles(gopath + "/src/github.com/ayiga/go-kit-middlewarer/tmpl/endpoint.tmpl")
	if err != nil {
		log.Fatalf("Template Parse Error: %s", err)
	}

	endpointPackage := createImportWithPath(f.pkg.dir + "/endpoint")
	basePackage := createImportWithPath(f.pkg.dir)

	for _, interf := range f.interfaces {
		err := tmpl.Execute(&buf, createTemplateBase(basePackage, endpointPackage, interf, f.imports))
		if err != nil {
			log.Fatalf("Template execution failed: %s\n", err)
		}
	}

	filename := "defs_gen.go"

	file := openFile("./endpoint", filename)
	defer file.Close()

	fmt.Fprint(file, string(formatBuffer(buf, filename)))
}

func init() {
	registerProcess("endpoint", processEndpoint)
}
