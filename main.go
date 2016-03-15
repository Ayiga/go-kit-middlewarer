package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	typeNames             = flag.String("type", "", "comma-separated list of type names; must be set")
	middlewaresToGenerate = flag.String("middleware", "logging,instrumenting,transport,zipkin", "comma-seperated list of middlewares to process. Options: [logging,instrumenting,transport,zipkin]")
	summarize             = flag.String("summarize", "", "Prints out the Summary of Found structures intead of generating code")
	binaryName            = ""
)

var bindings map[string]func(*Generator, *File)

func init() {
	bindings = make(map[string]func(*Generator, *File))
}
func registerProcess(argument string, fun func(*Generator, *File)) {
	bindings[argument] = fun
}

var extras map[string]string

// Usage is a replacement usage function for the flags package
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s [flags] -type T [directory]\n")
}

func main() {
	extras = make(map[string]string)
	binaryName = os.Args[0]
	log.SetPrefix(fmt.Sprintf("%s:", binaryName))
	flag.Usage = Usage
	flag.Parse()

	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	types := strings.Split(*typeNames, ",")

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	var (
		dir string
		g   Generator
	)

	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
		g.parsePackageDir(args[0])
		// parsePackageDir(args[0])
	} else {
		dir = filepath.Dir(args[0])
		g.parsePackageFiles(args)
		// parsePackageFiles(args)
	}
	dir = dir
	// Run generate for each type.
	for _, typeName := range types {
		g.generate(typeName)
	}

}

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func prefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}

	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}
