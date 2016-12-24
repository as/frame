package frame

import (
	"golang.org/x/image/math/f64"
	"image"
)

// Scale returns the matrix for scaling an image
func Scale(w, h float64) f64.Aff3 {
	return f64.Aff3{w, 0, 0, 0, h, 0}
}
func Translate(x, y float64) f64.Aff3 {
	return f64.Aff3{1, 0, x, 0, 1, y}
}
func ST(p image.Point, wh image.Point) f64.Aff3 {
	x, y := float64(p.X), float64(p.Y)
	w, h := float64(wh.X), float64(wh.Y)
	return f64.Aff3{
		w, 0, w * x,
		0, h, h * y,
	}
}

func (f *Frame) alignY(pt image.Point) image.Point {
	pt.Y -= pt.Y % f.FontHeight()
	return pt
}

func (f *Frame) upscale(pt image.Point) image.Point {
	pt.X *= f.Scale.X
	pt.Y *= f.Scale.Y
	return pt
}
func (f *Frame) downscale(pt image.Point) image.Point {
	pt.X /= f.Scale.X
	pt.Y /= f.Scale.Y
	return pt
}
