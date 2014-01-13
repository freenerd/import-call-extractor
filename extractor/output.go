package extractor

import (
	"fmt"
	"strings"
)

func PrintYAML(imports Imports) {
	for imp, calls := range imports {
		fmt.Printf("%s:\n", imp)

		for call, occurences := range calls {
			fmt.Printf("  %s:\n", call)

			for _, occurence := range occurences {
				fmt.Printf("    - %s\n", occurence.position)
			}
		}
	}
}

func PrintCSV(imports Imports) {
	pkgs := extractUniquePackages(imports)
	fmt.Println(strings.Join(pkgs, ", "))
}

func extractUniquePackages(imports Imports) []string {
	pkgs := []string{}

	for _, calls := range imports {
		for _, occurences := range calls {
			for _, occurence := range occurences {
				uniqueAppendToArray(&pkgs, occurence.pkg.ImportPath)
			}
		}
	}

	return pkgs
}

func uniqueAppendToArray(array *[]string, element string) {
	for _, v := range *array {
		if v == element {
			return
		}
	}

	*array = append(*array, element)
}
