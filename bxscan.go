package frame

import (
	"image"
)

func (f *Frame) bxscan(s []byte, ppt image.Point) (image.Point, image.Point) {
	f.ir.Reset(f.ir.Measure)
	f.ir.Bxscan(s, f.maxlines)
	ppt = f.lineWrap0(ppt, &f.ir.Box[0])
	return ppt, f.drawRun(f.ir, ppt)
}
