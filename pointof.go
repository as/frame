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
	for ; bn < f.Nbox; bn++ {
		b := &f.Box[bn]
		pt = f.wrapMax(pt, b)
		l := b.Len()
		if p < int64(l) {
			if b.Nrune > 0 {
				ptr := b.Ptr
				if p > 0{
					pt.X += f.Font.Dx(ptr[:p])
				}
//				bsb := len(ptr)
//				i := 0
//				for p > 0 {
//					if bsb == i {
//						break
//					}
//					size := 1
//					widthPx := f.Font.Dx(ptr[i : i+size])
//					i += size
//					p -= int64(size)
//					pt.X += widthPx
//					if pt.X > f.r.Max.X {
//						panic("PtOfCharPtBox")
//					}
//				}
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
