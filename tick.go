package frame

import (
	"image"
	"image/draw"
)

const (
	TickWidth = 3
	tickOff   = 0
	tickOn    = 1
)

func (f *Frame) Untick() {
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), false)
	}
}
func (f *Frame) Tick() {
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), true)
	}
}

func (f *Frame) SetTick(style int) {
	f.tickoff = style == tickOff
}
func (f *Frame) inittick() {
	h := f.Font.Dy()
	r := image.Rect(0, 0, TickWidth, h).Inset(-1)
	f.tickscale = 1 // TODO implement scalesize
	f.tick = image.NewRGBA(r)
	f.tickback = image.NewRGBA(r)
	drawtick := func(x0, y0, x1, y1 int) {
		draw.Draw(f.tick, image.Rect(x0+1, y0+1, x1+1, y1+1), f.Color.Text, image.ZP, draw.Src)
	}
	drawtick(TickWidth/2, 0, TickWidth/2+1, h)
	drawtick(0, 0, TickWidth, h/5)
	drawtick(0, h-h/5, TickWidth, h)
}

// Put
func (f *Frame) tickat(pt image.Point, ticked bool) {
	if f.Ticked == ticked || f.tick == nil || !pt.In(f.Bounds().Inset(-1)) {
		return
	}
	//pt.X--
	r := f.tick.Bounds().Add(pt)
	if r.Max.X > f.r.Max.X {
		r.Max.X = f.r.Max.X
	} //
	adj := image.Pt(1, 1)
	if ticked {
		draw.Draw(f.tickback, f.tickback.Bounds(), f.b, pt.Sub(adj), draw.Src)
		f.Draw(f.b, r.Sub(adj), f.tick, image.ZP, draw.Over)
	} else {
		f.Draw(f.b, r, f.tickback, image.ZP.Sub(adj), draw.Src)
	}
	f.Ticked = ticked
}
