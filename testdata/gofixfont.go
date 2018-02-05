// This is the raw driver file for gofix, to run it, put it in %GOROOT%\src\cmd\fix\
// and then rebuild the `go fix` tool.
package main

import (
	"go/ast"
	"go/token"
)

func init() {
	register(frameFix)
}

var frameFix = fix{
	name: "font",
	date: "2018-02-03",
	f:    frame,
	desc: `fix font.Font -> font.Face and so forth`,
}

func frame(f *ast.File) bool {
	fixed := false
	if !imports(f, "github.com/as/font") {
		return false
	}
	walk(f, func(n interface{}) {
		fcall, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}
		se, ok := fcall.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		pkgname, ok := se.X.(*ast.Ident)
		if !ok {
			return
		}
		if pkgname.Name != "frame" {
			return
		}
		x := se.Sel
		if x.Name != "New" {
			return
		}
		args := fcall.Args
		var flags ast.Expr
		switch len(args) {
		default:
			flags = args[4]
			fallthrough
		case 4:
			dst := args[2]
			r := args[0]
			ft := args[1]
			col := args[3]
			fcall.Args = []ast.Expr{dst, r, mkConfig(ft, col, flags)}
			fixed = true
		case 3:
		case 2:
		case 1:
		case 0:
		}
	})

	return fixed
}

func colorFix(a ast.Expr) ast.Expr {
	return a
}

func mkConfig(ft, col, flags ast.Expr) (exp *ast.UnaryExpr) {
	list := []ast.Expr{
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "Font",
			},
			Value: ft,
		},
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "Color",
			},
			Value: col,
		},
	}
	if flags != nil {
		list = append(list, &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "Flag",
			},
			Value: flags,
		})
	}
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "frame",
				},
				Sel: &ast.Ident{
					Name: "Config",
				},
			},
			Elts: list,
		},
	}
}
