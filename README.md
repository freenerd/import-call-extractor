# import-call-extractor

import-call-extractor is a program for extracing calls to imported packages in a go file

## Install

    go get github.com/freenerd/

## Use

For basic usage give the absolute go source file path as first argument. In this example, we analyze the code of this package:

    import-call-extractor $GOPATH/src/github.com/freenerd/import-call-extractor/main.go

The output goes to STDOUT and is formated as yaml. It looks like this:

```
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
