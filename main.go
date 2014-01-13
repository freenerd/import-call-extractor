package main

import (
	"flag"
	"github.com/freenerd/import-call-extractor/extractor"
	"log"
)

var (
	file     = flag.String("f", "", "a go source file to extract calls from")
	pkg      = flag.String("p", "", "a go package")
	suspects = flag.Bool("s", false, "filter output through a suspects list")
)

func main() {
	flag.Parse()

	if *file == "" && *pkg == "" {
		log.Printf("need go source file path or go package to continue\n\n")
		flag.Usage()
		return
	}

	imports := extractor.Imports{}
	var err error

	if *file != "" {
		imports, err = extractor.FileImportCalls(*file)
	} else if *pkg != "" {
		imports, err = extractor.PackageImportCalls(*pkg)
	}
	if err != nil {
		log.Fatal(err)
		return
	}

	if *suspects {
		imports, err = extractor.FilterForSuspectPackages(imports)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	extractor.PrintYAML(imports)
}
