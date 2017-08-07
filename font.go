package frame

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
	"unicode"
)

type Font struct {
	font.Face
	height int
	data   []byte
}

func NewGoRegular(size int) Font {
	return NewTTF(goregular.TTF, size)
}

func NewGoMono(size int) Font {
	return NewTTF(gomono.TTF, size)
}

func NewTTF(data []byte, size int) Font {
	f, err := truetype.Parse(data)
	if err != nil {
		panic(err)
	}
	return Font{
		Face: truetype.NewFace(f,
			&truetype.Options{
				Size:              float64(size),
				GlyphCacheEntries: 512 * 2,
				SubPixelsX:        1,
			}),
		height: size + size/5,
		data:   data,
	}
}

func (f Font) Dx(s string) int {
	return f.MeasureBytes([]byte(s))
}
func (f Font) Dy() int {
	return f.height
}
func (f Font) Size() int {
	return 5 * f.Dy() / 6
}
func fix(i fixed.Int26_6) int {
	return i.Round()
}
func (f Font) MeasureBytes(p []byte) (w int) {
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
		l, _ = f.Face.GlyphAdvance('_')
	}
	return fix(l)
}

func (f Font) measureHex() int {
	return f.Measure('_') * 2
}

func (f Font) TTF() []byte {
	return f.data
}
