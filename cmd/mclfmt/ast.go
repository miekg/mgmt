package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/purpleidea/mgmt/lang/ast"
	"github.com/purpleidea/mgmt/lang/funcs/operators"
)

func Print(a any, w *LineWriter, opt Option) {
	switch a := a.(type) {
	case *ast.StmtProg:
		StmtProg(a, w, opt)
	case *ast.StmtRes:
		StmtRes(a, w, opt)
	case *ast.StmtBind:
		StmtBind(a, w, opt)
	case *ast.StmtResField:
		StmtResField(a, w, opt)
	case *ast.StmtEdge:
		StmtEdge(a, w, opt)
	case *ast.StmtEdgeHalf:
		StmtEdgeHalf(a, w, opt)
	case *ast.ExprStr:
		ExprStr(a, w, opt)
	case *ast.ExprInt:
		ExprInt(a, w, opt)
	case *ast.ExprFloat:
		ExprFloat(a, w, opt)
	case *ast.ExprBool:
		ExprBool(a, w, opt)
	case *ast.ExprVar:
		ExprVar(a, w, opt)
	case *ast.ExprCall:
		ExprCall(a, w, opt)
	default:
		panic("mclfmt: unhandled ast " + fmt.Sprintf("%T", a))
	}
}

func StmtProg(a *ast.StmtProg, w *LineWriter, opt Option) {
	for _, b := range a.Body {
		Print(b, w, opt)
	}
}

func StmtBind(a *ast.StmtBind, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "$%s =", a.Ident)
	} else {
		fmt.Fprintf(w, " $%s =", a.Ident)
	}
	Print(a.Value, w, opt)
}

func StmtRes(a *ast.StmtRes, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "%s", a.Kind)
	} else {
		fmt.Fprintf(w, " %s", a.Kind)
	}

	Print(a.Name, w, opt)

	io.WriteString(w, " {\n")
	w.Indent++
	for _, c := range a.Contents {
		Print(c, w, opt)
		io.WriteString(w, ",\n")
	}
	w.Indent--
	io.WriteString(w, "}\n")
}

func StmtEdge(a *ast.StmtEdge, w *LineWriter, opt Option) {
	// Start with new line?
	fmt.Fprintln(w)
	for i, e := range a.EdgeHalfList {
		Print(e, w, opt)
		if i < len(a.EdgeHalfList)-1 {
			fmt.Fprintf(w, " %s ", "->")
		}
	}
}

func StmtEdgeHalf(a *ast.StmtEdgeHalf, w *LineWriter, opt Option) {
	fmt.Fprintf(w, "%s%s[", strings.ToUpper(a.Kind[:1]), a.Kind[1:])
	opt.DropSpace = true
	Print(a.Name, w, opt)
	fmt.Fprintf(w, "].%s", a.SendRecv)
}

func ExprStr(a *ast.ExprStr, w *LineWriter, opt Option) {
	if opt.DropQuote {
		if opt.DropSpace {
			fmt.Fprintf(w, "%s", a.V)
		} else {
			fmt.Fprintf(w, " %s", a.V)
		}
		return
	}

	if opt.DropSpace {
		fmt.Fprintf(w, `"%s"`, a.V)
	} else {
		fmt.Fprintf(w, ` "%s"`, a.V)
	}
}

func ExprInt(a *ast.ExprInt, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "%d", a.V)
	} else {
		fmt.Fprintf(w, " %d", a.V)
	}
}

func ExprFloat(a *ast.ExprFloat, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "%g", a.V)
	} else {
		fmt.Fprintf(w, " %g", a.V)
	}
}

func ExprBool(a *ast.ExprBool, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "%t", a.V)
	} else {
		fmt.Fprintf(w, " %t", a.V)
	}
}

func ExprVar(a *ast.ExprVar, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "$%s", a.Name)
	} else {
		fmt.Fprintf(w, " $%s", a.Name)
	}
}

func StmtResField(a *ast.StmtResField, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "%s =>", a.Field)
	} else {
		fmt.Fprintf(w, " %s =>", a.Field)
	}

	Print(a.Value, w, opt)
}

func ExprCall(a *ast.ExprCall, w *LineWriter, opt Option) {
	switch a.Name {
	case operators.OperatorFuncName:
		Print(a.Args[1], w, opt)

		opt.DropQuote = true
		Print(a.Args[0], w, opt)
		opt.DropQuote = false

		Print(a.Args[2], w, opt)
	default:
		if opt.DropSpace {
			fmt.Fprintf(w, "%s(", a.Name)
		} else {
			fmt.Fprintf(w, " %s(", a.Name)
		}

		for i, arg := range a.Args {
			opt.DropSpace = false
			if i == 0 {
				opt.DropSpace = true
			}
			Print(arg, w, opt)
			if i < len(a.Args)-1 {
				fmt.Fprint(w, ",")
			}
		}
		fmt.Fprint(w, ")")
	}
}
