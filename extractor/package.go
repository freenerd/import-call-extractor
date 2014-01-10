package extractor

import (
	"fmt"
	"go/build"
	"os"
	"path"
)

func PackageImportCalls(pkgName string) (map[string]map[string][]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working dir: %s", err)
	}

	imports := make(map[string]map[string][]string)
	imports, err = processPackage(cwd, pkgName, imports)

	if err != nil {
		return nil, fmt.Errorf("failed to process package: %s", err)
	}

	return imports, nil
}

func processPackage(root, pkgName string, imports map[string]map[string][]string) (map[string]map[string][]string, error) {
	// read package
	pkg, err := build.Import(pkgName, root, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to import package %s: %s", pkgName, err)
	}

	// Don't worry about dependencies for stdlib packages
	if pkg.Goroot {
		return imports, nil
	}

	// analyze each file in package, merge results
	for _, file := range pkg.GoFiles {
		fileImports, err := FileImportCalls(path.Join(pkg.Dir, file))
		if err != nil {
			return nil, fmt.Errorf("failed to in file %s: %s", file, err)
		}

		imports = mergeImportMaps(imports, fileImports)
	}

	// recursively extract from each imported package
	for _, imp := range pkg.Imports {
		// TODO: Don't do already analyzed packages
		imports, err = processPackage(root, imp, imports)
		if err != nil {
			return nil, fmt.Errorf("failed to process package %s:", imp)
		}
	}

	return imports, nil
}

// merge two import maps one merged output map
func mergeImportMaps(import1, import2 map[string]map[string][]string) map[string]map[string][]string {
	for imp, _ := range import1 {
		_, present := import2[imp]
		if present {
			// since import is already there, we have to look at each call's array
			for call, _ := range import1[imp] {
				_, present := import2[imp][call]
				if present {
					// appending each call to call array
					for _, e := range import1[imp][call] {
						import2[imp][call] = append(import2[imp][call], e)
					}
				} else {
					// call not there yet, use whole call array
					import2[imp][call] = import1[imp][call]
				}
			}
		} else {
			// import not there yet, use whole import map
			import2[imp] = import1[imp]
		}
	}

	return import2
}
