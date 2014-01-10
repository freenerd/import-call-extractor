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

- In package mode, if a package is imported several times, it will be analyzed again
- Calls made in a `var` block will not be extracted
