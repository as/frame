package frame

import (
	"unicode"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	font.Face
	height int
}

func NewTTF(data []byte, size int) Font {
	f, err := truetype.Parse(data)
	if err != nil {
		panic(err)
	}
	return Font{
		Face: truetype.NewFace(f, &truetype.Options{
			Size: float64(size),
		}),
		height: size + size/5,
	}
}

func (f Font) Dx(s string) int {
	return f.stringwidth([]byte(s))
}
func (f Font) Dy() int {
	return f.height
}
func fix(i fixed.Int26_6) int {
	return i.Round()
}
func (f Font) stringwidth(p []byte) (w int) {
	for i := range p {
		w += f.Measure(rune(byte(p[i])))
	}
	return w
}
func (f Font) Measure(r rune) (q int) {
	if r == 0 || !unicode.IsGraphic(r) || r > 127 {
		return f.measureHex()
	}
	l, ok := f.Face.GlyphAdvance(r)
	if !ok {
		println("warn: glyph missing")
		l, _ = f.Face.GlyphAdvance('@')
	}
	return fix(l)
}

func (f Font) measureHex() int{
	return f.Measure('_')*7/4
}