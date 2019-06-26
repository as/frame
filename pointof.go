package frame

import (
	"image"
)

// Grid returns a grid-aligned point on the frame relative to pt
func (f *Frame) Grid(pt image.Point) image.Point {
	if f == nil {
		return image.ZP
	}
	return pt26point(f.grid(pt26(pt)))
}

// PointOf returns the point on the closest to index p.
func (f *Frame) PointOf(p int64) image.Point {
	if f == nil {
		return image.ZP
	}
	return pt26point(f.point0(p))
}

func (f *Frame) point0(p int64) p26 {
	return f.pointOf(p, f.r.Min, 0)
}

func (f *Frame) grid(pt p26) p26 {
	pt.Y -= f.r.Min.Y
	pt.Y -= pt.Y % f.dy()
	pt.Y += f.r.Min.Y
	if pt.X > f.r.Max.X {
		pt.X = f.r.Max.X
	}
	return pt
}
func (f *Frame) pointOf(p int64, pt p26, bn int) (x p26) {
	for ; bn < f.Nbox; bn++ {
		b := &f.Box[bn]
		pt = f.wrapMax(pt, b)
		l := b.Len()
		if p < int64(l) {
			if b.Nrune > 0 {
				ptr := b.Ptr
				if p > 0 {
					pt.X += f.Dx(ptr[:p])
				}
			}
			break
		}
		p -= int64(l)
		pt = f.advance(pt, b)
	}
	return pt
}

func (f *Frame) pointOfBox(p int64, nb int) (pt p26) {
	Nbox := f.Nbox
	f.Nbox = nb
	pt = f.pointOf(p, f.r.Min, 0)
	f.Nbox = Nbox
	return pt
}
