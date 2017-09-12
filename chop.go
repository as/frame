package frame

import (
	"image"
)

// TODO(as): rename this to something more appropriate
func (f *Frame) chopFrame(pt image.Point, p int64, bn int) {
	for ; bn < f.Nbox; bn++ {
		b := &f.Box[bn]
		if pt = f.lineWrap(pt, b); pt.Y >= f.r.Max.Y {
			break
		}
		p += int64(b.Len())
		pt = f.advance(pt, b)
	}
	f.Nchars = p
	f.Nlines = f.maxlines
	if bn < f.Nbox {
		f.Run.Delete(bn, f.Nbox-1)
	}
}

