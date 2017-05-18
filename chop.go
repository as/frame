package frame

import (
	"image"
	"fmt"
)

func (f *Frame) setlines(s string, i int){
	fmt.Printf("[%p] %s\t%d ->%d\n", f, s, f.Nlines, i)
	f.Nlines = i
}

func (f *Frame) ChopFrame(pt image.Point, p int64, bn int) {
	for ; bn < f.Nbox; bn++ {
		b := &f.Box[bn]
		pt = f.LineWrap(pt, b)
		if pt.Y >= f.r.Max.Y {
			break
		}
		p += int64(b.Len())
		pt = f.Advance(pt, b)
	}
	f.Nchars = p
	f.Nlines = f.maxlines
//	f.setlines("ChopFrame", f.maxlines)
	if bn < f.Nbox {
		f.Run.Delete(bn, f.Nbox-1)
	}
}
// Put
