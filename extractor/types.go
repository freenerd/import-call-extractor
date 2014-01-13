package extractor

import (
	"go/build"
)

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
type Imports map[string]Calls
type Calls map[string]Occurences
type Occurences []Occurence
type Occurence struct {
	position string
	pkg      *build.Package
}

type ImportPaths map[string]string
