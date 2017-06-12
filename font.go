package frame

import (
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
func (f Font) Measure(r rune) int {
	if r == 0 {
		r = 1
	}
	l, ok := f.Face.GlyphAdvance(r)
	if !ok {
		println("warn: glyph missing")
		l, _ = f.Face.GlyphAdvance('@')
	}
	return fix(l)
}
