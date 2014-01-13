# import-call-extractor

import-call-extractor is a program to analyze go programs. It extracts calls made to imported packages. It can work either on a single go source file or a go package. On a go package, it will also recursively analyze the imported packages.

## Install

    go get github.com/freenerd/import-call-extractor

## Use

For basic usage give the absolute go source file path as first argument. In this example, we analyze the code of this package:

    import-call-extractor $GOPATH/src/github.com/freenerd/import-call-extractor/main.go

The output goes to STDOUT and is formated as yaml. It looks like this:

```yaml
flag:
  Parse:
    - /Users/johan/Code/go/src/github.com/freenerd/import-call-extractor/main.go:87:3
  Args:
    - /Users/johan/Code/go/src/github.com/freenerd/import-call-extractor/main.go:88:10
fmt:
  Printf:
    - /Users/johan/Code/go/src/github.com/freenerd/import-call-extractor/main.go:108:3
    - /Users/johan/Code/go/src/github.com/freenerd/import-call-extractor/main.go:111:4
    - /Users/johan/Code/go/src/github.com/freenerd/import-call-extractor/main.go:114:5
```

To analyze a whole package, call like this:

    import-call-extractor -p github.com/freenerd/import-call-extractor

Output can be filtered by specific suspect package calls (TODO: make customizable)

    import-call-extractor -p github.com/freenerd/import-call-extractor -s

## Limitations

- It is assumed that all calls are made on the object with the package name string split after last "/". This is a casual convention and not how go imports are actually implemented. Therefore anything exported publicly by a package, that is not within the packages name, will not be detected.
  - This will be detected: `import "flag"; ...; flag.Parse()`
  - This will not be detected: `import "flag"; ...; flags := NewFlagSet(...)`
- In package mode, if a package is imported several times, it will be analyzed each time, resulting in duplicate call occurences
- Calls made in a `var` block will not be detected
- Overriding objects will not be respected
  - `import "fmt"; ...; fmt := "anything";`

## Todo

- How are package renames handled?
- Handle all publicl exports by a package
- Regarding finding network connections: Also within the stdlib, try to traceback network connections to `net.Dial`

