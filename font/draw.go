package font

import (
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"unicode/utf8"
)

func StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte, bg image.Image, bgp image.Point) int {
	panic("StringBG")
	if bg == nil {
		return StringNBG(dst, p, src, sp, ft, s)
	}
	cache := ft.imgCache
	quad := rgba{}
	{
		r, g, b, a := bg.(*image.Uniform).RGBA()
		quad = rgba{r, g, b, a}
	}
	sig := signature{
		dy:   ft.Dy(),
		rgba: quad,
	}
	r0 := image.Rectangle{p, p}
	for _, b := range s {
		sig.b = b
		if img, ok := cache[sig]; ok {
			draw.Draw(dst, img.Bounds().Add(p), img, img.Bounds().Min, draw.Src)
			p.X += img.Bounds().Dx() + ft.stride //Add(image.Pt(img.Bounds().Dx(), 0))
			continue
		}
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		if r0.Min == image.ZP {
			r0.Min = r.Add(p).Min
		}
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		img := image.NewRGBA(r)
		draw.Draw(img, img.Bounds(), bg, bgp, draw.Src)
		draw.Draw(img, img.Bounds(), dst, r.Add(p).Min, draw.Src)
		cache[sig] = img

		p.X += r.Dx() + ft.stride
		r0.Max.X += r.Dx() + ft.stride
		if r.Dy() > r0.Dy() {
			r0.Max.Y = r.Dy()
		}
	}
	return p.X
}

func StringNBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte) int {
	panic("RuneNBG")
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
	panic("RuneNBG")
	p.Y += ft.Size()
	for {
		b, size := utf8.DecodeRune(s)
		dr, mask, maskp, advance, ok := ft.Glyph(fixed.P(p.X, p.Y), b)
		if !ok {
			//panic("RuneBG")
		}
		//draw.Draw(dst, dr, bg, bgp, draw.Src)
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
	panic("RuneNBG")
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
