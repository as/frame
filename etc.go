package frame

import "image/draw"
import "image"
import "image/color"

func (f *Frame) ckindex(i int) int {
	nb := f.nbytes
	if i >= nb {
		i = nb - 1
	}
	if i < 0 {
		i = 0
	}
	return i
}

func (f *Frame) flushcache() {
	draw.Draw(f.cached, f.cached.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.ZP, draw.Src)
}

func (f *Frame) grow(n int) {
	f.s = append(f.s, make([]byte, n)...)
}
