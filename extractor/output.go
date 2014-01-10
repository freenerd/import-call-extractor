package extractor

import (
	"fmt"
)

func PrintYAML(imports map[string]map[string][]string) {
	for imp, calls := range imports {
		fmt.Printf("%s:\n", imp)

		for call, occurences := range calls {
			fmt.Printf("  %s:\n", call)

			for _, occurence := range occurences {
				fmt.Printf("    - %s\n", occurence)
			}
		}
	}
}
