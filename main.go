package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	typeNames   = flag.String("type", "", "comma-separated list of type names; must be set")
	middlewares = flag.String("middleware", "logging,instrumenting,transport,zipkin", "comma-seperated list of middlewares to process. Options: [logging,instrumenting,transport,zipkin]")
	binaryName  = ""
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
	log.SetFlags(0)
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

// Generator holds the state of the analysis.  Primarily used to buffer the
// output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated ouptu.
	pkg *Package
}

// Printf writes the given output to the internalized buffer.
func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

type Import struct {
	name string
	path string
	last string
}

func createImportWithPath(p string) Import {
	last := path.Base(p)
	return Import{
		name: last,
		path: p,
		last: last,
	}
}

func createImport(imp *ast.ImportSpec) Import {
	var name string
	pth := strings.TrimPrefix(strings.TrimSuffix(imp.Path.Value, "\""), "\"")
	last := path.Base(pth)
	if n := imp.Name; n == nil {
		name = last
	} else {
		name = n.String()
	}
	return Import{
		name: name,
		path: pth,
		last: last,
	}
}

func (i Import) ImportSpec() string {
	if i.name == i.last {
		return fmt.Sprintf("\"%s\"", i.path)
	}

	return fmt.Sprintf("%s \"%s\"", i.name, i.path)
}

type Type struct {
	name string
}

func createType(typ *ast.TypeSpec) string {
	return typ.Name.String()
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Packge to which this file belongs.
	file *ast.File // Parsed AST.
	path string
	// These fields are reset for each type being generated.
	typeName   string      // Name of the constant type.
	interfaces []Interface // Accumulator for constant values of that type.
	imports    []Import
	types      []string

	interfacesToResolve []ifaceArgs
}

type Package struct {
	dir      string
	name     string
	defs     map[*ast.Ident]types.Object
	files    []*File
	typesPkg *types.Package
}

// Value represents a declared constant.
type Interface struct {
	name    string // the name of the constant.
	methods []Method
	types   []Param
}

func createInterface(name string, iface *ast.InterfaceType, reservedNames []string, file File) Interface {
	names := append([]string{}, reservedNames...)
	names = append(names, name)

	interf := Interface{
		name:    name,
		methods: make([]Method, 0, iface.Methods.NumFields()),
		types:   nil,
	}
	for _, f := range iface.Methods.List {
		if len(f.Names) > 0 {
			// This is a method
			interf.methods = append(interf.methods, createMethod(f, names, file))
		} else {
			// this is an interface.
			var suggestedName = ""

			n := resolveFieldTypes(f.Type)
			potentialNamePieces := strings.Split(n, ".")
			if len(potentialNamePieces) > 0 {
				suggestedName = strings.ToLower(potentialNamePieces[len(potentialNamePieces)-1])
			}

			p := createParam(f, reservedNames, suggestedName, file)

			if len(p.names) > 0 {
				reservedNames = append(reservedNames, p.names[0])
				// add it to the reserved name
			}

			interf.types = append(interf.types, p)
		}
	}
	return interf
}

type Method struct {
	name    string
	params  []Param
	results []Param
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

type Param struct {
	names []string
	typ   string
}

func createParam(field *ast.Field, reservedNames []string, suggestion string, file File) Param {
	typ := resolveFieldTypes(field.Type)
	if sliceContains(file.types, typ) {
		typ = file.pkg.name + "." + typ
	}
	p := Param{
		names: make([]string, 0, len(field.Names)),
		typ:   typ,
	}

	for _, n := range field.Names {
		p.names = append(p.names, n.Name)
	}

	// no name specified, let's create one...
	if len(p.names) == 0 {
		n := determineLocalName(suggestion, reservedNames)
		p.names = []string{n}
	}

	return p
}

func (p Param) ParamSpec() string {
	if len(p.names) > 0 {
		return fmt.Sprintf("%s %s", strings.Join(p.names, ", "), p.typ)
	}
	return p.typ
}

func (g *Generator) parsePackageDir(directory string) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		log.Fatalf("cannot process directory %s: %s", directory, err)
	}

	d, e := os.Getwd()
	gopath := os.Getenv("GOPATH")
	if e != nil {
		log.Fatalf("Error Grabbing WD: %s\n", e)
	}

	d = strings.TrimPrefix(d, fmt.Sprintf("%s/src/", gopath))

	var names []string
	names = append(names, pkg.GoFiles...)
	names = append(names, pkg.CgoFiles...)

	names = append(names, pkg.SFiles...)
	names = prefixDirectory(directory, names)
	g.parsePackage(d, names, nil)
}

func (g *Generator) parsePackageFiles(names []string) {
	g.parsePackage(".", names, nil)
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

// parsePackage analyzes the signle package constructed from the named files.
// If text is non-nil, it is a string to be used instead of the content of the file,
// to be used for testing.  parsePackage exists if there is an error.
func (g *Generator) parsePackage(directory string, names []string, text interface{}) {
	var files []*File
	var astFiles []*ast.File
	g.pkg = new(Package)
	fs := token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}

		parsedFile, err := parser.ParseFile(fs, name, text, parser.ParseComments)

		for _, v := range parsedFile.Comments {
			str := v.Text()
			if strings.HasPrefix(str, log.Prefix()) {
				lines := strings.Split(str, "\n")
				if len(lines) <= 0 {
					continue
				}
				var firstLine = lines[0]

				typ := strings.TrimPrefix(firstLine, log.Prefix())
				if len(lines) > 1 {
					extras[typ] = strings.Join(lines[1:], "\n")
				}
			}
		}

		if err != nil {
			log.Fatalf("parsing package: %s: %s", name, err)
		}
		astFiles = append(astFiles, parsedFile)
		files = append(files, &File{
			file: parsedFile,
			pkg:  g.pkg,
			// path:
		})
	}

	if len(astFiles) == 0 {
		log.Fatalf("%s: no buildable Go files", directory)
	}
	g.pkg.name = astFiles[0].Name.Name
	g.pkg.files = files
	g.pkg.dir = directory
	g.pkg.check(fs, astFiles)
}

// check type-checks the package.  The package must be OK to proceed.
func (pkg *Package) check(fs *token.FileSet, astFiles []*ast.File) {
	pkg.defs = make(map[*ast.Ident]types.Object)
	config := types.Config{Importer: importer.Default(), FakeImportC: true}
	info := &types.Info{
		Defs: pkg.defs,
	}
	typesPkg, err := config.Check(pkg.dir, fs, astFiles, info)
	if err != nil {
		log.Fatalf("checking package: %s", err)
	}
	pkg.typesPkg = typesPkg
}

// generate does 'things'
func (g *Generator) generate(typeName string) {
	found := false

	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.interfaces = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genImportsAndTypes)
			file.resolveInterfaces()
			if len(file.interfaces) > 0 {
				found = true
				list := strings.Split(*middlewares, ",")
				list = append(list, "endpoint")
				for _, l := range list {
					if bindings[l] != nil {
						bindings[l](g, file)
					}
				}
			}
		}
	}

	if !found {
		log.Fatalf("no interfaces defined for type %s", typeName)
	}
}

// format returns gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the pacakge to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

type ifaceArgs struct {
	iface *ast.InterfaceType
	name  string
	names []string
}

// genImportsAndTypes
func (f *File) genImportsAndTypes(node ast.Node) bool {
	imports, ok := node.(*ast.ImportSpec)
	if ok {
		imp := createImport(imports)
		f.imports = append(f.imports, imp)
	}

	decl, ok := node.(*ast.TypeSpec)
	if !ok {
		// this is not an interface...
		return true
	}

	f.types = append(f.types, createType(decl))

	if decl.Name.Name != f.typeName {
		// we don't have the correct interface
		return true
	}

	iface, ok := decl.Type.(*ast.InterfaceType)
	if !ok {
		log.Fatalf("The specified type: %s is not an interface!", f.typeName)
	}

	names := make([]string, 0, len(f.imports))
	for _, i := range f.imports {
		names = append(names, i.name)
	}

	f.interfacesToResolve = append(f.interfacesToResolve, ifaceArgs{
		iface: iface,
		name:  decl.Name.Name,
		names: names,
	})
	return false
}

func (f *File) resolveInterfaces() {
	for _, args := range f.interfacesToResolve {
		f.interfaces = append(f.interfaces, createInterface(args.name, args.iface, args.names, *f))
	}
}

func resolveFieldNames(t *ast.Field) string {
	var result = []string{}
	for _, n := range t.Names {
		result = append(result, n.Name)
	}

	return strings.Join(result, ", ")
}

func resolveFieldTypes(t ast.Expr) string {
	switch t1 := t.(type) {
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
