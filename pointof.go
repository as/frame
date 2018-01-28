package frame

import (
	"image"
)

func (f *Frame) Grid(pt image.Point) image.Point {
	return f.grid(pt)
}

// PointOf returns the point on the frame's
// bitmap that intersects the index p.
func (f *Frame) PointOf(p int64) image.Point {
	return f.pointOf(p, f.r.Min, 0)
}

func (f *Frame) grid(pt image.Point) image.Point {
	pt.Y -= f.r.Min.Y
	pt.Y -= pt.Y % f.Font.Dy()
	pt.Y += f.r.Min.Y
	if pt.X > f.r.Max.X {
		pt.X = f.r.Max.X
	}
	return pt
}

func (f *Frame) pointOf(p int64, pt image.Point, bn int) (x image.Point) {
	for _, b := range f.Box[bn:f.Nbox] {
		pt = f.wrapMax(pt, &b)
		l := int64(b.Len())
		if p >= l {
			p -= l
			pt = f.advance(pt, &b)
			continue
		}
		if b.Nrune < 0 {
			break
		}
		pt.X += f.Font.Dx(b.Ptr[:p])
	}
	return pt
}
func (f *Frame) pointOfBox(p int64, nb int) (pt image.Point) {
	Nbox := f.Nbox
	f.Nbox = nb
	pt = f.pointOf(p, f.r.Min, 0)
	f.Nbox = Nbox
	return pt
}
