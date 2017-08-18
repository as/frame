package font

import (
	"image"
	"image/draw"
)

func StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte, bg image.Image, bgp image.Point) int {
	for _, b := range s {
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		draw.Draw(dst, r.Add(p), bg, bgp, draw.Src)
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		p.X += r.Dx()
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
		p.X += r.Dx()
	}
	return p.X
}
