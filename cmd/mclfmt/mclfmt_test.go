package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/purpleidea/mgmt/lang/parser"
)

func TestPrint(t *testing.T) {
	testcases := []struct {
		name string
		code string
	}{
		{
			name: "addition",
			code: `
		test "t1" {
			int64ptr => 13 + 42,
		}
		`,
		},
		{
			name: "multiple float addition",
			code: `
			test "t1" {
				float32 => -25.38789 + 32.6 + 13.7,
			}
			`,
		},
		{
			name: "func call dotted 1",
			code: `
			$x1 = pkg.foo1()
			`,
		},
		{
			name: "func call dotted 2",
			code: `
			$x1 = pkg.foo1(true, "hello")
			`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			prog, err := parser.LexParse(strings.NewReader(tc.code))
			if err != nil {
				t.Fatal(err)
			}
			lw := &LineWriter{Indent: 0, Start: true, b: &bytes.Buffer{}}

			Print(prog, lw, Option{})
			fmt.Println(lw.String())
		})
	}
}
