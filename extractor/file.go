package extractor

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"strings"
)

func FileImportCalls(file string, pkg *build.Package) (Imports, error) {
	v, err := newVisitor(file, pkg)
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
	Imports     Imports
	ImportPaths ImportPaths // to retain full import paths for display
	fileAst     *ast.File
	fset        *token.FileSet
	pkg         *build.Package
}

func newVisitor(file string, pkg *build.Package) (*visitor, error) {
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, file, nil, parser.ParseComments)

	if err != nil {
		return nil, err
	}

	v := visitor{
		Imports:     Imports{},
		ImportPaths: ImportPaths{},
		fileAst:     fileAst,
		fset:        fset,
		pkg:         pkg,
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
					pkg := pathSplit[len(pathSplit)-1]

					v.Imports[pkg] = Calls{}

					// to reconstruct the path from the package names, we save them to importPaths
					v.ImportPaths[pkg] = path
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
					v.Imports[object][call] = Occurences{}
				}

				occurence := Occurence{
					position: token.Position.String((v.fset.Position(t.Pos()))),
					pkg:      v.pkg,
				}

				v.Imports[object][call] = append(v.Imports[object][call], occurence)
			}
		}
	}

	return v
}

func reconstructImportPaths(imports *Imports, importPaths ImportPaths) (Imports, error) {
	output := Imports{}

	for imp, calls := range *imports {
		path, present := importPaths[imp]

		if !present {
			return nil, fmt.Errorf("import package name %s not in import paths", imp)
		}

		output[path] = calls
	}

	return output, nil
}
