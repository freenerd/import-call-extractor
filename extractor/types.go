package extractor

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
type Imports map[string]map[string][]string
