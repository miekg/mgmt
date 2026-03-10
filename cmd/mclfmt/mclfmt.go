package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/purpleidea/mgmt/lang/interfaces"
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

	xast.Apply(func(node interfaces.Node) error {
		fmt.Printf("node %+v\n", node)
		fmt.Printf("node %T\n", node)
		fmt.Printf("node %s\n", node)
		return nil
	})

	fmt.Printf("behold, the AST: %+v", xast)
}
