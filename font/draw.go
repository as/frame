package font

import (
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"unicode/utf8"
)

func StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte, bg image.Image, bgp image.Point) int {
	for _, b := range s {
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		//draw.Draw(dst, r.Add(p), bg, bgp, draw.Src)
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		p.X += r.Dx() + ft.stride
	}
	return p.X
}

func StringNBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte) int {
	for _, b := range s {
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		p.X += r.Dx() + ft.stride
	}
	return p.X
}

func RuneBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte, bg image.Image, bgp image.Point) int {
	p.Y += ft.Size()
	for {
		b, size := utf8.DecodeRune(s)
		dr, mask, maskp, advance, ok := ft.Glyph(fixed.P(p.X, p.Y), b)
		if !ok {
			panic("RuneBG")
		}
		draw.Draw(dst, dr, bg, bgp, draw.Src)
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
		if len(s)-size == 0 {
			break
		}
		s = s[size:]
	}
	return p.X
}

func RuneNBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte) int {
	p.Y += ft.Size()
	for {
		b, size := utf8.DecodeRune(s)
		dr, mask, maskp, advance, ok := ft.Glyph(fixed.P(p.X, p.Y), b)
		if !ok {
			panic("RuneBG")
		}
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
		if len(s)-size == 0 {
			break
		}
		s = s[size:]
	}
	return p.X
}
