package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func FileImportCalls(file string) (map[string]map[string][]string, error) {
	v, err := newVisitor(file)
	if err != nil {
		return nil, err
	}

	v.walk()

	imps, err := reconstructImportPaths(&v.Imports, v.ImportPaths)
	if err != nil {
		return nil, err
	}

	return imps, nil
}

type visitor struct {
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
	Imports map[string]map[string][]string

	// retain full import paths (see Visit())
	ImportPaths map[string]string

	fileAst *ast.File
	fset    *token.FileSet
}

func newVisitor(file string) (*visitor, error) {
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, file, nil, parser.ParseComments)

	if err != nil {
		return nil, err
	}

	v := visitor{
		Imports:     make(map[string]map[string][]string),
		ImportPaths: make(map[string]string),
		fileAst:     fileAst,
		fset:        fset,
	}

	return &v, nil
}

func (v *visitor) walk() *visitor {
	ast.Walk(v, v.fileAst) // calls the Visit method for each ast node
	return v
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
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

					v.Imports[pkgName] = make(map[string][]string)

					// to reconstruct the path from the package names, we save them to importPaths
					v.ImportPaths[pkgName] = path
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
			if _, present := v.Imports[object]; present {
				if _, present = v.Imports[object][call]; !present {
					// first call, make call occurence array
					v.Imports[object][call] = make([]string, 0)
				}

				position := token.Position.String((v.fset.Position(t.Pos())))
				v.Imports[object][call] = append(v.Imports[object][call], position)
			}
		}
	}

	return v
}

func reconstructImportPaths(imports *map[string]map[string][]string, importPaths map[string]string) (map[string]map[string][]string, error) {
	output := make(map[string]map[string][]string)

	for imp, calls := range *imports {
		path, present := importPaths[imp]

		if !present {
			return nil, fmt.Errorf("import package name %s not in import paths", imp)
		}

		output[path] = calls
	}

	return output, nil
}
