package main

import (
	"flag"
	"fmt"
	"github.com/freenerd/import-call-extractor/extractor"
	"log"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		log.Fatal("need go source file path to continue")
	}

	imports, err := extractor.GetImportCalls(args[0])

	if err != nil {
		log.Fatal(err)
		return
	}

	extractor.PrintYAML(imports)
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
