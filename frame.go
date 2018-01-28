package frame

import (
	"errors"
	. "github.com/as/font"
	"github.com/as/frame/box"
	"golang.org/x/image/font"
	"image"
	"image/draw"
)

var (
	ForceElastic bool
	ForceUTF8    bool
)

const (
	FrElastic = 1 << iota
	FrUTF8
)

var (
	ErrBadDst = errors.New("bad dst")
)

func (f *Frame) Config() *Config {
	return &Config{
		Flag:   f.flags,
		Color:  &f.Color,
		Font:   f.Font,
		Drawer: f.Drawer,
	}
}

func (c *Config) check() *Config {
	if c.Color == nil {
		c.Color = &A
	}
	if c.Font == nil {
		c.Font = NewGoMono(11)
	}
	if c.Drawer == nil {
		c.Drawer = &defaultDrawer{}
	}
	return c
}

func New(dst draw.Image, r image.Rectangle, conf *Config) *Frame {
	if dst == nil {
		return nil
	}
	if conf == nil {
		conf = &Config{}
	}
	conf.check()
	fl := conf.Flag
	var face Face
	switch f := conf.Font.(type) {
	case Face:
		face = f
	case font.Face:
		face = Open(f)
	}
	mintab, maxtab := tabMinMax(face, fl&FrElastic != 0)

	f := &Frame{
		Font:   face,
		Color:  *conf.Color,
		Drawer: conf.Drawer,
		Run:    box.NewRun(mintab, 5000, face),
		op:     draw.Src,
		mintab: mintab,
		maxtab: maxtab,
		flags:  fl,
	}
	f.setrects(r, dst)
	f.inittick()
	run := box.NewRun(mintab, 5000, face)
	f.ir = &run
	return f
}

// Frame is a write-only container for editable text
type Frame struct {
	box.Run
	p0 int64
	p1 int64
	b  draw.Image
	r  image.Rectangle
	ir *box.Run

	Font Face
	Color
	Ticked bool
	Scroll func(int)
	Drawer
	op draw.Op

	mintab int
	maxtab int
	full   int

	tick      draw.Image
	tickback  draw.Image
	tickscale int
	tickoff   bool
	maxlines  int
	modified  bool
	noredraw  bool

	pts [][2]image.Point

	flags int
}

// Flags returns the flags currently set for the frame
func (f *Frame) Flags() int {
	return f.flags
}

// Flag sets the flags for the frame. At this time
// only FrElastic is supported.
func (f *Frame) SetFlags(flags int) {
	fl := getflag(flags)
	f.flags = fl
	f.mintab, f.maxtab = tabMinMax(f.Font, f.elastic())
	//	f.Reset( f.r, f.RGBA(),f.Font)
	//	f.mintab, f.maxtab = tabMinMax(f.Font, f.elastic())
}

func (f *Frame) elastic() bool {
	return f.flags&FrElastic != 0
}

func tabMinMax(ft font.Face, elastic bool) (min, max int) {
	mintab := 5 //ft.Measure(' ')
	maxtab := mintab * 4
	if elastic {
		mintab = maxtab
	}
	return mintab, maxtab
}

func newRuneFrame(r image.Rectangle, ft font.Face, b draw.Image, cols Color, flag ...int) *Frame {
	/*
		fl := getflag(flag...)
		mintab, maxtab := tabMinMax(ft, fl&FrElastic != 0)

		f := &Frame{
			Font:         ft,
			mintab:       mintab,
			maxtab:       maxtab,
			Color:        cols,
			Run:          box.NewRun(mintab, 5000, ft, box.NewRuneRuler),
			stringBG:     font.RuneBG,
			stringNBG:    font.RuneNBG,
			op:           draw.Src,
			flags:        fl,
		}
		f.setrects(r, b)
		f.inittick()
		run := box.NewRun(mintab, 5000, ft, box.NewRuneRuler)
		f.ir = &run
		f.Drawer = drawcache.New()
		return f
	*/
	panic("disabled")
}

func getflag(flag ...int) (fl int) {
	if len(flag) != 0 {
		fl = flag[0]
	}
	if ForceElastic {
		fl |= FrElastic
	}
	if ForceUTF8 {
		fl |= FrUTF8
	}
	return fl
}

func (f *Frame) RGBA() *image.RGBA {
	return f.b.(*image.RGBA)
}
func (f *Frame) Size() image.Point {
	r := f.RGBA().Bounds()
	return image.Pt(r.Dx(), r.Dy())
}

// Dirty returns true if the contents of the frame have changes since the last redraw
func (f *Frame) Dirty() bool {
	return f.modified
}

// SetDirty alters the frame's internal state
func (f *Frame) SetDirty(dirty bool) {
	f.modified = dirty
}

func (f *Frame) SetOp(op draw.Op) {
	f.op = op

}

// Close closes the frame
func (f *Frame) Close() error {
	return nil
}

// Reset resets the frame to display on image b with bounds r and font ft.
func (f *Frame) Reset(r image.Rectangle, b *image.RGBA, ft font.Face) {
	f.r = r
	f.b = b
	f.SetFont(ft)
}

func (f *Frame) SetFont(ft font.Face) {
	f.Font = Open(ft)
	f.Run.Reset(f.Font)
	f.Refresh()
}

// Bounds returns the frame's clipping rectangle
func (f *Frame) Bounds() image.Rectangle {
	return f.r.Bounds()
}

// Full returns true if the last line in the frame is full
func (f *Frame) Full() bool {
	return f.full == 1
}

// Maxline returns the max number of wrapped lines fitting on the frame
func (f *Frame) MaxLine() int {
	return f.maxlines
}

// Line returns the number of wrapped lines currently in the frame
func (f *Frame) Line() int {
	return f.Nlines
}

// Len returns the number of bytes currently in the frame
func (f *Frame) Len() int64 {
	return f.Nchars
}

// Dot returns the range of the selected text
func (f *Frame) Dot() (p0, p1 int64) {
	return f.p0, f.p1
}

func (f *Frame) setrects(r image.Rectangle, b draw.Image) {
	f.b = b
	f.r = r
	h := f.Font.Dy()
	f.r.Max.Y -= f.r.Dy() % h
	f.maxlines = f.r.Dy() / h
}

func (f *Frame) clear(freeall bool) {
	if f.Nbox != 0 {
		f.Run.Delete(0, f.Nbox-1)
	}
	if freeall {
		f.tick = nil
		f.tickback = nil
	}
	f.Box = nil
	f.Ticked = false
}
