package font

import "image"

type Cache [256]*Glyph

type Glyph struct {
	mask *image.Alpha
	image.Rectangle
}

func (g Glyph) Mask() *image.Alpha {
	return g.mask
}

func (g Glyph) Bounds() image.Rectangle {
	return g.Rectangle
}
