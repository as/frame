package font

import (
	"fmt"
	"image"
	"image/draw"
	"unicode"

	"github.com/golang/freetype/truetype"
	gofont "golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	gofont.Face
	hexDx   int
	data    []byte
	size    int
	ascent  int
	descent int
	stride  int
	letting int
	dy      int

	cache    Cache
	hexCache Cache
	decCache Cache
}

func NewGoRegular(size int) *Font {
	return NewTTF(goregular.TTF, size)
}

func NewGoMedium(size int) *Font {
	return NewTTF(gomedium.TTF, size)
}

func NewGoMono(size int) *Font {
	return NewTTF(gomono.TTF, size)
}

// NewBasic always returns a 7x13 basic font
func NewBasic(size int) *Font {
	f := basicfont.Face7x13
	size = 13
	ft := &Font{
		Face:    f,
		size:    size,
		ascent:  2,
		descent: 1,
		letting: 0,
		stride:  0,
	}
	ft.dy = ft.ascent + ft.descent + ft.size
	hexFt := fromTTF(gomono.TTF, ft.Dy()/4+3)
	ft.hexDx = ft.genChar('_').Bounds().Dx()
	for i := 0; i != 256; i++ {
		ft.cache[i] = ft.genChar(byte(i))
		if ft.cache[i] == nil {
			ft.cache[i] = hexFt.genHexChar(ft.Dy(), byte(i))
		}
	}
	return ft
}

func NewTTF(data []byte, size int) *Font {
	ft := fromTTF(data, size)
	hexFt := fromTTF(gomono.TTF, ft.Dy()/4+3)

	ft.hexDx = ft.genChar('_').Bounds().Dx()
	for i := 0; i != 256; i++ {
		ft.cache[i] = ft.genChar(byte(i))
		if ft.cache[i] == nil {
			ft.cache[i] = hexFt.genHexChar(ft.Dy(), byte(i))
		}
	}
	return ft
}
func fromTTF(data []byte, size int) *Font {
	f, err := truetype.Parse(data)
	if err != nil {
		panic(err)
	}
	ft := &Font{
		Face: truetype.NewFace(f,
			&truetype.Options{
				Size:              float64(size),
				GlyphCacheEntries: 512 * 2,
				SubPixelsX:        1,
			}),
		size:    size,
		ascent:  2,
		descent: +(size / 3),
		stride:  0,
		data:    data,
	}
	ft.dy = ft.ascent + ft.descent + ft.size
	return ft
}

func (f *Font) NewSize(dy int) *Font {
	if dy == f.Dy() {
		return f
	}
	if f.data == nil {
		return NewBasic(dy)
	}
	return NewTTF(f.data, dy)
}

func (f *Font) SetAscent(px int) {
	f.ascent = px
	f.dy = f.ascent + f.descent + f.size
}
func (f *Font) SetDescent(px int) {
	f.descent = px
	f.dy = f.ascent + f.descent + f.size
}
func (f *Font) SetStride(px int) {
	f.stride = px
}
func (f *Font) SetLetting(px int) {
	f.letting = px
}

func (f *Font) genChar(b byte) *Glyph {
	dr, mask, maskp, adv, _ := f.Face.Glyph(fixed.P(0, f.size), rune(b))
	if !f.Printable(b) {
		return nil
	}
	r := image.Rect(0, 0, Fix(adv), f.Dy())
	m := image.NewAlpha(r)
	r = r.Add(image.Pt(dr.Min.X, dr.Min.Y))
	draw.Draw(m, r, mask, maskp, draw.Src)
	return &Glyph{mask: m, Rectangle: m.Bounds()}
}
func (f *Font) genHexChar(dy int, b byte) *Glyph {
	s := fmt.Sprintf("%02x", b)
	g0 := f.genChar(s[0])
	g1 := f.genChar(s[1])
	r := image.Rect(2, f.descent+f.ascent, g0.Bounds().Dx()+g1.Bounds().Dx()+6, dy)
	m := image.NewAlpha(r)
	draw.Draw(m, r, g0.Mask(), image.ZP, draw.Over)
	r.Min.X += g0.Mask().Bounds().Dx()
	draw.Draw(m, r.Add(image.Pt(-f.descent/4, f.descent*2)), g1.Mask(), image.ZP, draw.Over)
	return &Glyph{mask: m, Rectangle: m.Bounds()}
}

func (f *Font) Char(b byte) (mask *image.Alpha) {
	return f.cache[b].mask
}

func (f *Font) Descent() int {
	return f.descent
}

func (f *Font) Dx(s string) int {
	return f.MeasureBytes([]byte(s))
}
func (f *Font) Dy() int {
	return f.dy + f.letting
}
func (f *Font) Size() int {
	return f.size
}
func Fix(i fixed.Int26_6) int {
	return i.Round()
}
func (f *Font) MeasureBytes(p []byte) (w int) {
	for i := range p {
		w += f.Measure(rune(byte(p[i])))
	}
	return w
}

func (f *Font) MeasureByte(b byte) (n int) {
	return f.cache[b].Dx() + f.stride
}

func (f *Font) MeasureRune(r rune) (q int) {
	advance, _ := f.Face.GlyphAdvance(r)
	return Fix(advance)
}

func (f *Font) Measure(r rune) (q int) {
	return f.cache[byte(r)].Dx() + f.stride
}

func (f *Font) MeasureHex() int {
	return f.hexDx
}

func (f *Font) TTF() []byte {
	return f.data
}

func (f *Font) Printable(b byte) bool {
	if b == 0 || b > 127 {
		return false
	}
	if unicode.IsGraphic(rune(b)) {
		return true
	}
	return false
}
