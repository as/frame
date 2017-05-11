package frame

import (
	"github.com/as/frame/box"
	"image"
)

func (f *Frame) Delete(p0, p1 int64) int {
	var (
		b              *box.Box
		pt0, pt1, ppt0 image.Point
		n0, n1, n      int
		cn1            int64
		r              image.Rectangle
		nn0            int
		col            image.Image
	)

	if p0 >= f.Nchars || p0 == p1 || f.b == nil {
		return 0
	}

	if p1 > f.Nchars {
		p1 = f.Nchars
	}
	n0 = f.Find(0, 0, p0)
	if n0 == f.Nbox {
		panic("delete")
	}
	n1 = f.Find(n0, p0, p1)
	pt0 = f.PtOfCharNBox(p0, n0)
	pt1 = f.PtOfChar(p1)
	if f.P0 == f.P1 {
		f.tickat(f.PtOfChar(int64(f.P0)), false)
	}

	nn0 = n0
	ppt0 = pt0
	f.Free(n0, n1-1)
	f.modified = true

	// pt0, pt1 - beginning, end
	// n0 - has beginning of deletion
	// n1, b - first box kept after deletion
	// cn1 char pos of n1
	//
	// adjust f.p0 and f.p1 after deletion is finished

	if n1 > f.Nbox {
		panic("DeleteBytes: Split bug: nul terminators removed")
	}

	b = &f.Box[n1]
	cn1 = int64(p1)

	for pt1.X != pt0.X && n1 < f.Nbox {
		pt0 = f.LineWrap0(pt0, b)
		pt1 = f.LineWrap(pt1, b)
		r.Min = pt0
		r.Max = pt0
		r.Max.Y += f.Font.height

		if b.Nrune > 0 { // non-newline
			n = f.CanFit(pt0, b)
			if n == 0 {
				panic("delete: canfit==0")
			}
			if n != b.Nrune {
				f.Split(n1, n)
				b = &f.Box[n1]
			}
			r.Max.X += b.Width
			f.draw(f.b, r, f.b, pt1)
			//drawBorder(f.b, r.Add(pt1).Inset(-4), Red, image.ZP, 8)
			//drawBorder(f.b, r.Inset(-4), Green, image.ZP, 8)
			cn1 += int64(b.Nrune)
		} else {
			r.Max.X += f.NewWid0(pt0, b)
			if r.Max.X > f.r.Max.X {
				r.Max.X = f.r.Max.X
			}
			col = f.Color.Back
			if f.P0 <= cn1 && cn1 < f.P1 {
				col = f.Color.Hi.Back
			}
			f.draw(f.b, r, col, pt0)
			cn1++
		}
		pt1 = f.Advance(pt1, b)
		pt0.X += f.NewWid(pt0, b)
		f.Box[n0] = f.Box[n1]
		n0++
		n1++
		b = &f.Box[n1]
	}

	if n1 == f.Nbox && pt0.X != pt1.X {
		f.SelectPaint(pt0, pt1, f.Color.Back)
	}

	if pt1.Y != pt0.Y {
		pt2 := f.PtOfCharPtBox(32767, pt1, n1)
		if pt2.Y > f.r.Max.Y {
			//panic("delete: PtOfCharPtBox")
		}
		if n1 < f.Nbox {
			h := f.Font.height
			q0 := pt0.Y + h
			q1 := pt1.Y + h
			q2 := pt2.Y + h
			if q2 > f.r.Max.Y {
				q2 = f.r.Max.Y
			}
			f.draw(f.b, image.Rect(pt0.X, pt0.Y, pt0.X+(f.r.Max.X-pt1.X), q0), f.b, pt1)
			f.draw(f.b, image.Rect(f.r.Min.X, q0, f.r.Max.X, q0+(q2-q1)), f.b, image.Pt(f.r.Min.X, q1))
			f.SelectPaint(image.Pt(pt2.X, pt2.Y-(pt1.Y-pt0.Y)), pt2, f.Color.Back)
		} else {
			f.SelectPaint(pt0, pt2, f.Color.Back)
		}
	}

	f.Close(n0, n1-1)
	if nn0 > 0 && f.Box[nn0-1].Nrune >= 0 && ppt0.X-f.Box[nn0-1].Width >= f.r.Min.X {
		nn0--
		ppt0.X -= f.Box[nn0].Width
	}

	if n0 < f.Nbox-1 {
		f.Clean(ppt0, nn0, n0+1)
	} else {
		f.Clean(ppt0, nn0, n0)
	}

	if f.P1 > p1 {
		f.P1 -= p1 - p0
	} else if f.P1 > p0 {
		f.P1 = p0
	}

	if f.P0 > p1 {
		f.P0 -= p1 - p0
	} else if f.P0 > p0 {
		f.P0 = p0
	}

	f.Nchars -= p1 - p0
	if f.P0 == f.P1 {
		f.tickat(f.PtOfChar(f.P0), true)
	}
	pt0 = f.PtOfChar(f.Nchars)
	n = f.Nlines
	extra := 0
	if pt0.X > f.r.Min.X {
		extra = 1
	}
	extra = 1 // todo
	f.Nlines = (pt0.Y-f.r.Min.Y)/f.Font.height + extra
	return n - f.Nlines
}
