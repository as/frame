package font

import (
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"unicode/utf8"
)

type rgba struct{
	r,g,b,a uint32
}
type signature struct{
	s string
	dy int
	rgba
}
var cache map[signature]*image.RGBA
func init(){
	cache = make(map[signature]*image.RGBA)
}

func StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte, bg image.Image, bgp image.Point) int {
	quad := rgba{}
	{
		r,g,b,a := src.(*image.Uniform).RGBA()
		quad=rgba{r,g,b,a}
	}
	sig := signature{
		s: string(s),
		dy: ft.Dy(),
		rgba: quad,
	}
	if img, ok := cache[sig]; ok{
		draw.Draw(dst, img.Bounds().Add(p), img, img.Bounds().Min, draw.Src)
		return p.X+img.Bounds().Dx() //Add(image.Pt(img.Bounds().Dx(), 0))
	}
	r0 := image.Rectangle{p, p}
	for _, b := range s {
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		if r0.Min == image.ZP{
			r0.Min = r.Add(p).Min
		}
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		p.X += r.Dx() + ft.stride
		r0.Max.X += r.Dx() + ft.stride
		if r.Dy() > r0.Dy(){
			r0.Max.Y=r.Dy()
		}
	}
	img := image.NewRGBA(image.Rect(0,0,r0.Max.X,r0.Max.Y))
	draw.Draw(img, img.Bounds(), bg, bgp, draw.Src)
	draw.Draw(img, img.Bounds(), dst, r0.Min, draw.Src)
	cache[sig]=img
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
