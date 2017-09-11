package frame

import (
	"image"
	"image/draw"
	"testing"
	"github.com/as/etch"
	"github.com/as/frame/font"
)

var (
	R = image.Rect(0, 0, 232, 232)
	fsize = 16
	ft = font.NewBasic(fsize)
)

func abtest(r image.Rectangle) (fr0, fr1 *Frame, a, b *image.RGBA) {
	a = image.NewRGBA(r)
	b = image.NewRGBA(r)
	fr0 = New(r, font.NewBasic(fsize), a, A)
	fr1 = New(r, font.NewBasic(fsize), b, A)
	return fr0, fr1, a, b
}

func TestDeleteLastLineNoNL(t *testing.T) {
	w, h, want, have := abtest(R)
	draw.Draw(want, want.Bounds(), w.Color.Back, image.ZP, draw.Src)
	draw.Draw(have, have.Bounds(), h.Color.Back, image.ZP, draw.Src)
	w.Insert([]byte("1234\ncccc\ndddd\n"), 0)
	h.Insert([]byte("1234\ncccc\ndddd"), 0)
	h.Delete(5, 10)
	w.Delete(5, 10)
	// We can untick because have has an extra newline
		h.Untick()
		w.Untick()
	etch.Assertf(t, have, want, "delta.png", "TestDeleteLastLineNoNL: failed")
}



















func abtestPad16(r image.Rectangle) (fr0, fr1 *Frame, a, b *image.RGBA) {
	a = image.NewRGBA(r)
	b = image.NewRGBA(r)
	fr0 = New(r.Inset(fsize), font.NewBasic(fsize), a, A)
	fr1 = New(r.Inset(fsize), font.NewBasic(fsize), b, A)
	return fr0, fr1, a, b
}