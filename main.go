package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type Visitor struct {
}

var (
	// keeping calls grouped by import with each call occurence in an array
	// example:
	//   map[
	//     fmt:map[
	//       PrintF:[12:32 43:1]
	//       PrintLn:[10:0]
	//     ]
	//     strings:map[
	//       Replace:[11:10]
	//     ]
	//   ]
	imports map[string]map[string][]string

	importPaths map[string]string
	fset        *token.FileSet
)

func (v *Visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	// adding a key per imported package to the imports map
	case *ast.GenDecl:
		if t.Tok == token.IMPORT {
			for _, spec := range t.Specs {
				if f, ok := spec.(*ast.ImportSpec); ok {
					// i'm not sure why, but f.Path.Value returns the value in "doublequotes"
					// so let's filter them out the hard way
					path := strings.Replace(f.Path.Value, "\"", "", -1)

					// a package is locally only referenced by it's name, not by path
					pathSplit := strings.Split(path, "/")
					pkgName := pathSplit[len(pathSplit)-1]

					imports[pkgName] = make(map[string][]string)

					// to reconstruct the path from the package names, we save them to importPaths
					importPaths[pkgName] = path
				}
			}
		}
		return nil

	// A SelectorExp might be a call to an imported package
	case *ast.SelectorExpr:
		object, call := "", ""
		if i, ok := t.X.(*ast.Ident); ok {
			object = i.Name
		}
		call = t.Sel.Name

		if object != "" && call != "" {
			if _, present := imports[object]; present {
				if _, present = imports[object][call]; !present {
					// first call, make call occurence array
					imports[object][call] = make([]string, 0)
				}

				position := token.Position.String((fset.Position(t.Pos())))
				imports[object][call] = append(imports[object][call], position)
			}
		}
	}

	return v
}

func main() {
	fset = token.NewFileSet() // positions are relative to fset
	imports = make(map[string]map[string][]string)
	importPaths = make(map[string]string)

  flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		log.Fatal("need go source file path to continue")
	}

	fileast, err := parser.ParseFile(fset, args[0], nil, parser.ParseComments)

	if err != nil {
		log.Fatal(err)
		return
	}

	ast.Walk(new(Visitor), fileast)

	printImports(imports, importPaths)
}

func printImports(imports map[string]map[string][]string, importPaths map[string]string) {
	for imp, calls := range imports {
		fmt.Printf("%s:\n", importPaths[imp])

		for call, occurences := range calls {
			fmt.Printf("  %s:\n", call)

			for _, occurence := range occurences {
				fmt.Printf("    - %s\n", occurence)
			}
		}
	}
}

