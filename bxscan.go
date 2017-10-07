package frame

import (
	"image"
)

// bxscan resets the measuring function and calls Bxscan in the embedded run
func (f *Frame) bxscan(s []byte, ppt image.Point) (image.Point, image.Point) {
	f.ir.Reset(f.Font)
	f.ir.Bxscan(s, f.maxlines)
	ppt = f.wrapMin(ppt, &f.ir.Box[0])
	return ppt, f.drawRun(f.ir, ppt)
}
