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

func TestIndexOfWrap(t *testing.T) {
	h, _, _, _ := abtestbg(R)
	x := []byte("a")
	dx := h.Font.Dx(x)
	for i := 1; i < 100000; i++{
		h.Insert(x, h.Len())
		have := h.IndexOf(R.Max)
		want := int64(dx*i)// h.IndexOf(h.PointOf(int64(dx*i)))
		if have != want{
			t.Logf("round %d: have: %d want %d\n", i, have, want)
			t.Fail()
		}
		if tdx := dx*i; tdx > h.Bounds().Max.X{
			t.Logf("round %d: didnt wrap: %d\n", i, tdx)
			t.Fail()
		}
	}
}

func TestIndexOf(t *testing.T) {
	tab := []pointTest{
		{"hello1", 0, 0, image.Pt(0, 0)},
		{"hello2", 0, 1, image.Pt(7, 0)},
		{"he\nsaid2", 0, 3, image.Pt(0, 16)},
		{"he\n\n\n\nsaid2", 0, 4, image.Pt(0, 32)},
		{"\n", 0, 0, image.Pt(0, 0*8)},
		{"\n", 100, 0, image.Pt(0, 1*8)},
		{"\n\n", 0, 1, image.Pt(0, 2*8)},
		{"\n\n\n", 0, 1, image.Pt(0, 3*8)},
		{"\t\n", 0, 0, image.Pt(0, 0*8)},
		{"\t\n", 100, 1+0, image.Pt(0, 1*8)},
		{"\t\n\n", 0, 1+1, image.Pt(0, 2*8)},
		{"\t\n\n\n", 0, 1+1, image.Pt(0, 3*8)},
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
