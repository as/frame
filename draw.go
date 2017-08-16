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

// Refresh renders the entire frame.
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

// Redraw0 renders the frame's bitmap at pt.
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

// Redraw draws the range [p0:p1] at the given pt.
func (f *Frame) Redraw(pt image.Point, p0, p1 int64, issel bool) {
	if f.Ticked {
		f.tickat(f.PointOf(f.p0), false)
	}

	if p0 == p1 {
		f.tickat(pt, issel)
		return
	}

	pal := f.Color.Palette
	if issel {
		pal = f.Color.Hi
	}
	f.drawsel(pt, p0, p1, pal.Back, pal.Text)
}

// Recolor redraws the range p0:p1 with the given palette
func (f *Frame) Recolor(pt image.Point, p0, p1 int64, cols Palette) {
	f.drawsel(pt, p0, p1, cols.Back, cols.Text)

}

// Put
func (f *Frame) tickat(pt image.Point, ticked bool) {
	if f.Ticked == ticked || f.tick == nil || !pt.In(f.r) {
		return
	}
	//pt.X--
	r := f.tick.Bounds().Add(pt)
	if r.Max.X > f.r.Max.X {
		r.Max.X = f.r.Max.X
	} //
	adj := image.Pt(0, -(f.Font.height / 6))
	if ticked {
		f.Draw(f.tickback, f.tickback.Bounds(), f.b, pt.Add(adj), draw.Src)
		f.Draw(f.b, r.Add(adj), f.tick, image.ZP, draw.Over)
	} else {
		f.Draw(f.b, r.Add(adj), f.tickback, image.ZP, draw.Over)
	}
	f.Ticked = ticked
}

func (f *Frame) drawAt(pt image.Point) image.Point {
	n := 0
	for nb := 0; nb < f.Nbox; nb++ {
		b := &f.Box[nb]
		pt = f.lineWrap0(pt, b)
		if pt.Y == f.r.Max.Y {
			f.Nchars -= f.Count(nb)
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
			//print("before wrap ")
			//println(qt.String())
			//print("after wrap ")
			//println(pt.String())
			if qt.X > f.r.Max.X {
				//println(f.r.Max.String())
				qt.X -= 5000
				//f.DumpBoxes()
			}
			if pt.Y > qt.Y {
				f.Draw(f.b, image.Rect(qt.X, qt.Y, f.r.Max.X, pt.Y), back, qt, f.op)
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
			if nr > -1 {
				w = f.Font.MeasureBytes(ptr[:nr])
			}  
		}
		x = pt.X + w
		if x > f.r.Max.X {
			x = f.r.Max.X
		}
		f.Draw(f.b, image.Rect(pt.X, pt.Y, x, pt.Y+f.Font.height), back, pt, f.op)
		if b.Nrune > 0 {
			//TODO: must be stringnbg....
			if nr <= int64(len(ptr)) && nr >= 0{
				f.stringbg(f.b, pt, text, image.ZP, f.Font, ptr[:nr], back, image.ZP)
			}
		}
		pt.X += w
	Continue:
		if nb+1 >= f.Nbox {
			break
		}
		b = &f.Box[nb+1]
		p += nr
	}

	if p1 > p0 && nb != 0 && nb != f.Nbox && (&f.Box[nb-1]).Nrune > 0 && !trim {
		qt = pt
		pt = f.lineWrap(pt, b)
		if qt.X > f.r.Max.X {
			//println(f.r.Max.String())
			//qt.X-=5000
			//f.DumpBoxes()
		}
		if pt.Y > qt.Y {
			f.Draw(f.b, image.Rect(qt.X, qt.Y, f.r.Max.X, pt.Y), back, qt, f.op) 
		}
	}
	return pt
}

func (f *Frame) renderDec() {
	if f.hexFont == nil {
		x := NewTTF(gomono.TTF, f.Dy()/4+3)
		f.hexFont = &x
	}
	f.hex = make([]draw.Image, 256)
	for i := range f.hex {
		sizer := f.Font.measureHex()
		f.hex[i] = image.NewRGBA(image.Rect(0, 0, sizer, f.Dy()))
		s := fmt.Sprintf("%03d", i)
		pt := image.Pt(sizer-((MeasureBytes(*f.hexFont, s)+sizer)/2), 0)
		stringnbg(f.hex[i], pt, f.Color.Text, image.ZP, *f.hexFont, []byte(s),
			image.NewUniform(color.RGBA{0, 0, 0, 255}), image.ZP)
	}
}

func (f *Frame) renderHex() {
	if f.hexFont == nil {
		x := NewTTF(gomono.TTF, f.Dy()/2+3)
		f.hexFont = &x
	}
	f.hex = make([]draw.Image, 256)
	for i := range f.hex {
		sizer := f.Font.measureHex()
		f.hex[i] = image.NewRGBA(image.Rect(0, 0, sizer, f.Dy()))
		s := fmt.Sprintf("%02x", i)
		pt := image.Pt(sizer-((MeasureBytes(*f.hexFont, s)+sizer)/2), 0)
		stringnbg(f.hex[i], pt, f.Color.Text, image.ZP, *f.hexFont, []byte(s),
			image.NewUniform(color.RGBA{0, 0, 0, 255}), image.ZP)
	}
}

func (f *Frame) stringbg(dst draw.Image, p image.Point, src image.Image,
	sp image.Point, font Font, s []byte, bg image.Image, bgp image.Point) int {
	h := font.height
	h = int(float64(h) - float64(h)/float64(5))
	for _, v := range s {
		fp := fixed.P(p.X, p.Y)
		dr, mask, maskp, _, ok := font.Glyph(fp, rune(v))
		//		if v == 0 {
		//			dr, mask, maskp, _, ok = font.Glyph(fp, 1)
		//		}
		if !ok {
			panic("Frame.stringbg")
		}
		dr.Min.Y += h
		dr.Max.Y += h
		//src = image.NewUniform(Rainbow)

		advance := f.Font.Measure(rune(v))
		if v == 0 || !unicode.IsGraphic(rune(v)) || v > 127 {
			if len(f.hex) == 0 {
				f.renderHex()
			}
			dr, _, _, _, _ := font.Glyph(fp, rune('@'))
			dr.Max.X = dr.Min.X + advance
			dr.Min.Y += font.height - 5
			dr.Max.Y += font.height - 5
			//draw.Draw(dst, dr, bg, bgp, draw.Src)
			draw.Draw(dst, dr, f.hex[byte(v)], bgp, draw.Over)
		} else {
			//draw.Draw(dst, dr, bg, bgp, draw.Src)
			//next()
			//
			draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		}
		//next()
		p.X += advance
	}
	return int(p.X)
}

func MeasureBytes(f Font, p string) (w int) {
	for i := range p {
		w += measure(f, rune(byte(p[i])))
	}
	return w
}
func measure(f Font, r rune) int {
	l, ok := f.Face.GlyphAdvance(r)
	if !ok {
		println("warn: glyph missing")
		l, _ = f.Face.GlyphAdvance('@')
	}
	return fix(l)
}

func stringbg(dst draw.Image, p image.Point, src image.Image,
	sp image.Point, font Font, s []byte, bg image.Image, bgp image.Point) int {
	h := font.height
	h = int(float64(h) - float64(h)/float64(5))
	for _, v := range s {
		fp := fixed.P(p.X, p.Y)
		dr, mask, maskp, _, ok := font.Glyph(fp, rune(v))
		if !ok {
			panic("stringbg")

		}
		dr.Min.Y += h
		dr.Max.Y += h
		//src = image.NewUniform(Rainbow)
		draw.Draw(dst, dr, bg, bgp, draw.Src)
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		//next()
		p.X += measure(font, rune(v)) //fix(advance)
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
