package frame

import (
	"image"
	"image/draw"
)

func (f *Frame) Mark() {
	f.modified = true
}

// boxalign collects a list of pts of each box
// on the frame before and after an insertion occurs
func (f *Frame) boxalign(cb0 int64, b0 int, pt0, pt1 image.Point) (int64, int, image.Point, image.Point) {
	type Pts [2]image.Point
	f.pts = f.pts[:0]

	// collect the start pts for each box on the plane
	// before and after the insertion
	for {
		if pt0.X == pt1.X || pt1.Y == f.r.Max.Y || b0 == f.Nbox {
			break
		}
		b := &f.Box[b0]
		pt0 = f.wrapMax(pt0, b)
		pt1 = f.wrapMin(pt1, b)
		if b.Nrune > 0 {
			if n := f.fits(pt1, b); n != b.Nrune {
				f.Split(b0, n)
				b = &f.Box[b0]
			}
		}

		// early exit - point went off the frame
		if pt1.Y == f.r.Max.Y {
			break
		}

		f.pts = append(f.pts, Pts{pt0, pt1})

		pt0 = f.advance(pt0, b)
		pt1.X += f.plot(pt1, b)

		cb0 += int64(b.Len())
		b0++
	}
	return cb0, b0, pt0, pt1
}

// boxpush moves boxes down the frame to make room for an insertion
// from pt0 to pt1
func (f *Frame) boxpush(p0 int64, b0, b1 int, pt0, pt1, ppt1 image.Point) {

	// delete boxes that ran off the frame
	// and update the char count
	if pt1.Y == f.r.Max.Y && b0 < f.Nbox {
		f.Nchars -= f.Count(b0)
		f.Run.Delete(b0, f.Nbox-1)
	}

	// update the line count
	if b0 == f.Nbox {
		f.Nlines = (pt1.Y - f.r.Min.Y) / f.Font.Dy()
		if pt1.X > f.r.Min.X {
			f.Nlines++
		}
		return
	}

	if pt1.Y == pt0.Y {
		// insertion won't propagate down
		return
	}

	qt0 := pt0.Y + f.Font.Dy()
	qt1 := pt1.Y + f.Font.Dy()
	f.Nlines += (qt1 - qt0) / f.Font.Dy()
	if f.Nlines > f.maxlines {
		f.trim(ppt1, p0, b1)
	}

	// shift down the existing boxes
	// on the bitmap
	if r := f.r; pt1.Y < r.Max.Y {
		r.Min.Y = qt1

		// rectangular group of boxes
		if qt1 < f.r.Max.Y {
			f.Draw(f.b, r, f.b, image.Pt(f.r.Min.X, qt0), f.op)
		}

		// partial line
		r.Min = pt1
		r.Max.X = pt1.X + (f.r.Max.X - pt0.X)
		r.Max.Y = qt1
		f.Draw(f.b, r, f.b, pt0, f.op)
	}
}

func (f *Frame) boxmove(cb0 int64, b0 int, pt0, pt1, opt0 image.Point) {
	h := f.Font.Dy()
	y := 0
	if pt1.Y == f.r.Max.Y {
		y = pt1.Y
	}
	x := len(f.pts)
	run := f.Box[b0-x:]
	x--
	_, back := f.pick(cb0, f.p0, f.p1)
	for ; x >= 0; x-- {
		b := &run[x]
		br := image.Rect(0, 0, b.Width, h)
		pt := f.pts[x]
		if b.Nrune > 0 {
			f.Draw(f.b, br.Add(pt[1]), f.b, pt[0], f.op)
			// clear bit hanging off right
			if x == 0 && pt[1].Y > pt0.Y {
				_, back = f.pick(cb0, f.p0, f.p1)
				// line wrap - new char bigger than first char displaced
				r := br.Add(opt0)
				r.Max.X = f.r.Max.X
				f.Draw(f.b, r, back, r.Min, f.op)
			} else if pt[1].Y < y {
				// copy from left to right
				_, back = f.pick(cb0, f.p0, f.p1)

				r := image.ZR.Add(pt[1])
				r.Min.X += b.Width
				r.Max.X += f.r.Max.X
				r.Max.Y += h
				f.Draw(f.b, r, back, r.Min, f.op)
			}
			y = pt[1].Y
			cb0 -= int64(b.Nrune)
		} else {
			r := br.Add(pt[1])
			if r.Max.X >= f.r.Max.X {
				r.Max.X = f.r.Max.X
			}
			cb0--
			_, back = f.pick(cb0, f.p0, f.p1)
			f.Draw(f.b, r, back, r.Min, f.op)
			y = 0
			if pt[1].X == f.r.Min.X {
				y = pt[1].Y
			}
		}
	}
}

type offset struct {
	p0   int64
	b0   int
	cb0  int64
	b1   int64
	pt0  image.Point
	pt1  image.Point
	opt0 image.Point
}

func (f *Frame) Insert(s []byte, p0 int64) (wrote int) {
	type Pts [2]image.Point
	if p0 > f.Nchars || len(s) == 0 || f.b == nil {
		return
	}

	// find p0, it's box, and its point in the box its in
	b0 := f.Find(0, 0, p0)
	//ob0 := b0
	cb0 := p0
	b1 := b0
//	ebn := f.StartCell(b0)
//	ept := f.pointOfBox(p0, ebn)
	pt0 := f.pointOfBox(p0, b0)
	opt0 := pt0

	// find p1
	ppt0, pt1 := f.boxscan(s, pt0)
	ppt1 := pt1
	
	if ForceElasticTabstopExperiment {
		// With elastic tabstops, the scanned boxes need to be
		// connected to the frame's run to determine whether
		// anything above the scan point was affected by the
		// insertion
		
//		f.ir.Box = append(f.Box[ebn:b0], f.ir.Box)
//		f.ir.Nbox += len(f.Box[ebn:b0])
		
		bn := f.Nbox
		for bn > 0 {
			bn = f.Stretch(bn)
		}
		f.Stretch(bn)
	}
	
	// Line wrap
	if b0 < f.Nbox {
		b := &f.Box[b0]
		pt0 = f.wrapMax(pt0, b)
		ppt1 = f.wrapMin(ppt1, b)
	}
	f.modified = true


	if f.p0 == f.p1 {
		f.tickat(f.PointOf(int64(f.p0)), false)
	}
	

	cb0, b0, pt0, pt1 = f.boxalign(cb0, b0, pt0, pt1)
	f.boxpush(p0, b0, b1, pt0, pt1, ppt1)
	f.boxmove(cb0, b0, pt0, pt1, opt0)
	text, back := f.pick(p0, f.p0+1, f.p1+1)
	f.Paint(ppt0, ppt1, back)
	f.redrawRun0(f.ir, ppt0, text, back)
	
	f.Run.Combine(f.ir, b1)
//	f.Add(b1, f.ir.Nbox)
//	copy(f.Box[b1:], f.ir.Box[:f.ir.Nbox])
	if b1 > 0 && f.Box[b1-1].Nrune >= 0 && ppt0.X-f.Box[b1-1].Width >= f.r.Min.X {
		b1--
		ppt0.X -= f.Box[b1].Width
	}
	b0 += f.ir.Nbox
	if b0 < f.Nbox-1 {
		b0++
	}
	f.clean(ppt0, b1, b0)
	f.Nchars += f.ir.Nchars
	
	
	f.p0, f.p1 = coInsert(p0, p0+f.Nchars, f.p0, f.p1)
	if f.p0 == f.p1 {
		f.tickat(f.PointOf(f.p0), true)
	}
	if ForceElasticTabstopExperiment {
		// Just to see if the algorithm works not ideal to sift through all of
		// the boxes per insertion, although surprisingly faster than expected
		// to the point of where its almost unnoticable without the print
		// statements
		bn := f.Nbox
		for bn > 0 {
			bn = f.Stretch(bn)
		}
		f.Stretch(bn)
		f.Refresh() // must do this until line mapper is fixed
	}
	return int(f.ir.Nchars)
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
