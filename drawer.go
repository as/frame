package frame

import (
	"image"
	"image/draw"

	. "github.com/as/font"
	"golang.org/x/image/font"
)

func NewDefaultDrawer() Drawer {
	return &defaultDrawer{}
}

type defaultDrawer struct{}

func (d *defaultDrawer) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, op draw.Op) {
	draw.Draw(dst, r, src, sp, op)
}

func (d *defaultDrawer) StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft font.Face, s []byte, bg image.Image, bgp image.Point) int {
	return StringBG(dst, p, src, sp, ft, s, bg, bgp)
}

func (d *defaultDrawer) Flush(r ...image.Rectangle) error {
	return nil
}
