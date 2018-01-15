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
	pt.Y -= pt.Y % Dy(f.Font)
	pt.Y += f.r.Min.Y
	if pt.X > f.r.Max.X {
		pt.X = f.r.Max.X
	}
	return pt
}

func (f *Frame) pointOf(p int64, pt image.Point, bn int) (x image.Point) {
	for ; bn < f.Nbox; bn++ {
		b := &f.Box[bn]
		pt = f.wrapMax(pt, b)
		l := b.Len()
		if p < int64(l) {
			if b.Nrune > 0 {
				br := f.newRulerFunc(b.Ptr, f.Font)
				for p > 0 {
					size, width, err := br.Next()
					if err != nil {
						break
					}
					p -= int64(size)
					pt.X += width
					if pt.X > f.r.Max.X {
						panic("PtOfCharPtBox")
					}
				}
			}
			break
		}
		p -= int64(l)
		pt = f.advance(pt, b)
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
