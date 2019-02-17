package frame

import (
	"image"
	"image/draw"

	. "github.com/as/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font"
)

// Drawer implements the set of methods a frame needs to draw on a draw.Image. The frame's default behavior is to use
// the native image/draw package and x/exp/font packages to satisfy this interface.
type Drawer interface {
	MaxFit(p []byte, dx i26) int
	Dx(p []byte) i26
	Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, op draw.Op)
}

func NewDefaultDrawer() Drawer {
	return &defaultDrawer{}
}

type defaultDrawer struct{
	drawer
}

func (d *defaultDrawer) Flush(r ...image.Rectangle) error {
	return nil
}

type drawer struct {
	Dst    draw.Image
	Src    image.Image
	Face   font.Face
	Dot    fixed.Point26_6
	height fixed.Int26_6
	width  fixed.Int26_6
	glyphDesc
}

func (d *drawer) Reset(dst draw.Image, dot fixed.Point26_6, src image.Image) {
	d.Dst = dst
	d.Src = src
	d.Dot = dot
}

func (d *drawer) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, op draw.Op) {
	draw.Draw(dst, r, src, sp, op)
}

func newDrawer(dst draw.Image, src image.Image, face font.Face, dot p26) *drawer {
	h := face.Metrics().Height
	d := &drawer{
		Dst:    dst,
		Src:    src,
		Face:   face,
		height: h,
	}
	d.move(dot)
	return d
}

func (d *drawer) Move(p image.Point) {
	d.move(fixed.P(p.X, p.Y))
}
func (d *drawer) move(p p26) {
	d.Dot = p
	d.Dot.Y += d.height
}

func (d *drawer) Dy() i26 {
	return d.height + d.height/2
}
func (d *drawer) Dx(p []byte) (dx i26) {
	if d == nil || d.Face == nil{
		return 0
	}
	for _, c := range p {
		w, _ := d.Face.GlyphAdvance(rune(c))
		dx += w
	}
	return dx
}
func (d *drawer) MaxFit(p []byte, dx i26) (n int) {
	var c byte
	for n, c = range p {
		w, _ := d.Face.GlyphAdvance(rune(c))
		dx -= w
		if dx < 0 {
			return n
		}
	}
	return n
}
func (d *drawer) DrawBG(bg image.Image, maxx int) {
	p := image.Point{d.Dot.X.Ceil(), (d.Dot.Y-d.height).Ceil()}
	r := image.Rectangle{Min: p, Max: p}
	h := d.height.Ceil()
	dy := h + h/2
	r.Max.Y += dy+1
	r.Max.X = maxx
	draw.Draw(d.Dst, r, bg, image.ZP, draw.Src)
}
func (d *drawer) drawBG(bg image.Image, maxx i26) {
	d.DrawBG(bg, maxx.Ceil())
}

func (d *drawer) WriteString(text string) (n int, err error) {
	c := rune(-1)
	x := d.Dot.X
	for n, c = range text {
		d.glyph(c)
		draw.DrawMask(d.Dst, d.dr, d.Src, image.Point{}, d.mask, d.maskp, draw.Over)
		d.Dot.X += d.advance
	}
	d.width = d.Dot.X - x
	return
}

func (d *drawer) Write(p []byte) (n int, err error) {
	c := byte(0)
	x := d.Dot.X
	for n, c = range p {
		d.glyph(rune(c))
		draw.DrawMask(d.Dst, d.dr, d.Src, image.Point{}, d.mask, d.maskp, draw.Over)
		d.Dot.X += d.advance
	}
	d.width = d.Dot.X - x
	return
}

func (d *drawer) glyph(c rune) {
	d.dr, d.mask, d.maskp, d.advance, d.ok = d.Face.Glyph(d.Dot, c)
	d.r = c
}

type glyphDesc struct {
	dr      image.Rectangle
	mask    image.Image
	maskp   image.Point
	advance fixed.Int26_6
	ok      bool
	r       rune
}

func negotiateFace(f font.Face, flags int) Face {
	if flags&FrUTF8 != 0 {
		return NewCache(NewRune(f))
	}
	switch f := f.(type) {
	case Face:
		return f
	case font.Face:
		return Open(f)
	}
	return Open(f)
}


