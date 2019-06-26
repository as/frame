package frame

import (
	"image"
)

const enablePointProximity = false

// IndexOf returns the chracter index under the
// point pt.
func (f *Frame) IndexOf(pt image.Point) (p int64) {
	pt.X += 1
	return f.indexOf(pt26(pt))
}

func (f *Frame) indexOf(pt p26) (p int64) {
	pt = f.grid(pt)
	qt := f.r.Min
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
				left := pt.X - qt.X
				p += int64(f.Face.Fits(b.Ptr, left.Ceil()))
				qt.X += left
			}
		} else {
			p += int64(b.Len())
			qt = f.advance(qt, b)
		}
	}

	if enablePointProximity {
		// NOTE(as): Slow, but nice. If the cursor 2/3s of the way
		// out of the glyph, assume we wanted the next glyph
		// instead.
		//
		// This should probably be in its own function, and not
		// call point twice (or even once).
		pt0 := f.point0(p)
		qt0 := f.point0(p + 1)
		if qt0.Y == pt0.Y {
			dx := pt.X - pt0.X
			dx = 3 * dx / 2
			if pt0.X+dx > qt0.X {
				p++
			}
		}
	}

	return p
}
