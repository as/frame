package frame

import (
	"github.com/as/etch"
	"github.com/as/frame/font"
	"image"
	"image/draw"
	"testing"
)

var (
	R     = image.Rect(0, 0, 232, 232)
	fsize = 16
	ft    = font.NewBasic(fsize)
)

func abtest(r image.Rectangle) (fr0, fr1 *Frame, a, b *image.RGBA) {
	a = image.NewRGBA(r)
	b = image.NewRGBA(r)
	fr0 = New(r, font.NewBasic(fsize), a, A)
	fr1 = New(r, font.NewBasic(fsize), b, A)
	return fr0, fr1, a, b
}

func abtestbg(r image.Rectangle) (fa, fb *Frame, a, b *image.RGBA) {
	fa, fb, a, b = abtest(r)
	draw.Draw(a, a.Bounds(), fa.Color.Back, image.ZP, draw.Src)
	draw.Draw(b, b.Bounds(), fb.Color.Back, image.ZP, draw.Src)
	return fa, fb, a, b
}

func TestDeleteOneChar(t *testing.T) {
	h, w, have, want := abtestbg(R)
	h.Insert(string("1"), 0)
	h.Delete(0, h.Len())
	h.Untick()
	w.Untick()
	etch.Assert(t, have, want, "TestDelete.png")
}

func TestDeleteLastLineNoNL(t *testing.T) {
	w, h, want, have := abtestbg(R)
	draw.Draw(want, want.Bounds(), w.Color.Back, image.ZP, draw.Src)
	draw.Draw(have, have.Bounds(), h.Color.Back, image.ZP, draw.Src)
	w.Insert(string("1234\ncccc\ndddd\n"), 0)
	h.Insert(string("1234\ncccc\ndddd"), 0)
	h.Delete(5, 10)
	w.Delete(5, 10)
	// We can untick because have has an extra newline
	h.Untick()
	w.Untick()
	etch.Assert(t, have, want, "TestDeleteLastLineNoNL.png")
}
