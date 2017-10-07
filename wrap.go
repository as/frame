package frame

import (
	"fmt"
	"github.com/as/frame/box"
	"image"
)

// LineWrap checks whether the box would wrap across a line boundary
// if it were inserted at pt. If it wraps, the line-wrapped point is
// returned.
func (f *Frame) lineWrap(pt image.Point, b *box.Box) image.Point {
	width := b.Width
	if b.Nrune < 0 {
		width = b.Minwidth
	}
	if width > f.r.Max.X-pt.X {
		pt.X = f.r.Min.X
		pt.Y += f.Font.Dy()
	}
	return pt
}

// LineWrap0 returns the line-wrapped point if the box doesn't
// fix on the line
func (f *Frame) lineWrap0(pt image.Point, b *box.Box) image.Point {
	if f.canFit(pt, b) == 0 {
		pt.X = f.r.Min.X
		pt.Y += f.Font.Dy()
	}
	return pt
}

// NewLineTrim returns the number of characters that would
// underflow on the left if b terminated at point pt.
func (f *Frame) canFitLeft(pt image.Point, b *box.Box) int {
	pt.X -= b.Width
	pt.X = f.r.Max.X - pt.X
	n := f.canFit(pt, b)
	return b.Len() - n
}

// canFit returns the number of runes that can fit
// on the line at pt. A newline yields 1.
func (f *Frame) canFit(pt image.Point, b *box.Box) (nr int) {
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
	for bx := f.newRulerFunc(b.Ptr, f.Font); ; {
		_, px, err := bx.Next()
		if err != nil {
			break
		}
		left -= px
		if left < 0 {
			return nr
		}
		nr++
	}
	// The box was too short and didn't end on a line boundary
	panic(fmt.Sprintf("CanFit: short box: len=%d left=%d box=%q\n", len(b.Ptr), left, b.Ptr))
}

func (f *Frame) advance(pt image.Point, b *box.Box) (x image.Point) {
	if b.Nrune < 0 && b.BC == '\n' {
		pt.X = f.r.Min.X
		pt.Y += f.Font.Dy()
	} else {
		pt.X += b.Width
	}
	return pt
}

// TODO: Naming
func (f *Frame) newWid(pt image.Point, b *box.Box) int {
	b.Width = f.newWid0(pt, b)
	return b.Width
}
func (f *Frame) newWid0(pt image.Point, b *box.Box) int {
	c := f.r.Max.X
	x := pt.X
	if b.Nrune >= 0 || b.BC != '\t' {
		return b.Width
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
