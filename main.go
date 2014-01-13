package main

import (
	"flag"
	"github.com/freenerd/import-call-extractor/extractor"
	"log"
	"strings"
)

var (
	file     = flag.String("f", "", "a go source file to extract calls from")
	pkg      = flag.String("p", "", "a go package")
	suspects = flag.Bool("s", false, "filter output through a suspects list")
	format   = flag.String("format", "YAML", "format of output. valid values: yaml (default), csv (uniqued list of calling packages, useful with suspect list)")
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
		imports, err = extractor.FileImportCalls(*file, nil)
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

	format := strings.ToLower(*format)
	if format == "yaml" {
		extractor.PrintYAML(imports)
	} else if format == "csv" {
		extractor.PrintCSV(imports)
	} else {
		log.Printf("unknown format %s\n\n", format)
		flag.Usage()
		return
	}
}
