package frame

import (
	"github.com/as/drawcache"
	"github.com/as/frame/box"
	"github.com/as/frame/font"
	"image"
	"image/draw"
)

var (
	// Enables the UTF-8 experiment
	ForceUTF8Experiment = false
	// Enables the Elastic Tabstop experiement
	ForceElasticTabstopExperiment = false
)

// Frame is a write-only container for editable text
type Frame struct {
	box.Run
	Color
	Font         *font.Font
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
	
	// Points to the font subpackage's StringN?BG or RuneN?BG functions
	stringBG  func (draw.Image, image.Point, image.Image, image.Point, *font.Font, []byte, image.Image, image.Point) int
	stringNBG func (draw.Image, image.Point, image.Image, image.Point, *font.Font, []byte) int
	newRulerFunc func(s []byte, ft *font.Font) box.Ruler

	drawcache.Drawer
	pts     [][2]image.Point
	Scroll  func(int)
	ir      *box.Run
	hexFont *font.Font
	hex     []draw.Image
}

func newRuneFrame(r image.Rectangle, ft *font.Font, b *image.RGBA, cols Color, runes ...bool) *Frame{
	spaceDx := ft.Measure(' ')
	f := &Frame{
		Font:   ft,
		maxtab: 4 * spaceDx,
		Color:  cols,
		Run:    box.NewRun(spaceDx, 5000, ft, box.NewRuneRuler),
		stringBG: font.RuneBG,
		stringNBG: font.RuneNBG,
		newRulerFunc: box.NewRuneRuler,
		op:     draw.Src,
	}
	f.setrects(r, b)
	f.inittick()
	run := box.NewRun(spaceDx, 5000, ft, box.NewRuneRuler)
	f.ir = &run
	f.Drawer = drawcache.New()
	return f
}

// New creates a new frame on b with bounds r. The image b is used
// as the frame's internal bitmap cache.
func New(r image.Rectangle, ft *font.Font, b *image.RGBA, cols Color, runes ...bool) *Frame {
	if (len(runes) > 0 && runes[0]) || ForceUTF8Experiment {
		return newRuneFrame(r,ft,b,cols)
	}
	spaceDx := ft.Measure(' ')
	f := &Frame{
		Font:   ft,
		maxtab: 4 * spaceDx,
		Color:  cols,
		Run:    box.NewRun(spaceDx, 5000, ft),
		stringBG: font.StringBG,
		stringNBG: font.StringNBG,
		newRulerFunc: box.NewByteRuler,
		op:     draw.Src,
	}
	f.setrects(r, b)
	f.inittick()
	run := box.NewRun(spaceDx, 5000, ft)
	f.ir = &run
	f.Drawer = drawcache.New()
	return f
}

func (f *Frame) RGBA() *image.RGBA {
	return f.b
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
func (f *Frame) Close() error{
	return nil
}

// Reset resets the frame to display on image b with bounds r and font ft.
func (f *Frame) Reset(r image.Rectangle, b *image.RGBA, ft *font.Font) {
	f.r = r
	f.b = b
	f.SetFont(ft)
}

func (f *Frame) SetFont(ft *font.Font) {
	f.Font = ft
	f.Run.Reset(ft)
	f.Refresh()
}

// Bounds returns the frame's clipping rectangle
func (f *Frame) Bounds() image.Rectangle {
	return f.r.Bounds()
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

// Len returns the number of bytes currently in the frame
func (f *Frame) Len() int64 {
	return f.Nchars
}

// Dot returns the range of the selected text
func (f *Frame) Dot() (p0, p1 int64) {
	return f.p0, f.p1
}

func (f *Frame) setrects(r image.Rectangle, b *image.RGBA) {
	f.b = b
	f.entire = r
	f.r = r
	f.r.Max.Y -= f.r.Dy() % f.Font.Dy()
	f.maxlines = f.r.Dy() / f.Font.Dy()
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
