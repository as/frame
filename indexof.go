package frame

import (
	"image"
)

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
				ptr := b.Ptr
				bsb := len(ptr)
				i := 0
				for {
					if bsb == i {
						break
					}
					i++
					size, width := 1, f.Font.Dx(ptr[:i])
					//i+=size
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
