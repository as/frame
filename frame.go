// Package frame provides plan9-like editable text images on a raster display. This implementation
// preserves NUL bytes, and uses a set of replacement characters for undenderable text glyphs generated
// with a smaller sized font (hexidecimal or ascii representation).
//
// A frame's text is not addressable
package frame

import (
	"github.com/as/frame/box"
	"github.com/as/text"
	"image"
	"image/draw"
)

const (
	TickWidth = 3
	tickOff   = 0
	tickOn    = 1
)

func (f *Frame) RGBA() *image.RGBA {
	return f.b
}
func (f *Frame) Size() image.Point {
	r := f.RGBA().Bounds()
	return image.Pt(r.Dx(), r.Dy())
}

type Frame struct {
	box.Run
	Color
	Font         Font
	b            *image.RGBA
	r, entire    image.Rectangle
	maxtab       int
	lastlinefull int

	p0 int64
	p1 int64

	tick      draw.Image
	tickback  draw.Image
	Ticked    bool
	tickscale int
	tickoff   bool
	maxlines  int
	modified  bool
	noredraw  bool
	op        draw.Op

	text.Drawer

	//	npts int
	pts [][2]image.Point

	Scroll func(int)
	fr     *Frame

	hexFont *Font
	hex     []draw.Image
}

// New creates a new frame on b with bounds r. The image b is used
// as the frame's internal bitmap cache.
func New(r image.Rectangle, ft Font, b *image.RGBA, cols Color) *Frame {
	f := &Frame{
		Font:   ft,
		maxtab: 4 * ft.Measure(' '),
		Color:  cols,
		Run:    box.Run{Measure: ft.MeasureBytes},
		op:     draw.Src,
	}
	f.setrects(r, b)
	f.inittick()
	f.fr = new(Frame)
	f.renderHex()
	f.Drawer = text.NewCached()
	return f
}

// Dirty returns true if the contents of the frame have changes since the last redraw
func (f *Frame) Dirty() bool {
	return f.modified
}

// SetDirty alters the frame's internal state
func (f *Frame) SetDirty(dirty bool) {
	f.modified = dirty
}

// Reset resets the frame to display on image b with bounds r and font ft.
func (f *Frame) Reset(r image.Rectangle, b *image.RGBA, ft Font) {
	f.r = r
	f.b = b
	f.SetFont(ft)
}

func (f *Frame) SetFont(ft Font) {
	f.Font = ft
	f.Run.Reset(f.Font.MeasureBytes)
	f.Refresh()
}

// Dx returns the width of s in pixels
func (f *Frame) Dx(s string) int {
	return f.Font.Dx(s)
}

// Dy returns the height of a glyphs bounding box
func (f *Frame) Dy() int {
	return f.Font.Dy()
}

// Bounds returns the frame's clipping rectangle
func (f *Frame) Bounds() image.Rectangle {
	return f.r.Bounds()
}

func (f *Frame) SetTick(style int) {
	f.tickoff = style == tickOff
}

// Full returns true if the last line in the frame is full
func (f *Frame) Full() bool {
	return f.lastlinefull == 1
}

// Maxline returns the max number of wrapped lines fitting on the frame
func (f *Frame) MaxLine() int {
	return f.maxlines
}

// Line returns the number of wrapped lines currently in the frame
func (f *Frame) Line() int {
	return f.Nlines
}

// Dot returns the range of the selected text
func (f *Frame) Dot() (p0, p1 int64) {
	return f.p0, f.p1
}

// Select sets the range of the selected text
func (f *Frame) Select(p0, p1 int64) {
	f.modified = true
	f.p0, f.p1 = p0, p1
}

func (f *Frame) inittick() {
	h := f.Font.height + (f.Font.height / 5)
	r := image.Rect(0, 0, TickWidth, h)
	f.tickscale = 1 // TODO implement scalesize
	f.tick = image.NewRGBA(r)
	f.tickback = image.NewRGBA(r)
	drawtick := func(x0, y0, x1, y1 int) {
		draw.Draw(f.tick, image.Rect(x0, y0, x1, y1), f.Color.Text, image.ZP, draw.Src)
	}
	drawtick(TickWidth/2, 0, TickWidth/2+1, h)
	drawtick(0, 0, TickWidth, h/5)
	drawtick(0, h-h/5, TickWidth, h)
}

func (f *Frame) setrects(r image.Rectangle, b *image.RGBA) {
	f.b = b
	f.entire = r
	f.r = r
	f.r.Max.Y -= f.r.Dy() % f.Font.height
	f.maxlines = f.r.Dy() / f.Font.height
}

func (f *Frame) clear(freeall bool) {
	if f.Nbox != 0 {
		f.Run.Delete(0, f.Nbox-1)
	}
	if f.Box != nil {
		free(f.Box)
	}
	if freeall {
		// TODO: unnecessary
		freeimage(f.tick)
		freeimage(f.tickback)
		f.tick = nil
		f.tickback = nil
	}
	f.Box = nil
	f.Ticked = false
}

func free(i interface{}) {
}
func freeimage(i image.Image) {
}
