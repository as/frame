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
		c.Font = NewGoMono(11)
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
	return Dy(f) / 2
}
func Height(f font.Face) int  { return f.Metrics().Height.Ceil() }
func Ascent(f font.Face) int  { return f.Metrics().Ascent.Ceil() }
func Descent(f font.Face) int { return f.Metrics().Descent.Ceil() }

func Dy(f font.Face) int {
	return Height(f) + Height(f)/2
}
func NewGoMono(size int) font.Face {
	return truetype.NewFace(gomonoTTF, &truetype.Options{
		SubPixelsX: 64,
		SubPixelsY: 64,
		Hinting:    font.HintingFull,
		Size:       float64(size),
	})
}

func StaticFace(f font.Face) font.Face{
	return &staticFace{
		cache: make(map[signature]*image.RGBA),
		Face: f,
	}
}

type signature struct {
	b byte
	fg color.RGBA
	bg color.RGBA
}
type staticFace struct {
	cache map[signature]*image.RGBA
	font.Face
}

func (s *staticFace) RawGlyph(b byte, fg color.RGBA, bg color.RGBA) image.Image {
	if img, ok := s.cache[signature{b, fg,bg}]; ok {
		return img
	}
	mask, r := s.genChar(b)
	img := image.NewRGBA(r)
	draw.Draw(img, img.Bounds(), image.NewUniform(bg), image.ZP, draw.Src)
	draw.DrawMask(img, img.Bounds(), image.NewUniform(fg), image.ZP, mask, r.Min, draw.Over)
	s.cache[signature{b, fg,bg}] = img
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
	return &staticFace{
		cache: make(map[signature]*image.RGBA),
		Face: truetype.NewFace(gomonoTTF, &truetype.Options{
			SubPixelsX: 64,
			SubPixelsY: 64,
			Hinting:    font.HintingFull,
			Size:       float64(size),
		}),
	}
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
