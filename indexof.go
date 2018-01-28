package frame

import (
	"image"
)

// IndexOf returns the chracter index under the
// point pt.
func (f *Frame) IndexOf(pt image.Point) int64 {
	pt = f.grid(pt)
	qt := f.r.Min
	p := int64(0)
	bn := 0
	for ; bn < f.Nbox && qt.Y < pt.Y; bn++ {
		b := &f.Box[bn]
		qt = f.wrapMax(qt, b)
		if qt.Y >= pt.Y {
			break
		}
		qt = f.advance(qt, b)
		p += int64(b.Len())
	}

	for ; bn < f.Nbox && qt.X <= pt.X; bn++ {
		b := &f.Box[bn]
		qt = f.wrapMax(qt, b)
		if qt.Y > pt.Y {
			break
		}
		if qt.X+b.Width > pt.X {
			if b.Nrune < 0 {
				qt = f.advance(qt, b)
			} else {
				step := f.Font.Fits(b.Ptr, qt.X-pt.X)
				if step == len(b.Ptr) {
					qt.X += b.Width
				} else {
					qt.X += f.Font.Dx(b.Ptr[:step])
				}
				p += int64(step)
			}
		} else {
			p += int64(b.Len())
			qt = f.advance(qt, b)
		}
	}
	return p
}
