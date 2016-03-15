package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

func processLogging(g *Generator, f *File) {
	gopath := os.Getenv("GOPATH")

	var buf bytes.Buffer

	extra, err := template.New("extra").Parse(extras["logging"])
	if err != nil {
		log.Fatalf("Extra Template Parsing Error: %s", err)
	}

	files := []string{
		filepath.Join(gopath, "src", "github.com", "ayiga", "go-kit-middlewarer", "tmpl", "logging.tmpl"),
	}
	tmpl, err := extra.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Template Parse Error: %s", err)
	}

	convertedPath := filepath.ToSlash(f.pkg.dir)

	endpointPackage := createImportWithPath(path.Join(convertedPath, "endpoint"))
	basePackage := createImportWithPath(convertedPath)

	for _, interf := range f.interfaces {
		err := tmpl.ExecuteTemplate(&buf, "logging.tmpl", createTemplateBase(basePackage, endpointPackage, interf, f.imports))
		if err != nil {
			log.Fatalf("Template execution failed: %s\n", err)
		}
	}

	filename := "middleware_gen.go"

	file := openFile(filepath.Join(".", "logging"), filename)
	defer file.Close()

	fmt.Fprint(file, string(formatBuffer(buf, filename)))
}

func init() {
	registerProcess("logging", processLogging)
}
