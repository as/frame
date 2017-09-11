package frame

import (
	"image"
)

// Put

func (f *Frame) Delete(p0, p1 int64) int {
	var (
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
	n1 = f.Find(n0, p0, p1)
	pt0 = f.ptOfCharNBox(p0, n0)
	pt1 = f.PointOf(p1)
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), false)
	}
	h := f.Font.Dy()
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

	//	b = &f.Box[n1]
	cn1 = int64(p1)

	for pt1.X != pt0.X && n1 < f.Nbox {
		b := &f.Box[n1]
		pt0 = f.lineWrap0(pt0, b)
		pt1 = f.lineWrap(pt1, b)
		r.Min = pt0
		r.Max = pt0
		r.Max.Y += h

		if b.Nrune > 0 { // non-newline
			n = f.canFit(pt0, b)
			if n != b.Nrune {
				f.Split(n1, n)
				b = &f.Box[n1]
			}
			r.Max.X += b.Width
			f.Draw(f.b, r, f.b, pt1, f.op)
			//drawBorder(f.b, r.Inset(-4), Green, image.ZP, 8)
			cn1 += int64(b.Nrune)
		} else {
			r.Max.X = min(r.Max.X+f.newWid0(pt0, b), f.r.Max.X)
			_, col = f.pick(cn1, f.p0, f.p1)
			f.Draw(f.b, r, col, pt0, f.op)
			cn1++
		}
		pt1 = f.advance(pt1, b)
		pt0.X += f.newWid(pt0, b)
		f.Box[n0] = f.Box[n1]
		n0++
		n1++
	}

	if n1 == f.Nbox && pt0.X != pt1.X {
		f.Paint(pt0, pt1, f.Color.Back)
	}

	if pt1.Y != pt0.Y {
		pt2 := f.ptOfCharPtBox(32768, pt1, n1)
		if pt2.Y > f.r.Max.Y {
			pt2.Y = f.r.Max.Y - h
			//panic(fmt.Sprintf("delete: PtOfCharPtBox %s > %s", pt2, f.r.Max))
		}
		if n1 < f.Nbox {
			q0 := pt0.Y + h
			q1 := pt1.Y + h
			q2 := pt2.Y + h
			if q2 > f.r.Max.Y {
				q2 = f.r.Max.Y
			}
			f.Draw(f.b, image.Rect(pt0.X, pt0.Y, pt0.X+(f.r.Max.X-pt1.X), q0), f.b, pt1, f.op)
			f.Draw(f.b, image.Rect(f.r.Min.X, q0, f.r.Max.X, q0+(q2-q1)), f.b, image.Pt(f.r.Min.X, q1), f.op)
			f.Paint(image.Pt(pt2.X, pt2.Y-(pt1.Y-pt0.Y)), pt2, f.Color.Back)
		} else {
			f.Paint(pt0, pt2, f.Color.Back)
		}
	}

	f.Close(n0, n1-1)
	if nn0 > 0 && f.Box[nn0-1].Nrune >= 0 && ppt0.X-f.Box[nn0-1].Width >= f.r.Min.X {
		nn0--
		ppt0.X -= f.Box[nn0].Width
	}

	if n0 < f.Nbox-1 {
		n0++
	}
	f.clean(ppt0, nn0, n0)

	if f.p1 > p1 {
		f.p1 -= p1 - p0
	} else if f.p1 > p0 {
		f.p1 = p0
	}
	if f.p0 > p1 {
		f.p0 -= p1 - p0
	} else if f.p0 > p0 {
		f.p0 = p0
	}

	f.Nchars -= p1 - p0
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(f.p0), true)
	}
	pt0 = f.PointOf(f.Nchars)
	n = f.Nlines
	extra := 0
	if pt0.X > f.r.Min.X {
		extra = 1
	}
	f.Nlines = (pt0.Y-f.r.Min.Y)/h + extra
	return int(p1 - p0 + 1) //n - f.Nlines
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a < b {
		return a
	}
	return b
}
