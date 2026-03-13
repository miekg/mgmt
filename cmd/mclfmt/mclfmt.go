package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/purpleidea/mgmt/lang/ast"
	"github.com/purpleidea/mgmt/lang/parser"
)

func main() {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	xast, err := parser.LexParse(f)
	if err != nil {
		log.Fatal(err)
	}

	prog, ok := xast.(*ast.StmtProg)
	if !ok {
		log.Fatal("AST should start with *ast.StmtProg")
	}
	lw := &LineWriter{Indent: 0, Start: true, b: &bytes.Buffer{}}

	Print(prog, lw, Option{})
	fmt.Println(lw.String())
}

type Option struct {
	DropQuote bool // Print ExprStr without quotes.
	DropSpace bool // Print any without a starting space.
}

// LineWrtiter is a writer that ignores a single space written as the first
// character on the line, but will apply any indentation that is required.
type LineWriter struct {
	Indent int
	Start  bool // set at start and after newline

	b *bytes.Buffer
}

func (lw *LineWriter) Write(p []byte) (int, error) {
	if lw.Start {
		p = bytes.TrimLeft(p, " ")
		lw.b.Write(bytes.Repeat([]byte("\t"), lw.Indent))
	}
	lw.b.Write(p)
	lw.Start = bytes.HasSuffix(p, []byte("\n"))
	return len(p), nil
}

func (lw *LineWriter) String() string {
	return lw.b.String()
}
