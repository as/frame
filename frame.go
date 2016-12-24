package frame

import (
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"golang.org/x/image/font/gofont/goregular"
	"github.com/golang/freetype/truetype"
)

var (
	AcmeColors = &Colors{
		Back:  image.NewUniform(color.RGBA{0,0,0,0}),
		Text:  image.NewUniform(color.RGBA{255,0,0,255}),
		HText: image.NewUniform(color.RGBA{0,0,0,255}),
		HBack: image.NewUniform(color.RGBA{255,0,0,128}),
	}
	DefaultColors = defaultColors
	DarkGrayColors = &Colors{
		Back:  image.NewUniform(color.RGBA{33,33,33,4}),
		Text:  image.NewUniform(color.RGBA{128,128,128,255}),
		HText: image.NewUniform(color.RGBA{0,0,0,255}),
		HBack: image.NewUniform(color.RGBA{0,128,128,64}),
	}
	GrayColors = &Colors{
		Back:  image.NewUniform(color.RGBA{48,48,48,0}),
		Text:  image.NewUniform(color.RGBA{99,99,99,255}),
		HText: image.NewUniform(color.RGBA{0,0,0,255}),
		HBack: image.NewUniform(color.RGBA{0,128,128,64}),
	}
	defaultColors = &Colors{
		Back:  image.NewUniform(color.RGBA{66,66,66,0}),
		Text:  image.NewUniform(color.RGBA{0,255,255,255}),
		HText: image.NewUniform(color.RGBA{0,0,0,255}),
		HBack: image.NewUniform(color.RGBA{132,255,255,0}),
	}
	defaultOption = &Option{
		Font:   parseDefaultFont(12),
		Wrap:   80,
		Scale:  image.Pt(1, 1),
		Colors: *defaultColors,
	}
	largeScale = &Option{
		Font:   parseDefaultFont(24),
		Wrap:   80,
		Scale:  image.Pt(1, 1),
		Colors: *defaultColors,
	}
)

type Frame struct {
	disp   draw.Image
	r      image.Rectangle
	origin image.Point

	Option
	Tick *Tick

	s      []byte
	width  int
	nbytes int
	Dirty  bool

	// cache for the transformation
	cached draw.Image
	last   rune
}

type Colors struct {
	Text, Back   image.Image
	HText, HBack image.Image
}

type Option struct {
	// Font is the font face for the frame
	Font font.Face

	// Number of glyphs drawn on one line before wrapping
	Wrap int

	// Multiplicative scale factor for X and Y coordinates
	// (1, 1) means no scale.
	Scale image.Point

	// Colors define the text and background colors for the rame
	// Text: glyph color
	// Back: background color
	// HText: highlighted glyph color
	// HBack: highlighted background color
	Colors Colors

	fontheight int
}

// New initializes a new frame on disp. The first glyph of
// text is inserted at point p. If opt is nil, the default
// color, wrapping, and font settings are used.
func New(disp draw.Image, p image.Point, opt *Option) *Frame {
	if opt == nil {
		opt = defaultOption
	}
	f := &Frame{
		disp:   disp,
		r:      disp.Bounds(),
		origin: p,
		Option: *opt,
		s:      make([]byte, 64),
	}
	f.cached = image.NewRGBA(image.Rectangle{image.ZP, image.Pt(f.FontHeight(), f.FontHeight())})
	f.flushcache()
	return f
}

func (o Option) FontHeight() int {
	if o.Font == nil {
		return 0
	}
	if o.fontheight == 0 {
		o.fontheight = int(o.Font.Metrics().Height>>6) + 1
	}
	return o.fontheight
}

func ParseDefaultFont(size float64) font.Face{
	return parseDefaultFont(size)
}

func parseDefaultFont(size float64) font.Face{
	f, err := truetype.Parse(goregular.TTF)
	if err != nil{
		panic(err)
	}
	return truetype.NewFace(f, &truetype.Options{
		Size: size,
	})
}

func (f *Frame) SetDisp(disp draw.Image) {
	f.disp = disp
}

// Insert inserts s starting from index i in the
// the frame buffer.
func (f *Frame) Insert(s []byte, i int) (err error) {
	if i >= len(f.s) {
		i = len(f.s) - 1
	}
	if i < 0 {
		i = 0
	}
	if s == nil {
		return nil
	}
	f.grow(len(s))
	f.nbytes += len(s)
	copy(f.s[i+len(s):], f.s[i:])
	copy(f.s[i:], s)
	f.Dirty = true
	return nil
}

// Delete erases the range [i:j] in the framebuffer
// TODO: fix i == j
func (f *Frame) Delete(i, j int) (err error) {
	if i > j {
		i, j = j, i
	}
	if i < 0 {
		i = 0
	}
	if j >= len(f.s) {
		j = len(f.s) - 1
	}
	copy(f.s[i:], f.s[j:])
	f.nbytes -= j - i
	if f.nbytes < 0{
		f.nbytes = 0
	}
	f.Dirty = true
	return nil
}

func (f *Frame) Bytes() []byte{
	return f.s
}

func (f *Frame) Bounds() (r image.Rectangle){
	return f.disp.Bounds()
}
