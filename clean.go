package frame

import (
	"github.com/as/frame/box"
	"image"
)

func (f *Frame) Clean(pt image.Point, n0, n1 int) {
	var b *box.Box
	c := f.r.Max.X
	nb := n0
	for ; nb < n1-1; nb++ {
		b = &f.Box[nb]
		b1 := &f.Box[nb+1]
		pt = f.LineWrap(pt, b)
		for b.Nrune >= 0 && nb < n1-1 && b1.Nrune >= 0 && pt.X+b.Width+b1.Width < c {
			f.Merge(nb)
			n1--
			b = &f.Box[nb]
		}
		pt = f.Advance(pt, &f.Box[nb])
	}

	for ; nb < f.Nbox; nb++ {
		b = &f.Box[nb]
		pt = f.LineWrap(pt, b)
		pt = f.Advance(pt, &f.Box[nb])
	}
	f.lastlinefull = 0
	if pt.Y >= f.r.Max.Y {
		f.lastlinefull = 1
	}
}
