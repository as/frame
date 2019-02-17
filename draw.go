package frame

import (
	"image"
	"image/draw"

	"github.com/as/frame/box"
)

func (f *Frame) draw26(dst draw.Image, r r26, src image.Image, sp p26, op draw.Op) {
	f.Draw(dst, r26rect(r), src, pt26point(sp), op)
}

func (f *Frame) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, op draw.Op) {
	if f == nil {
		panic("nil frame")
	}
	if f.drawer == nil {
		f.drawer = newDrawer(dst, src, f.Face, pt26(r.Min))
	}
	f.drawer.Draw(dst, r, src, sp, op)
}

// Paint paints the color col on the frame at points pt0-pt1. The result is a Z shaped fill
// consisting of at-most 3 rectangles. No text is redrawn.
func (f *Frame) paint26(p0, p1 p26, col image.Image) {
	if f.b == nil {
		panic("selectpaint: b == 0")
	}
	if f.r.Max.Y == p0.Y {
		return
	}
	h := f.dy()
	q0, q1 := p0, p1
	q0.Y+= h
	q1.Y+= h
	n := (p1.Y - p0.Y) / h

	if n == 0 { // one line
		f.draw26(f.b, r26{p0, q1}, col, p26{}, draw.Over)
	} else {
		if p0.X >= f.r.Max.X {
			p0.X = f.r.Max.X // - 1
		}
		f.draw26(f.b, r26{p0, p26{f.r.Max.X, q0.Y}}, col, p26{}, draw.Over)
		if n > 1 {
			f.draw26(f.b, r26{p26{f.r.Min.X, q0.Y}, p26{f.r.Max.X, p1.Y}}, col, p26{}, draw.Over)
		}
		f.draw26(f.b, r26{p26{f.r.Min.X, p1.Y}, q1}, col, p26{}, draw.Over)
	}
}
func (f *Frame) Paint(p0, p1 image.Point, col image.Image) {
	f.paint26(pt26(p0), pt26(p0), col)
}

// Refresh renders the entire frame, including the underlying
// bitmap. Refresh should not be called after insertion and deletion
// unless the frame's RGBA bitmap was painted over by another
// draw operation.
func (f *Frame) Refresh() {
	cols := f.Color
	if f.p0 == f.p1 {
		ticked := f.Ticked
		if ticked {
			f.tickat(f.point0(f.p0), false)
		}
		f.drawsel(f.point0(0), 0, f.Nchars, cols.Back, cols.Text)
		if ticked {
			f.tickat(f.point0(f.p0), true)
		}
		return
	}
	pt := f.point0(0)
	pt = f.drawsel(pt, 0, f.p0, cols.Back, cols.Text)
	pt = f.drawsel(pt, f.p0, f.p1, cols.Hi.Back, cols.Hi.Text)
	f.drawsel(pt, f.p1, f.Nchars, cols.Back, cols.Text)
}

// RedrawAt renders the frame's bitmap starting at pt and working downwards.
func (f *Frame) RedrawAt(pt image.Point, text, back image.Image) {
	f.redrawRun0(&(f.Run), pt26(pt), text, back)
}

// Redraw draws the range [p0:p1] at the given pt.
func (f *Frame) Redraw(pt image.Point, p0, p1 int64, issel bool) {
	f.redraw(pt26(pt), p0,p1, issel)
}

func (f *Frame) redraw(pt p26, p0, p1 int64, issel bool) {
	if f.Ticked {
		f.tickat(f.point0(f.p0), false)
	}
	if p0 == p1 {
		f.tickat(pt, issel)
		return
	}
	pal := f.Color.Palette
	if issel {
		pal = f.Color.Hi
	}
	f.drawsel(pt, p0, p1, pal.Back, pal.Text)
}


// Recolor redraws the range p0:p1 with the given palette
func (f *Frame) Recolor(pt image.Point, p0, p1 int64, cols Palette) {
	f.drawsel(pt26(pt), p0, p1, cols.Back, cols.Text)
	f.modified = true
}

// widthBox returns the width of box n. If the length of
// alt is different than the box, alt is measured and
// returned instead.
func (f *Frame) widthBox(b *box.Box, alt []byte) i26 {
	if b.Nrune < 0 || len(alt) == b.Len() {
		return b.Width
	}
	return f.Dx(alt)
}

func (f *Frame) redrawRun0(r *box.Run, pt p26, text, back image.Image) p26 {
	// d := newDrawer(f.b, text, f.Face, pt)
	f.drawer.Reset(f.b, pt, text)
	nb := 0
	for ; nb < r.Nbox; nb++ {
		b := &r.Box[nb]
		pt = f.wrapMax(pt, b)
		if b.Nrune >= 0 {
			f.drawer.move(pt)
			f.drawer.Write(b.Ptr)
		}
		pt.X += b.Width
	}
	return pt
}

func (f *Frame) drawsel(pt p26, p0, p1 int64, back, text image.Image) p26 {
	f.drawer.Reset(f.b, pt, text)
	maxx := f.r.Max.X
	{
		// doubled
		p0, p1 := int(p0), int(p1)
		q0 := 0
		trim := false

		var (
			nb int
			b  *box.Box
		)

		for nb = 0; nb < f.Nbox; nb++ {
			b = &f.Box[nb]
			l := q0 + b.Len()
			if l > p0 {
				break
			}
			q0 = l
		}

		// Step into box, start coloring it
		// How much does this lambda slow things down?
		stepFill := func() {
			pt0 := pt
			pt = f.wrapMax(pt, b)
			if pt.Y > pt0.Y {
				f.drawer.drawBG(back, maxx)
			}
			f.drawer.move(pt)
		}
		for ; nb < f.Nbox && q0 < p1; nb++ {
			b = &f.Box[nb]
			if q0 >= p0 { // region 0 or 1 or 2
				stepFill()
			}
			ptr := b.Ptr[:b.Len()]
			if q0 < p0 {
				// region -1: shift p right inside the selection
				ptr = ptr[p0-q0:]
				q0 = p0
			}

			trim = false
			if q1 := q0 + len(ptr); q1 > p1 {
				// region 1: would draw too much, retract the selection
				lim := len(ptr) - (q1 - p1)
				ptr = ptr[:lim]
				trim = true
			}
			w := f.widthBox(b, ptr)
			f.drawer.drawBG(back, min(pt.X+w, maxx))
			if b.Nrune > 0 {
				f.drawer.Write(ptr)
			}
			pt.X += w
			f.drawer.move(pt)
			if q0 += len(ptr); q0 > p1 {
				break
			}
		}
		if p1 > p0 && nb != 0 && nb < f.Nbox && f.Box[nb-1].Len() > 0 && !trim {
			b = &f.Box[nb]
			stepFill()
		}
		return pt
	}
}
