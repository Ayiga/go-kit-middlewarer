package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

func processRequestResponse(gopath string, tb TemplateBase) {
	var buf bytes.Buffer

	tmpl, err := template.ParseFiles(filepath.Join(gopath, "src", "github.com", "ayiga", "go-kit-middlewarer", "tmpl", "transport-request-response.tmpl"))
	if err != nil {
		log.Fatalf("Template Parse Error: %s", err)
	}

	err = tmpl.Execute(&buf, tb)
	if err != nil {
		log.Fatalf("Template execution failed: %s\n", err)
	}

	filename := "request-response_gen.go"

	file := openFile(filepath.Join(".", "transport", "http"), filename)
	defer file.Close()

	fmt.Fprint(file, string(formatBuffer(buf, filename)))
}

func processMakeEndpoint(gopath string, tb TemplateBase) {
	var buf bytes.Buffer

	tmpl, err := template.ParseFiles(filepath.Join(gopath, "src", "github.com", "ayiga", "go-kit-middlewarer", "tmpl", "transport-make-endpoint.tmpl"))
	if err != nil {
		log.Fatalf("Template Parse Error: %s", err)
	}

	err = tmpl.Execute(&buf, tb)
	if err != nil {
		log.Fatalf("Template execution failed: %s\n", err)
	}
	filename := "make-endpoint_gen.go"

	file := openFile(filepath.Join(".", "transport", "http"), filename)
	defer file.Close()

	fmt.Fprint(file, string(formatBuffer(buf, filename)))
}

func processHTTPServer(gopath string, tb TemplateBase) {
	var buf bytes.Buffer

	tmpl, err := template.ParseFiles(filepath.Join(gopath, "src", "github.com", "ayiga", "go-kit-middlewarer", "tmpl", "transport-http-server.tmpl"))
	if err != nil {
		log.Fatalf("Template Parse Error: %s", err)
	}

	err = tmpl.Execute(&buf, tb)
	if err != nil {
		log.Fatalf("Template execution failed: %s\n", err)
	}

	filename := "http-server_gen.go"

	file := openFile(filepath.Join(".", "transport", "http"), filename)
	defer file.Close()

	fmt.Fprint(file, string(formatBuffer(buf, filename)))
}

func processTransportClient(gopath string, tb TemplateBase) {
	var buf bytes.Buffer

	tmpl, err := template.ParseFiles(filepath.Join(gopath, "src", "github.com", "ayiga", "go-kit-middlewarer", "tmpl", "transport-client.tmpl"))
	if err != nil {
		log.Fatalf("Template Parse Error: %s", err)
	}

	err = tmpl.Execute(&buf, tb)
	if err != nil {
		log.Fatalf("Template execution failed: %s\n", err)
	}

	filename := "client_gen.go"

	file := openFile(filepath.Join(".", "transport", "http"), filename)
	defer file.Close()

	fmt.Fprint(file, string(formatBuffer(buf, filename)))
}

func processHTTPInstanceClient(gopath string, tb TemplateBase) {
	var buf bytes.Buffer

	tmpl, err := template.ParseFiles(filepath.Join(gopath, "src", "github.com", "ayiga", "go-kit-middlewarer", "tmpl", "transport-http-client.tmpl"))
	if err != nil {
		log.Fatalf("Template Parse Error: %s", err)
	}

	err = tmpl.Execute(&buf, tb)
	if err != nil {
		log.Fatalf("Template execution failed: %s\n", err)
	}

	filename := "http-client_gen.go"

	file := openFile(filepath.Join(".", "transport", "http"), filename)
	defer file.Close()

	fmt.Fprint(file, string(formatBuffer(buf, filename)))
}

func processHTTPLoadBalancedClient(gopath string, tb TemplateBase) {
	var buf bytes.Buffer

	tmpl, err := template.ParseFiles(filepath.Join(gopath, "src", "github.com", "ayiga", "go-kit-middlewarer", "tmpl", "transport-http-loadbalanced.tmpl"))
	if err != nil {
		log.Fatalf("Template Parse Error: %s", err)
	}

	err = tmpl.Execute(&buf, tb)
	if err != nil {
		log.Fatalf("Template execution failed: %s\n", err)
	}

	filename := "http-client-loadbalanced_gen.go"

	file := openFile(filepath.Join(".", "transport", "http"), filename)
	defer file.Close()

	fmt.Fprint(file, string(formatBuffer(buf, filename)))
}

func processTransport(g *Generator, f *File) {
	gopath := os.Getenv("GOPATH")

	convertedPath := filepath.ToSlash(f.pkg.dir)

	endpointPackage := createImportWithPath(path.Join(convertedPath, "endpoint"))
	basePackage := createImportWithPath(convertedPath)

	for _, interf := range f.interfaces {
		tb := createTemplateBase(basePackage, endpointPackage, interf, f.imports)
		processRequestResponse(gopath, tb)
		processMakeEndpoint(gopath, tb)
		processHTTPServer(gopath, tb)
		processTransportClient(gopath, tb)
		processHTTPInstanceClient(gopath, tb)
		processHTTPLoadBalancedClient(gopath, tb)
	}
}

func init() {
	registerProcess("transport", processTransport)
}
