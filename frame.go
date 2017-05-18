package frame

import (
	"github.com/as/frame/box"
	"image"
	"image/draw"
)

const (
	TEXT  = 0
	BACK  = 1
	HTEXT = 2
	HIGH  = 3

	TickWidth = 3
	TickOff = 0
	TickOn  = 1
)

type Frame struct {
	box.Run
	Color
	Font         Font
	disp, b      draw.Image
	r, entire    image.Rectangle
	maxtab       int
	lastlinefull int

	P0 int64
	P1 int64

	tick      draw.Image
	tickback  draw.Image
	ticked    bool
	tickscale int
	tickoff   bool
	maxlines  int
	modified  bool
	noredraw  bool
	op        draw.Op

	npts     int
	pts      []Pts
	Cache    []image.Rectangle

	Scroll func(int)
	fr *Frame
}

func New(r image.Rectangle, ft Font, b draw.Image, cols Color) *Frame {
	f := &Frame{
		Font:   ft,
		maxtab: 4 * ft.Measure(' '),
		Color:  cols,
		Run:    box.Run{Measure: ft.stringwidth},
		op:     draw.Src,
	}
	f.setrects(r, b)
	f.cacheinit()
	f.inittick()
	f.fr = new(Frame)
	return f
}

// Reset resets the frame to display on image b with bounds r and font ft.
func (f *Frame) Reset(r image.Rectangle, b draw.Image, ft Font) {
	f.r = r
	f.b = b
	f.Font = ft
	f.Run.Reset(f.Font.stringwidth)
}

func (f *Frame) Dx(s string) int {
	return f.Font.Dx(s)
}
func (f *Frame) Dy() int {
	return f.Font.Dy()
}
func (f *Frame) Bounds() image.Rectangle {
	return f.r.Bounds()
}

func (f *Frame) SetTick(style int) {
	f.tickoff = style == TickOff
}
func (f *Frame) SetOp(op draw.Op) {
	f.op = op
}
func (f *Frame) cacheinit() {
	f.Cache = make([]image.Rectangle, 0, 1024)
}
func (f *Frame) flushcache() {
	f.Cache = f.Cache[:0]
}

func (f *Frame) inittick() {
	h := f.Font.height + (f.Font.height / 5)
	r := image.Rect(0, 0, TickWidth, h)
	f.tickscale = 1 // TODO implement scalesize
	f.tick = image.NewRGBA(r)
	f.tickback = image.NewRGBA(r)
	//draw.Draw(f.tick, f.tick.Bounds(), f.cols[BACK], image.ZP, draw.Src)
	drawtick := func(x0, y0, x1, y1 int) {
		draw.Draw(f.tick, image.Rect(x0, y0, x1, y1), f.Color.Text, image.ZP, draw.Src)
	}
	drawtick(TickWidth/2, 0, TickWidth/2+1, h)
	drawtick(0, 0, TickWidth, h/5)
	drawtick(0, h-h/5, TickWidth, h)
}

func (f *Frame) setrects(r image.Rectangle, b draw.Image) {
	f.b = b
	f.entire = r
	f.r = r
	f.r.Max.Y -= f.r.Dy() % f.Font.height
	f.maxlines = f.r.Dy() / f.Font.height
}

func (f *Frame) Clear(freeall bool) {
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
	f.ticked = false
}

// Full returns true if the last line in the frame is full
func (f *Frame) Full() bool {
	return f.lastlinefull == 1
}

func (f *Frame) MaxLine() int {
	return f.maxlines
}
func (f *Frame) Line() int {
	return f.Nlines
}

func free(i interface{}) {
}
func freeimage(i image.Image) {
}
