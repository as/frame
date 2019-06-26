package frame

import (
	"image"

	"github.com/as/frame/box"
	"golang.org/x/image/math/fixed"
)

type (
	p26 = fixed.Point26_6
	i26 = fixed.Int26_6
	r26 = fixed.Rectangle26_6
)

func int26(i int) i26 {
	return i26(i << 6)
}
func pt26(p image.Point) p26 {
	return p26{i26(p.X << 6), i26(p.Y << 6)}
}
func rect26(r image.Rectangle) r26 {
	return r26{pt26(r.Min), pt26(r.Max)}
}
func pt26point(p p26) image.Point {
	return image.Point{p.X.Ceil(), p.Y.Ceil()}
}
func r26rect(r r26) image.Rectangle {
	return image.Rectangle{
		pt26point(r.Min),
		pt26point(r.Max),
	}
}

// wrapMax returns the point where b should go on a plane
// if b doesn't fit entirely on the plane at pt, wrapMax
// returns a pt on the next line
func (f *Frame) wrapMax(pt p26, b *box.Box) p26 {
	width := b.Width
	if b.Nrune < 0 {
		width = b.Minwidth
	}
	if width > f.r.Max.X-pt.X {
		return f.wrap(pt)
	}
	return pt
}

// wrapMin is like wrapMax, except it lazily wraps lines if
// no chars in the box fit on the plane at pt.
func (f *Frame) wrapMin(pt p26, b *box.Box) p26 {
	if f.fits(pt, b) == 0 {
		return f.wrap(pt)
	}
	return pt
}

func (f *Frame) dy() i26 {
	h := f.Face.Metrics().Height
	return h + h/2
}
func (f *Frame) wrap(pt p26) p26 {
	pt.X = f.r.Min.X
	pt.Y += f.dy()
	return pt
}

func (f *Frame) advance(pt p26, b *box.Box) (x p26) {
	if b.Nrune < 0 && b.Break() == '\n' {
		pt = f.wrap(pt)
	} else {
		pt.X += b.Width
	}
	return pt
}

// fits returns the number of runes that can fit on the line at pt. A newline yields 1.
func (f *Frame) fits(pt p26, b *box.Box) (nr int) {
	left := f.r.Max.X - pt.X
	if b.Nrune < 0 {
		if b.Minwidth <= left {
			return 1
		}
		return 0
	}
	if left >= b.Width {
		return b.Nrune
	}
	return f.Drawer.MaxFit(b.Ptr, left)
}
func (f *Frame) plot(pt p26, b *box.Box) i26 {
	b.Width = f.project(pt, b)
	return b.Width
}
func (f *Frame) project(pt p26, b *box.Box) i26 {
	c := f.r.Max.X
	x := pt.X
	if b.Nrune >= 0 || b.Break() != '\t' { //
		return b.Width
	}
	if f.elastic() && b.Break() == '\t' {
		return b.Minwidth
	}
	if x+b.Minwidth > c {
		pt.X = f.r.Min.X
		x = pt.X
	}
	x += f.maxtab
	x -= (x - f.r.Min.X) % f.maxtab
	if x-pt.X < b.Minwidth || x > c {
		x = pt.X + b.Minwidth
	}
	return x - pt.X
}
