package win

import (
	"github.com/as/frame"
	"image"
	"image/draw"
)

const minSbWidth = 10

func (w *Win) scrollinit(pad image.Point) {
	w.Scrollr = image.ZR
	if pad.X > minSbWidth+2 {
		sr := w.Frame.RGBA().Bounds()
		sr.Max.X = minSbWidth
		w.Scrollr = sr
	}
	w.Frame.Draw(w.Frame.RGBA(), w.Scrollr, frame.ATag0.Back, image.ZP, draw.Src)
}

/*
 */
func (w *Win) Clicksb(pt image.Point, dir int) {
	var (
		rat float64
	)
	//	pt.Y -= w.pad.Y
	fl := float64(w.Frame.Len())
	n := w.org
	barY0 := float64(w.bar.Min.Y)
	barY1 := float64(w.bar.Max.Y)
	ptY := float64(pt.Y)
	switch dir {
	case -1:
		rat = barY1 / ptY
		delta := int64(fl * rat)
		n -= delta
	case 0:
		rat = (ptY - barY0) / (barY1 - barY0)
		delta := int64(fl * rat)
		n += delta
	case 1:
		rat = (barY1 / ptY)
		delta := int64(fl * rat)
		n += delta
	}
	w.SetOrigin(n, false)
	w.updatesb()
	w.dirty = true
}

func (w *Win) realsbr(r image.Rectangle) image.Rectangle {
	return r.Add(w.Sp).Add(image.Pt(0, w.pad.Y))
}

func (w *Win) drawsb() {
	w.Frame.Draw(w.Frame.RGBA(), w.Scrollr, frame.ATag0.Back, image.ZP, draw.Src)
	w.Frame.Draw(w.Frame.RGBA(), w.bar, LtGray, image.ZP, draw.Src)
}

func (w *Win) updatesb() {
	r := w.Scrollr
	dy := float64(w.Frame.Bounds().Dy() + w.pad.Y*2)
	rat0 := float64(w.org) / float64(w.Len()) // % scrolled
	r.Min.Y = +int(dy * rat0)

	rat1 := float64(w.org+w.Frame.Len()) / float64(w.Len()) // % covered by screen
	r.Max.Y = int(dy * rat1)
	if r.Max.Y-r.Min.Y < 3 {
		r.Max.Y = r.Min.Y + 3
	}
	w.dirty = true
	w.bar = r
	w.drawsb()
}
