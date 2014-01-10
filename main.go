package main

import (
	"flag"
	"github.com/freenerd/import-call-extractor/extractor"
	"log"
)

var (
	file = flag.String("f", "", "a go source file to extract calls from")
	pkg  = flag.String("p", "", "a go package")
)

func main() {
	flag.Parse()

	if *file == "" && *pkg == "" {
		log.Fatal("need go source file path or go package to continue")
	}

	if *file != "" {
		imports, err := extractor.FileImportCalls(*file)
		if err != nil {
			log.Fatal(err)
			return
		}
		extractor.PrintYAML(imports)
	} else if *pkg != "" {
		imports, err := extractor.PackageImportCalls(*pkg)
		if err != nil {
			log.Fatal(err)
			return
		}
		extractor.PrintYAML(imports)
	}
}
