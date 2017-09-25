package frame

import (
	"image"
)

func (f *Frame) ptOfCharPtBox(p int64, pt image.Point, bn int) (x image.Point) {
	for ; bn < f.Nbox; bn++ {
		b := &f.Box[bn]
		pt = f.lineWrap(pt, b)
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
func (f *Frame) ptOfCharNBox(p int64, nb int) (pt image.Point) {
	Nbox := f.Nbox
	f.Nbox = nb
	pt = f.ptOfCharPtBox(p, f.r.Min, 0)
	f.Nbox = Nbox
	return pt
}

func (f *Frame) PointOf(p int64) image.Point {
	return f.ptOfCharPtBox(p, f.r.Min, 0)
}
func (f *Frame) Grid(pt image.Point) image.Point {
	return f.grid(pt)
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
func (f *Frame) IndexOf(pt image.Point) int64 {
	pt = f.grid(pt)
	qt := f.r.Min
	p := int64(0)
	bn := 0
	for ; bn < f.Nbox && qt.Y < pt.Y; bn++ {
		b := &f.Box[bn]
		qt = f.lineWrap(qt, b)
		if qt.Y >= pt.Y {
			break
		}
		qt = f.advance(qt, b)
		p += int64(b.Len())
	}

	for ; bn < f.Nbox && qt.X <= pt.X; bn++ {
		b := &f.Box[bn]
		qt = f.lineWrap(qt, b)
		if qt.Y > pt.Y {
			break
		}
		if qt.X+b.Width > pt.X {
			if b.Nrune < 0 {
				qt = f.advance(qt, b)
			} else {
				bs := f.newRulerFunc(b.Ptr, f.Font)
				for {
					size, width, err := bs.Next()
					if err != nil {
						break
					}
					qt.X += width
					if qt.X > pt.X {
						break
					}
					p += int64(size)
				}
			}
		} else {
			p += int64(b.Len())
			qt = f.advance(qt, b)
		}
	}
	return p
}
