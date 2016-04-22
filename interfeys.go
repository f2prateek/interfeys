package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	typeName = flag.String("type", "", "type name; must be set")
	output   = flag.String("output", "", "output file name; default srcdir/<type>_interface.go")
)

func main() {
	flag.Parse()

	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	buildPkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}

	fs := token.NewFileSet()
	var files []*ast.File
	for _, goFile := range buildPkg.GoFiles {
		parsedFile, err := parser.ParseFile(fs, goFile, nil, 0)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, parsedFile)
	}

	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check(buildPkg.Name, fs, files, nil)
	if err != nil {
		log.Fatal(err)
	}

	target := pkg.Scope().Lookup(*typeName)
	templateData := struct {
		Package string
		Imports []string
		Type    string
		Methods []string
	}{
		pkg.Name(), make([]string, 0, 0), *typeName, make([]string, 0, 0),
	}

	methodSet := types.NewMethodSet(types.NewPointer(target.Type()))
	for i := 0; i < methodSet.Len(); i++ {
		method := methodSet.At(i)
		methodName := method.Obj().Name()
		methodType := method.Type().(*types.Signature)

		methodString := strings.Replace(methodType.String(), "func", methodName, 1)
		for _, imp := range pkg.Imports() {
			if !strings.Contains(methodString, imp.Path()) {
				continue
			}

			templateData.Imports = append(templateData.Imports, imp.Path())

			methodString = strings.Replace(methodString, imp.Path(), imp.Name()+"", -1)
			break
		}
		templateData.Methods = append(templateData.Methods, methodString)
	}

	// todo: imports
	tmpl, err := template.New("module").Parse(`package {{ .Package }}
{{range .Imports}}
import "{{ . }}"
{{- end}}

var _ {{ .Type }}Interface = (*{{ .Type }})(nil)

type {{ .Type }}Interface interface { {{- range .Methods}}
	{{ . }}{{end}}
}
`)
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		log.Fatal(err)
	}

	fileName := *output
	if fileName == "" {
		dir := filepath.Dir(".")
		fileName = filepath.Join(dir, fmt.Sprintf("%s_interface.go", strings.ToLower(*typeName)))
	}
	err = ioutil.WriteFile(fileName, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
