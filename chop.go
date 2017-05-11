package frame

import (
	//	"fmt"
	"image"
)

func (f *Frame) ChopFrame(pt image.Point, p int64, bn int) {
	//	fmt.Printf("nlines=%d maxlines=%d\n", f.Nlines, f.maxlines)
	//	fmt.Printf("chop at pt=%s p=%d bn=%d\n", pt, p, bn)
	for ; bn < f.Nbox; bn++ {
		b := &f.Box[bn]
		pt = f.LineWrap(pt, b)
		//fmt.Printf("%d >= %d\n", pt, f.r.Max.Y)
		if pt.Y >= f.r.Max.Y {
			break
		}
		p += int64(b.Len())
		pt = f.Advance(pt, b)
	}
	f.Nchars = p
	f.Nlines = f.maxlines
	if bn < f.Nbox {
		f.Run.Delete(bn, f.Nbox-1)
	}
}
