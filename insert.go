package frame

import (
	"fmt"
	"github.com/as/frame/box"
	"image"
	"image/draw"
)
// Put
type Pts struct {
	pt0, pt1 image.Point
}

var (
	pts    []Pts
	Nalloc int
)

func (f *Frame) Mark() {
	f.modified = true
}

func (f *Frame) Insert(s []byte, p0 int64) (wrote int64) {
	var (
		pt0, pt1,
		ppt0, ppt1,
		opt0,
		pt image.Point

		b             *box.Box
		n, n0, nn0, y int
		cn0           int64
		back, text    image.Image

		r image.Rectangle
	)

	if p0 > f.Nchars || len(s) == 0 || f.b == nil {
		return
	}

	// find p0, it's box, and its point in the box its in
	n0 = f.Find(0, 0, p0)
	cn0 = p0
	nn0 = n0
	pt0 = f.PtOfCharNBox(p0, n0)
	//	fmt.Printf("insert: f.PtOfCharNBox <- %d: %s\n", p0, pt0)
	ppt0 = pt0
	opt0 = pt0

	// find p1
	ppt0, pt1 = f.bxscan(s, ppt0)
	ppt1 = pt1
	//	fmt.Printf("insert: pt0, pt1: %s, %s\n", pt0, pt1)
	//	fmt.Printf("insert: ppt0, ppt1: %s, %s\n", ppt0, ppt1)
	// Line wrap
	if n0 < f.Nbox {
		b = &f.Box[n0]
		pt0 = f.LineWrap(pt0, b)
		ppt1 = f.LineWrap0(ppt1, b)
	}
	f.modified = true

	// pt0, pt1   - start and end of insertion (current; and without line wrap)
	// ppt0, ppt1 - start and end of insertion when its complete

	if f.P0 == f.P1 {
		f.tickat(f.PtOfChar(int64(f.P0)), false)
	}

	// Find the points where all the old x and new x line up
	// Invariants:
	//   pt0 is where the next box (b, n0) is now
	//   pt1 is where it will be after insertion	ChopF
	// If pt1 goes off the rect, toss everything from there on
	f.npts = 0
	if n0 < f.Nbox {
		b = &f.Box[n0]
	}
	for ; pt1.X != pt0.X && pt1.Y != f.r.Max.Y && n0 < f.Nbox; n0, f.npts = n0+1, f.npts+1 {
		b = &f.Box[n0]
		pt0 = f.LineWrap(pt0, b)
		pt1 = f.LineWrap0(pt1, b)

		if b.Nrune > 0 {
			n = f.CanFit(pt1, b)
			//				fmt.Println("can fit %d in box from pt %s\n", n, pt1)
			if n == 0 {
				panic("f. ==0")
			}
			if n != b.Nrune {
				f.Split(n0, n)
				b = &f.Box[n0]
			}
		}

		if f.npts == Nalloc {
			pts = append(pts, make([]Pts, DELTA)...)
			Nalloc += DELTA
			b = &f.Box[n0]
		}
		pts[f.npts].pt0 = pt0
		pts[f.npts].pt1 = pt1
		// check for text overflow off the frame
		if pt1.Y == f.r.Max.Y {
			break
		}
		pt0 = f.Advance(pt0, b)
		pt1.X += f.NewWid(pt1, b)
		cn0 += int64(b.Len())
	}

	if pt1.Y > f.r.Max.Y {
		panic(fmt.Sprintf("frame.Insert: pt1 too far: %s > %s\n", pt1, f.r.Max))
	}
	if pt1.Y == f.r.Max.Y && n0 < f.Nbox {
		f.Nchars -= f.Len(n0)
		f.Run.Delete(n0, f.Nbox-1)
	}
// Put
	h := f.Font.height
	if n0 == f.Nbox {
		extra := 0
		if pt1.X > f.r.Min.X {
			extra = 1
		}
		f.Nlines = (pt1.Y-f.r.Min.Y)/h + extra
// f.setlines("Insert: 1/2", (pt1.Y-f.r.Min.Y)/h + extra)
// lines = (pt.Y-f.r.Min.Y)/h + extra
	} else if pt1.Y != pt0.Y {
		y = f.r.Max.Y
		q0 := pt0.Y + h
		q1 := pt1.Y + h
// fmt.Printf("pt1(%s) != pt0(%s)\n", pt1, pt0)
// fmt.Printf("y=%d q0:q1=%d:%d\n", y,q0,q1)
		f.Nlines += (q1 - q0) / h
// f.setlines("Insert: 2/2", f.Nlines+((q1 - q0) / h))
// lines += (q1 - q0) / h
		if f.Nlines > f.maxlines {
// fmt.Printf("lines/max = %d/%d\n", f.Nlines, f.maxlines)
			f.ChopFrame(ppt1, p0, nn0)
		}
		if pt1.Y < y {
			r = f.r
			r.Min.Y = q1
			r.Max.Y = y
			if q1 < y {
				draw.Draw(f.b, r, f.b, image.Pt(f.r.Min.X, q0), draw.Src)
			}
			r.Min = pt1
			r.Max.X = pt1.X + (f.r.Max.X - pt0.X)
			r.Max.Y = q1
			draw.Draw(f.b, r, f.b, pt0, draw.Src)
		}
	}

	// Move the old stuff down to make rooms
	if pt1.Y == f.r.Max.Y {
		y = pt1.Y
	} else {
		y = 0
	}

	f.npts--
	for ctr := n0; f.npts >= 0; f.npts-- {
		ctr--
		b = &f.Box[ctr]
		pt = pts[f.npts].pt1
//		fmt.Printf("npts=%d selected point = %s\n", npts, pt)
		if b.Nrune > 0 {
			r.Min = pt
			r.Max = r.Min
			r.Max.X += b.Width
			r.Max.Y += f.Font.height
			draw.Draw(f.b, r, f.b, pts[f.npts].pt0, draw.Src)

			// clear bit hanging off right
			if f.npts == 0 && pt.Y > pt0.Y {
				// first new char bigger than first char displaced
				// so line wrap happened
				r.Min = opt0
				r.Max = opt0
				r.Max.X = f.r.Max.X
				r.Max.Y += f.Font.height
				if f.P0 <= cn0 && cn0 < f.P1 { // b+1 is in selection
					back = f.Color.Hi.Back
				} else {
					back = f.Color.Back
				}
				draw.Draw(f.b, r, back, r.Min, draw.Src)
			} else if pt.Y < y {
				r.Min = pt
				r.Max = pt
				r.Min.X += b.Width
				r.Max.Y += f.Font.height
				if f.P0 <= cn0 && cn0 < f.P1 {
					back = f.Color.Hi.Back
				} else {
					back = f.Color.Back
				}
				draw.Draw(f.b, r, back, r.Min, draw.Src)
			}
			y = pt.Y
			cn0 -= int64(b.Nrune)
		} else {
			r.Min = pt
			r.Max = pt
			r.Max.X += b.Width
			r.Max.Y += f.Font.height
			if r.Max.X >= f.r.Max.X {
				r.Max.X = f.r.Max.X
			}
			cn0--
			if f.P0 <= cn0 && cn0 < f.P1 { // box inside selection
				back = f.Color.Hi.Back
			} else {
				back = f.Color.Back
			}
			draw.Draw(f.b, r, back, r.Min, draw.Src)
			y = 0
			if pt.X == f.r.Min.X {
				y = pt.Y
			}
		}
	}

	// insertion can extend the selection; different condition
	if f.P0 < p0 && p0 <= f.P1 {
		text = f.Color.Hi.Back
		back = f.Color.Hi.Text
	} else {
		text = f.Color.Text
		back = f.Color.Back
	}
	fr := *f.fr

	f.SelectPaint(ppt0, ppt1, back)
	(&fr).DrawText(ppt0, text, back)
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
		f.Clean(ppt0, nn0, n0+1)
	} else {
		f.Clean(ppt0, nn0, n0)
	}
	wrote = fr.Nchars
	f.Nchars += fr.Nchars
	if f.P0 >= p0 {
		f.P0 += fr.Nchars
	}
	if f.P0 > f.Nchars {
		f.P0 = f.Nchars
	}
	if f.P1 >= p0 {
		f.P1 += fr.Nchars
	}
	if f.P1 > f.Nchars {
		f.P1 = f.Nchars
	}
	if f.P0 == f.P1 {
		f.tickat(f.PtOfChar(f.P0), true)
	}
	return wrote
}
