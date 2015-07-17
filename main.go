// goparse parses Go source code with the parser.Trace option set,
// yielding a trace of parsed productions.  This is ideal for human
// consumption and debugging, aiding in the understanding of the
// structure of Go programs as well as the behavior of the library
// parser.
//
// If a positional command line argument is specified, it will be used
// as the input file.  If no positional arguments are specified, the
// input will be read from standard in.
//
// There are other options as well that will enable other parser
// functionality.
//
//	usage: goparse [-h] [options] [file]
//	  -all-errors=false: report all errors (not just the first 10
//	                     on different lines)
//	  -declaration-errors=false: report declaration errors
//	  -imports-only=false: stop parsing after import declarations
//	  -parse-comments=false: parse comments and add them to AST
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"go/parser"
	"go/token"
)

var mode parser.Mode

func init() {
	mode = parser.Trace
	importsOnly := flag.Bool("imports-only", false, "stop parsing after import declarations")
	parseComments := flag.Bool("parse-comments", false, "parse comments and add them to AST")
	declarationErrors := flag.Bool("declaration-errors", false, "report declaration errors")
	allErrors := flag.Bool("all-errors", false, "report all errors (not just the first 10 on different lines)")
	flag.Parse()
	if (*importsOnly) {
		mode |= parser.ImportsOnly
	}
	if (*parseComments) {
		mode |= parser.ParseComments
	}
	if (*declarationErrors) {
		mode |= parser.DeclarationErrors
	}
	if (*allErrors) {
		mode |= parser.AllErrors
	}
}

func main() {
	var filename string
	var src io.Reader
	var err error

	fargs := flag.Args()
	if len(fargs) == 0 {
		filename = "-"
		src = os.Stdin
	} else if len(fargs) == 1 {
		filename = fargs[0]
		src, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("usage: %s [-h] [options] [file]\n", path.Base(os.Args[0]))
		os.Exit(2)
	}

	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, filename, src, mode)
	if err != nil {
		log.Fatal(err)
	}
}
