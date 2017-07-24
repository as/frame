package frame

import (
	"image"
	"image/draw"
)

// Put
type Pts [2]image.Point

func (f *Frame) Mark() {
	f.modified = true
}

/*
func (f *Frame) Insert(s []byte, p0 int64) (wrote int64) {
	var (
		pt0, pt1, ppt0, ppt1, opt0 image.Point
		n, n0, nn0, y              int
		cn0                        int64
		back, text                 image.Image
	)
	h := f.Font.height
	if p0 > f.Nchars || len(s) == 0 || f.b == nil {
		return
	}

	// find p0, it's box, and its point in the box its in
	n0 = f.Find(0, 0, p0)
	cn0 = p0
	nn0 = n0
	pt0 = f.ptOfCharNBox(p0, n0)
	ppt0 = pt0
	opt0 = pt0

	// find p1
	ppt0, pt1 = f.bxscan(s, ppt0)
	ppt1 = pt1
	// Line wrap
	if n0 < f.Nbox {
		b := &f.Box[n0]
		pt0 = f.lineWrap(pt0, b)
		ppt1 = f.lineWrap0(ppt1, b)
	}
	f.modified = true

	// pt0, pt1   - start and end of insertion (current; and without line wrap)
	// ppt0, ppt1 - start and end of insertion when its complete

	// Multiple ticks
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), false)
	}

	// Find the points where all the old x and new x line up
	// Invariants:
	//   pt[0] is where the next box (b, n0) is now
	//   pt[1] is where it will be after insertion
	// If pt[1] goes out of bounds, we're done

	f.pts = f.pts[:0]
	for ; pt1.X != pt0.X && pt1.Y != f.r.Max.Y && n0 < f.Nbox; n0++ {
		b := &f.Box[n0]
		pt0 = f.lineWrap(pt0, b)
		pt1 = f.lineWrap0(pt1, b)
		if b.Nrune > 0 {
			if n = f.canFit(pt1, b); n != b.Nrune {
				f.Split(n0, n)
				b = &f.Box[n0]
			}
		}
		f.pts = append(f.pts, Pts{pt0, pt1})
		// check for text overflow off the frame
		if pt1.Y == f.r.Max.Y {
			break
		}
		pt0 = f.advance(pt0, b)
		pt1.X += f.newWid(pt1, b)
		cn0 += int64(b.Len())
	}
	if pt1.Y == f.r.Max.Y && n0 < f.Nbox {
		f.Nchars -= f.Len(n0)
		f.Run.Delete(n0, f.Nbox-1)
	}

	if n0 == f.Nbox {
		f.Nlines = (pt1.Y - f.r.Min.Y) / h
		if pt1.X > f.r.Min.X {
			f.Nlines++
		}
	} else if pt1.Y != pt0.Y {
		y = f.r.Max.Y
		qt0 := pt0.Y + h
		qt1 := pt1.Y + h
		f.Nlines += (qt1 - qt0) / h
		if f.Nlines > f.maxlines {
			f.chopFrame(ppt1, p0, nn0)
		}
		if pt1.Y < y {
			r := f.r
			r.Min.Y = qt1
			r.Max.Y = y
			if qt1 < y {
				f.Draw(f.b, r, f.b, image.Pt(f.r.Min.X, qt0), f.op)
			}
			r.Min = pt1
			r.Max.X = pt1.X + (f.r.Max.X - pt0.X)
			r.Max.Y = qt1
			f.Draw(f.b, r, f.b, pt0, f.op)
		}
	}

	// Move the old stuff down to make rooms
	if pt1.Y == f.r.Max.Y {
		y = pt1.Y
	} else {
		y = 0
	}
	npts := len(f.pts)
	npts--
	for ctr := n0; npts >= 0; npts-- {
		ctr--
		b := &f.Box[ctr]
		pt := f.pts[npts]
		dx := b.Width
		br := image.Rect(0, 0, dx, h)
		if b.Nrune > 0 {
			f.Draw(f.b, br.Add(pt[1]), f.b, pt[0], f.op)
			// clear bit hanging off right
			if npts == 0 && pt[1].Y > pt0.Y {
				_, back = f.pick(cn0, f.p0, f.p1)
				// line wrap - new char bigger than first char displaced
				r := image.Rectangle{opt0, opt0}
				r.Max.X=f.r.Max.X
				r.Max.Y+=h
				// back
				f.Draw(f.b, r, back, r.Min, f.op)
			} else if pt[1].Y < y {
				// copy from left to right
				_, back = f.pick(cn0, f.p0, f.p1)

				r := image.ZR.Add(pt[1])
				r.Min.X += dx
				//r.Max.X = f.r.Max.X
				r.Max.Y += h
				f.Draw(f.b, r, back, r.Min, f.op)
				
				next()
			}
			y = pt[1].Y
			cn0 -= int64(b.Nrune)
		} else {
			r := br.Add(pt[1])
			if r.Max.X >= f.r.Max.X {
				r.Max.X = f.r.Max.X
			}
			cn0--
			_, back = f.pick(cn0, f.p0, f.p1)
			f.Draw(f.b, r, back, r.Min, f.op)
			y = 0
			if pt[1].X == f.r.Min.X {
				y = pt[1].Y
			}
		}
	}
	// insertion extends the selection
	text, back = f.pick(p0, f.p0+1, f.p1+1)
	fr := *f.fr
	//f.Paint(ppt0, ppt1, image.NewUniform(color.RGBA{123, 33, 234, 255}))
	f.Paint(ppt0, ppt1, back)
	//back
	(&fr).Redraw0(ppt0, text, back)
	f.Add(nn0, fr.Nbox)
	for n = 0; n < fr.Nbox; n++ {
		f.Box[nn0+n] = fr.Box[n]
	}
	if nn0 > 0 && f.Box[nn0-1].Nrune >= 0 && ppt0.X-f.Box[nn0-1].Width >= f.r.Min.X {
		nn0--
		ppt0.X -= f.Box[nn0].Width
	}

	n0 += fr.Nbox
	if n0 < f.Nbox-1 {
		n0++
	}
	f.clean(ppt0, nn0, n0)
	wrote = fr.Nchars
	f.Nchars += fr.Nchars
	if f.p0 >= p0 {
		f.p0 += fr.Nchars
	}
	if f.p0 > f.Nchars {
		f.p0 = f.Nchars
	}
	if f.p1 >= p0 {
		f.p1 += fr.Nchars
	}
	if f.p1 > f.Nchars {
		f.p1 = f.Nchars
	}
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(f.p0), true)
	}
	return wrote
}
*/

func (f *Frame) Insert(s []byte, p0 int64) (wrote int64) {
	var (
		pt0, pt1, ppt0, ppt1, opt0 image.Point
		n, n0, nn0, y int
		cn0           int64
		back, text    image.Image
	)
	h := f.Font.height
	if p0 > f.Nchars || len(s) == 0 || f.b == nil {
		return
	}

	// find p0, it's box, and its point in the box its in
	n0 = f.Find(0, 0, p0)
	cn0 = p0
	nn0 = n0
	pt0 = f.ptOfCharNBox(p0, n0)
	ppt0 = pt0
	opt0 = pt0

	// find p1
	ppt0, pt1 = f.bxscan(s, ppt0)
	ppt1 = pt1
	// Line wrap
	if n0 < f.Nbox {
		b := &f.Box[n0]
		pt0 = f.lineWrap(pt0, b)
		ppt1 = f.lineWrap0(ppt1, b)
	}
	f.modified = true

	// pt0, pt1   - start and end of insertion (current; and without line wrap)
	// ppt0, ppt1 - start and end of insertion when its complete

	// Multiple ticks
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), false)
	}

	// Find the points where all the old x and new x line up
	// Invariants:
	//   pt[0] is where the next box (b, n0) is now
	//   pt[1] is where it will be after insertion
	// If pt[1] goes out of bounds, we're done

	f.pts = f.pts[:0]
	for ; pt1.X != pt0.X && pt1.Y != f.r.Max.Y && n0 < f.Nbox; n0++ {
		b := &f.Box[n0]
		pt0 = f.lineWrap(pt0, b)
		pt1 = f.lineWrap0(pt1, b)
		if b.Nrune > 0 {
			if n = f.canFit(pt1, b); n != b.Nrune {
				f.Split(n0, n)
				b = &f.Box[n0]
			}
		}
		// check for text overflow off the frame
		if pt1.Y == f.r.Max.Y {
			break
		}
		f.pts = append(f.pts, Pts{pt0, pt1})
		pt0 = f.advance(pt0, b)
		pt1.X += f.newWid(pt1, b)
		cn0 += int64(b.Len())
	}
	if pt1.Y == f.r.Max.Y && n0 < f.Nbox {
		f.Nchars -= f.Len(n0)
		f.Run.Delete(n0, f.Nbox-1)
	}

	if n0 == f.Nbox {
		f.Nlines = (pt1.Y - f.r.Min.Y) / h
		if pt1.X > f.r.Min.X {
			f.Nlines++
		}
	} else if pt1.Y != pt0.Y {
		y = f.r.Max.Y
		qt0 := pt0.Y + h
		qt1 := pt1.Y + h
		f.Nlines += (qt1 - qt0) / h
		if f.Nlines > f.maxlines {
			f.chopFrame(ppt1, p0, nn0)
		}
		if pt1.Y < y {
			r := f.r
			r.Min.Y = qt1
			r.Max.Y = y
			if qt1 < y {
				f.Draw(f.b, r, f.b, image.Pt(f.r.Min.X, qt0), f.op)
			}
			r.Min = pt1
			r.Max.X = pt1.X + (f.r.Max.X - pt0.X)
			r.Max.Y = qt1
			f.Draw(f.b, r, f.b, pt0, f.op)
		}
	}

	// Move the old stuff down to make rooms
	if pt1.Y == f.r.Max.Y {
		y = pt1.Y
	} else {
		y = 0
	}
	
	x := len(f.pts)
	run := f.Box[n0-x:]
	x--
	for ; x >= 0; x-- {
		b := &run[x] // &f.Box[ctr]
		dx := b.Width
		br := image.Rect(0,0,dx,h)
		pt := f.pts[x]
		if b.Nrune > 0 {
			f.Draw(f.b, br.Add(pt[1]), f.b, pt[0], f.op)
			// clear bit hanging off right
			if x == 0 && pt[1].Y > pt0.Y {
				_, back = f.pick(cn0, f.p0, f.p1)
				// line wrap - new char bigger than first char displaced
				r := br.Add(opt0)
				r.Max.X = f.r.Max.X
				f.Draw(f.b, r, back, r.Min, f.op)
			} else if pt[1].Y < y {
				// copy from left to right
				_, back = f.pick(cn0, f.p0, f.p1)

				r := image.ZR.Add(pt[1])
				r.Min.X += dx
				r.Max.X += f.r.Max.X
				r.Max.Y += h
				f.Draw(f.b, r, back, r.Min, f.op)
			}
			y = pt[1].Y
			cn0 -= int64(b.Nrune)
		} else {
			r := br.Add(pt[1])
			if r.Max.X >= f.r.Max.X {
				r.Max.X = f.r.Max.X
			}
			cn0--
			_, back = f.pick(cn0, f.p0, f.p1)
			f.Draw(f.b, r, back, r.Min, f.op)
			y = 0
			if pt[1].X == f.r.Min.X {
				y = pt[1].Y
			}
		}
	}
	// insertion extends the selection
	text, back = f.pick(p0, f.p0+1, f.p1+1)
	fr := *f.fr
	f.Paint(ppt0, ppt1, back)
	(&fr).Redraw0(ppt0, text, back)
	f.Add(nn0, fr.Nbox)
	for n = 0; n < fr.Nbox; n++ {
		f.Box[nn0+n] = fr.Box[n]
	}
	if nn0 > 0 && f.Box[nn0-1].Nrune >= 0 && ppt0.X-f.Box[nn0-1].Width >= f.r.Min.X {
		nn0--
		ppt0.X -= f.Box[nn0].Width
	}

	n0 += fr.Nbox
	if n0 < f.Nbox-1 {
		n0++
	}
	f.clean(ppt0, nn0, n0)
	wrote = fr.Nchars
	f.Nchars += fr.Nchars
	if f.p0 >= p0 {
		f.p0 += fr.Nchars
	}
	if f.p0 > f.Nchars {
		f.p0 = f.Nchars
	}
	if f.p1 >= p0 {
		f.p1 += fr.Nchars
	}
	if f.p1 > f.Nchars {
		f.p1 = f.Nchars
	}
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(f.p0), true)
	}
	return wrote
}

func (f *Frame) pick(c, p0, p1 int64) (text, back image.Image) {
	if p0 <= c && c < p1 {
		return f.Color.Hi.Text, f.Color.Hi.Back
	}
	return f.Color.Text, f.Color.Back
}

func region(c, p0, p1 int64) int {
	if c < p0 {
		return -1
	}
	if c >= p1 {
		return 1
	}
	return 0
}

func drawBorder(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, thick int) {
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+thick), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-thick, r.Max.X, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X+thick, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Max.X-thick, r.Min.Y, r.Max.X, r.Max.Y), src, sp, draw.Src)
}
