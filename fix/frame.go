// Copyright 2018 (as). This package was lifted from cmd/fix and
// this file (along with frame_test.go) was added to fix API changes
// made in recent updates. The changes/additions retain the same
// license as the Go programming language (BSD like)

package main

import (
	"go/ast"
	"go/token"
)

func init() {
	register(frameFix)
}

var frameFix = fix{
	name: "frame",
	date: "2018-02-03",
	f:    frame,
	desc: `fix frame.New invokation`,
}

func frame(f *ast.File) bool {
	fixed := false
	if rewriteImport(f, "github.com/as/frame/font", "github.com/as/font") {
		fixed = true
	}
	if !imports(f, "github.com/as/frame") {
		return fixed
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
