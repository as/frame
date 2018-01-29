package frame

import (
	"image"
	"testing"
)

type indexTest struct {
	insert string
	p0     int64
	c0     int64
	pt     image.Point
}

func TestIndexOf(t *testing.T) {
	tab := []pointTest{
		{"hello1", 0, 0, image.Pt(0, 0)},
		{"hello2", 0, 1, image.Pt(7, 0)},
		{"he\nsaid2", 0, 3, image.Pt(0, 16)},
		{"he\n\n\n\nsaid2", 0, 4, image.Pt(0, 32)},
		//		{"hello3", 2, image.Pt(0, 0)},
		//		{"hello4", 3, image.Pt(0, 0)},
		//		{"hello5", 4, image.Pt(0, 0)},
	}
	for _, v := range tab {
		h, _, _, _ := abtestbg(R)
		h.Insert([]byte(v.insert), v.p0)
		have := h.IndexOf(v.pt)
		want := v.c0
		if have != want {
			t.Logf("%q: have %d, want %d", v.insert, have, want)
			t.Fail()
		}
	}
}

func TestIndexOfMultiInsert(t *testing.T) {
	t.Skip("not finished")
	type pointTest struct {
		s  string
		p0 int64
		pt image.Point
	}
	prog := `package main
import "fmt"

func main(){
	fmt.Println("take me to your leader")
}
`

	tab := []pointTest{
		{"package(sp)", 0, image.Pt(0, 0)},
		{"package(sp)", 1, image.Pt(7, 0)},
		{"package(ep)", 7, image.Pt(7*7, 0)},
		{"main(sp)", 7 + 1, image.Pt(7*(7+1), 0)},
		{"main(ep)", 7 + 1 + 4, image.Pt(7*(7+1+4), 0)},
		{"nl(1)", 7 + 1 + 4 + 1, image.Pt(0, 16)},
		{"import(sp)", 7 + 1 + 4 + 1 + 1, image.Pt(0, 16)},
		{"import(sp+1)", 7 + 1 + 4 + 1 + 1 + 1, image.Pt(0, 16)},
	}
	for _, v := range tab {
		h, _, _, _ := abtestbg(R)
		for i, c := range []byte(prog) {
			h.Insert([]byte{c}, int64(i))
		}
		have := h.IndexOf(v.pt)
		want := v.p0
		if have != want {
			t.Logf("%q: have %d, want %d", v.s, have, want)
			t.Fail()
		}
	}
}
