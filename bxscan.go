package frame

import (
	"image"
	"io"
	"unicode/utf8"
	"github.com/as/frame/box"
	"github.com/as/frame/font"
)

// bxscan resets the measuring function and calls Bxscan in the embedded run
func (f *Frame) bxscan(s []byte, ppt image.Point) (image.Point, image.Point) {
	f.ir.Reset(f.ir.Measure)
	f.ir.Bxscan(s, f.maxlines)
	ppt = f.lineWrap0(ppt, &f.ir.Box[0])
	return ppt, f.drawRun(f.ir, ppt)
}

type byteRuler struct {
	*font.Font
	b *box.Box
	i int
}

func newByteRuler(b *box.Box, ft *font.Font) *byteRuler {
	return &byteRuler{Font: ft, b: b}
}

func (bs *byteRuler) Next() (size, widthPx int, err error) {
	if bs.i == bs.b.Nrune {
		return 0, 0, io.EOF
	}
	size = 1
	widthPx = bs.Font.MeasureByte(bs.b.Ptr[bs.i])
	bs.i += size
	return
}

type runeRuler struct {
	*font.Font
	b *box.Box
	i int
}

func newRuneRuler(b *box.Box, ft *font.Font) runeRuler {
	return runeRuler{Font: ft, b: b}
}

func (bs runeRuler) Next() (size, widthPx int, err error) {
	if bs.i == bs.b.Nrune {
		return 0, 0, io.EOF
	}
	r, size := utf8.DecodeRune(bs.b.Ptr[bs.i:])
	widthPx = bs.MeasureRune(r)
	bs.i += size
	return
}
