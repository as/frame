package frame

import (
	"github.com/as/frame/box"
	"golang.org/x/image/font/gofont/gomono"

	"fmt"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"unicode"
)

// Put
func (f *Frame) draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	draw.Draw(dst, r, src, sp, draw.Src)
	if len(f.Cache) == 0 {
		f.Cache = append(f.Cache, r)
		return
	}
	c := f.Cache[len(f.Cache)-1]
	if r.Min.X == c.Max.X || r.Max.X == c.Min.X || r.Max.Y == c.Min.Y || r.Min.Y == c.Max.Y {
		f.Cache[0] = f.Cache[0].Union(r)
	} else {
		c := f.Cache[0]
		if c.Dx()*c.Dy() < r.Dx()*r.Dy() {
			f.Cache = append([]image.Rectangle{r}, f.Cache...)
		} else {
			f.Cache = append(f.Cache, r)
		}
	}
	//f.Cache = append(f.Cache, r)
}
func (f *Frame) drawover(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	f.Cache = append(f.Cache, r)
	draw.Draw(dst, r, src, sp, draw.Over)
}
func (f *Frame) tickat(pt image.Point, ticked bool) {
	if f.Ticked == ticked || f.tick == nil || !pt.In(f.r) {
		return
	}
	pt.X--
	r := f.tick.Bounds().Add(pt)
	if r.Max.X > f.r.Max.X {
		r.Max.X = f.r.Max.X
	} //
	adj := image.Pt(0, -(f.Font.height / 6))
	if ticked {
		f.drawover(f.tickback, f.tickback.Bounds(), f.b, pt.Add(adj))
		f.drawover(f.b, r.Add(adj), f.tick, image.ZP)
	} else {
		f.drawover(f.b, r.Add(adj), f.tickback, image.ZP)
	}
	f.Ticked = ticked
}

func (f *Frame) Refresh() {
	cols := f.Color
	if f.p0 == f.p1 {
		ticked := f.Ticked
		if ticked {
			f.tickat(f.PointOf(f.p0), false)
		}
		f.drawsel(f.PointOf(0), 0, f.Nchars, cols.Back, cols.Text)
		if ticked {
			f.tickat(f.PointOf(f.p0), true)
		}
		return
	}
	pt := f.PointOf(0)
	pt = f.drawsel(pt, 0, f.p0, cols.Back, cols.Text)
	pt = f.drawsel(pt, f.p0, f.p1, cols.Hi.Back, cols.Hi.Text)
	pt = f.drawsel(pt, f.p1, f.Nchars, cols.Back, cols.Text)
}

func (f *Frame) drawAt(pt image.Point) image.Point {
	n := 0
	for nb := 0; nb < f.Nbox; nb++ {
		b := &f.Box[nb]
		pt = f.lineWrap0(pt, b)
		if pt.Y == f.r.Max.Y {
			f.Nchars -= f.Len(nb)
			f.Run.Delete(nb, f.Nbox-1)
			break
		}

		if b.Nrune > 0 {
			n = f.canFit(pt, b)
			if n == 0 {
				panic("frame: draw: cant fit shit")
			}
			if n != b.Nrune {
				f.Split(nb, n)
				b = &f.Box[nb]
			}
			pt.X += b.Width
		} else {
			if b.BC == '\n' {
				pt.X = f.r.Min.X
				pt.Y += f.Font.height
			} else {
				pt.X += f.newWid(pt, b)
			}
		}
	}
	return pt
}

func (f *Frame) Redraw0(pt image.Point, text, back image.Image) {
	nb := 0
	for ; nb < f.Nbox; nb++ {
		b := &f.Box[nb]
		pt = f.lineWrap(pt, b)
		//if !f.noredraw && b.nrune >= 0 {
		if b.Nrune >= 0 {
			f.stringbg(f.b, pt, text, image.ZP, f.Font, b.Ptr, back, image.ZP)
		}
		pt.X += b.Width
	}
}

func (f *Frame) Redraw(pt image.Point, p0, p1 int64, issel bool) {
	if f.Ticked {
		f.tickat(f.PointOf(f.p0), false)
	}

	if p0 == p1 {
		f.tickat(pt, issel)
		return
	}

	pal := f.Color.Pallete
	if issel {
		pal = f.Color.Hi
	}
	f.drawsel(pt, p0, p1, pal.Back, pal.Text)
}

func (f *Frame) drawsel(pt image.Point, p0, p1 int64, back, text image.Image) image.Point {
	p := int64(0)
	nr := p
	w := 0
	trim := false
	qt := image.ZP
	var b *box.Box
	nb := 0
	x := 0
	var ptr []byte
	for ; nb < f.Nbox && p < p1; nb++ {
		b = &f.Box[nb]
		nr = int64(b.Nrune)
		if nr < 0 {
			nr = 1
		}
		if p+nr <= p0 {
			goto Continue
		}
		if p >= p0 {
			qt = pt
			pt = f.lineWrap(pt, b)
			// fill in the end of a wrapped line
			if pt.Y > qt.Y {
				//	cache = append(cache, image.Rect(qt.X, qt.Y, f.r.Max.X, pt.Y))
				f.draw(f.b, image.Rect(qt.X, qt.Y, f.r.Max.X, pt.Y), back, qt)
			}
		}
		ptr = b.Ptr
		if p < p0 {
			ptr = ptr[p0-p:] // todo: runes
			nr -= p0 - p
			p = p0
		}

		trim = false
		if p+nr > p1 {
			nr -= (p + nr) - p1
			trim = true
		}

		if b.Nrune < 0 || nr == int64(b.Nrune) {
			w = b.Width
		} else {
			// TODO: put stringwidth back
			w = f.Font.stringwidth(ptr[:nr])
		}
		x = pt.X + w
		if x > f.r.Max.X {
			x = f.r.Max.X
		}
		f.draw(f.b, image.Rect(pt.X, pt.Y, x, pt.Y+f.Font.height), back, pt)
		if b.Nrune >= 0 {
			//TODO: must be stringnbg....
			f.stringbg(f.b, pt, text, image.ZP, f.Font, ptr[:nr], back, image.ZP)
		}
		pt.X += w
	Continue:
		b = &f.Box[nb+1]
		p += nr
	}

	if p1 > p0 && nb != 0 && nb != f.Nbox && (&f.Box[nb-1]).Nrune > 0 && !trim {
		qt = pt
		pt = f.lineWrap(pt, b)
		if pt.Y > qt.Y {
			//cache =append(cache, image.Rect(qt.X, qt.Y, f.r.Max.X, pt.Y))
			f.draw(f.b, image.Rect(qt.X, qt.Y, f.r.Max.X, pt.Y), back, qt)
		}
	}
	return pt
}

var Rainbow = color.RGBA{255, 0, 0, 255}

func next() {
	Rainbow = nextcolor(Rainbow)
}

// nextcolor steps through a gradient
func nextcolor(c color.RGBA) color.RGBA {
	switch {
	case c.R == 255 && c.G == 0 && c.B == 0:
		c.G += 25
	case c.R == 255 && c.G != 255 && c.B == 0:
		c.G += 25
	case c.G == 255 && c.R != 0:
		c.R -= 25
	case c.R == 0 && c.B != 255:
		c.B += 25
	case c.B == 255 && c.G != 0:
		c.G -= 25
	case c.G == 0 && c.R != 255:
		c.R += 25
	default:
		c.B -= 25
	}
	return c
}

func (f *Frame) renderHex() {
	if f.hexFont == nil {
		x := NewTTF(gomono.TTF, f.Dy()/4+3)
		f.hexFont = &x
	}
	f.hex = make([]draw.Image, 256)
	for i := range f.hex {
		f.hex[i] = image.NewRGBA(image.Rect(0, 0, f.Dx("_"), f.Dy()))
		pt := image.Pt(1, f.Dy()/5)
		s := fmt.Sprintf("%02d", i)
		stringnbg(f.hex[i], pt, f.Color.Text, image.ZP, *f.hexFont, []byte(s),
			f.Color.Back, image.ZP)
	}
}

func (f *Frame) stringbg(dst draw.Image, p image.Point, src image.Image,
	sp image.Point, font Font, s []byte, bg image.Image, bgp image.Point) int {
	h := font.height
	h = int(float64(h) - float64(h)/float64(5))
	for _, v := range s {
		fp := fixed.P(p.X, p.Y)
		dr, mask, maskp, advance, ok := font.Glyph(fp, rune(v))
		if v == 0 {
			dr, mask, maskp, advance, ok = font.Glyph(fp, rune(1))
		}
		if !ok {
			panic("Frame.stringbg")
			break
		}
		dr.Min.Y += h
		dr.Max.Y += h
		//src = image.NewUniform(Rainbow)
		draw.Draw(dst, dr, bg, bgp, draw.Src)
		if !unicode.IsGraphic(rune(v)) {
			if len(f.hex) == 0 {
				f.renderHex()
			}
			draw.Draw(dst, dr, f.hex[byte(v)], bgp, draw.Over)
		} else {
			draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		}

		//next()
		p.X += fix(advance)
	}
	return int(p.X)
}

func stringbg(dst draw.Image, p image.Point, src image.Image,
	sp image.Point, font Font, s []byte, bg image.Image, bgp image.Point) int {
	h := font.height
	h = int(float64(h) - float64(h)/float64(5))
	for _, v := range s {
		fp := fixed.P(p.X, p.Y)
		dr, mask, maskp, advance, ok := font.Glyph(fp, rune(v))
		if !ok {
			panic("stringbg")
			break
		}
		dr.Min.Y += h
		dr.Max.Y += h
		//src = image.NewUniform(Rainbow)
		draw.Draw(dst, dr, bg, bgp, draw.Src)
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		//next()
		p.X += fix(advance)
	}
	return int(p.X)
}

func stringnbg(dst draw.Image, p image.Point, src image.Image,
	sp image.Point, font Font, s []byte, bg image.Image, bgp image.Point) int {
	h := font.height
	h = int(float64(h) - float64(h)/float64(5))
	for _, v := range s {
		fp := fixed.P(p.X, p.Y)
		dr, mask, maskp, advance, ok := font.Glyph(fp, rune(v))
		if !ok {
			panic("stringnbg")
			break
		}
		dr.Min.Y += h
		dr.Max.Y += h
		//src = image.NewUniform(Rainbow)
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		//next()
		p.X += fix(advance)
	}
	return int(p.X)
}
