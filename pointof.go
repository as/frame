package frame

import (
	"github.com/as/frame/box"
	"image"
)

func (f *Frame) PtOfCharPtBox(p int64, pt image.Point, bn int) (x image.Point) {
	var (
		b    *box.Box
		l, w int
		//r rune
	)
	for ; bn < f.Nbox; bn++ {
		b = &f.Box[bn]
		pt = f.LineWrap(pt, b)
		l = b.Len()
		if p < int64(l) {
			if b.Nrune > 0 {
				for s := b.Ptr; p > 0; p, s = p-1, s[w:] {
					// TODO: runes
					w = 1
					pt.X += f.Font.stringwidth(s[:1])
					if pt.X > f.r.Max.X {
						panic("PtOfCharPtBox")
					}
				}
			}
			break
		}
		p -= int64(l)
		pt = f.Advance(pt, b)
	}
	return pt
}
func (f *Frame) PtOfCharNBox(p int64, nb int) (pt image.Point) {
	Nbox := f.Nbox
	f.Nbox = nb
	pt = f.PtOfCharPtBox(p, f.r.Min, 0)
	f.Nbox = Nbox
	return pt
}

func (f *Frame) PtOfChar(p int64) image.Point {
	return f.PtOfCharPtBox(p, f.r.Min, 0)

}
func (f *Frame) Grid(pt image.Point) image.Point {
	pt.Y -= f.r.Min.Y
	pt.Y -= pt.Y % f.Font.height
	pt.Y += f.r.Min.Y
	if pt.X > f.r.Max.X {
		pt.X = f.r.Max.X
	}
	return pt
}
func (f *Frame) IndexOf(pt image.Point) int64 {
	pt = f.Grid(pt)
	qt := f.r.Min
	p := int64(0)
	bn := 0
	for ; bn < f.Nbox && qt.Y < pt.Y; bn++ {
		b := &f.Box[bn]
		qt = f.LineWrap(qt, b)
		if qt.Y >= pt.Y {
			break
		}
		qt = f.Advance(qt, b)
		p += int64(b.Len())
	}

	for ; bn < f.Nbox && qt.X <= pt.X; bn++ {
		b := &f.Box[bn]
		qt = f.LineWrap(qt, b)
		if qt.Y > pt.Y {
			break
		}
		if qt.X+b.Width > pt.X {
			if b.Nrune < 0 {
				qt = f.Advance(qt, b)
			} else {
				s := b.Ptr
				for {
					if s == nil {
						panic("Frame.IndexOf: s == nil")
					}
					r := s[0]
					//TODO: rune
					w := 1
					if r == 0 {
						//println("calm panic: nul in string")
					}
					qt.X += f.Font.stringwidth(s[:1])
					s = s[w:]
					if qt.X > pt.X {
						break
					}
					p++
				}
			}
		} else {
			p += int64(b.Len())
			qt = f.Advance(qt, b)
		}
	}
	return p
}
