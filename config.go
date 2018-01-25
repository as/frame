package frame

import (
	"image"
	"image/color"
	"image/draw"
	"unicode"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
)

type Config struct {
	Flag   int
	Scroll func(int)
	Color  *Color
	Font   font.Face
	Drawer Drawer
}

func (c *Config) check() *Config {
	if c.Color == nil {
		c.Color = &A
	}
	if c.Font == nil {
		c.Font = StaticFace(NewGoMono(11))
	}
	if c.Drawer == nil {
		c.Drawer = &defaultDrawer{}
	}
	return c
}

func (f *Frame) Config() *Config {
	return &Config{
		Flag:   f.flags,
		Color:  &f.Color,
		Font:   f.Font,
		Drawer: f.Drawer,
	}
}

var gomonoTTF, _ = truetype.Parse(gomono.TTF)

func Letting(f font.Face) int {
	if f, ok := f.(*staticFace); ok {
		return f.dy
	}
	return Dy(f) / 2
}
func Height(f font.Face) int {
	if f, ok := f.(*staticFace); ok {
		return f.h
	}
	return f.Metrics().Height.Ceil()
}
func Ascent(f font.Face) int {
	if f, ok := f.(*staticFace); ok {
		return f.a
	}
	return f.Metrics().Ascent.Ceil()
}
func Descent(f font.Face) int {
	if f, ok := f.(*staticFace); ok {
		return f.d
	}
	return f.Metrics().Descent.Ceil()
}
func Dy(f font.Face) int {
	if f, ok := f.(*staticFace); ok {
		return f.dy
	}
	return Height(f) + Height(f)/2
}
func NewGoMono(size int) font.Face {
	return truetype.NewFace(gomonoTTF, &truetype.Options{
		SubPixelsX: 64,
		SubPixelsY: 64,
		Size:       float64(size),
	})
}

func StaticFace(f font.Face) font.Face {
	if _, ok := f.(*staticFace); ok {
		return f
	}
	return &staticFace{
		a:     Ascent(f),
		d:     Descent(f),
		h:     Height(f),
		l:     Letting(f),
		dy:    Height(f) + Height(f)/2,
		cache: make(map[signature]*image.RGBA),
		Face:  f,
	}
}

func convert(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)}
}

type signature struct {
	b  byte
	fg color.RGBA
	bg color.RGBA
}
type staticFace struct {
	h, a, d, l, dy int
	cache          map[signature]*image.RGBA
	font.Face
}

func (s *staticFace) RawGlyph(b byte, fg, bg color.Color) image.Image {
	sig := signature{b, convert(fg), convert(bg)}
	if img, ok := s.cache[sig]; ok {
		return img
	}
	mask, r := s.genChar(b)
	img := image.NewRGBA(r)
	draw.Draw(img, img.Bounds(), image.NewUniform(bg), image.ZP, draw.Src)
	draw.DrawMask(img, img.Bounds(), image.NewUniform(fg), image.ZP, mask, r.Min, draw.Over)
	s.cache[sig] = img
	return img
}
func (f *staticFace) genChar(b byte) (*image.Alpha, image.Rectangle) {
	dr, mask, maskp, adv, _ := f.Face.Glyph(fixed.P(0, Height(f.Face)), rune(b))
	r := image.Rect(0, 0, Fix(adv), Dy(f.Face))
	m := image.NewAlpha(r)
	r = r.Add(image.Pt(dr.Min.X, dr.Min.Y))
	draw.Draw(m, r, mask, maskp, draw.Src)
	return m, m.Bounds()
}
func StaticGoMono(size int) font.Face {
	return StaticFace(NewGoMono(size))
}

func Printable(b byte) bool {
	if b == 0 || b > 127 {
		return false
	}
	if unicode.IsGraphic(rune(b)) {
		return true
	}
	return false
}
