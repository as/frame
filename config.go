package frame

import (
	"image"
	"image/draw"
	"unicode"

	"golang.org/x/image/font"
)

type Config struct {
	Flag   int
	Scroll func(int)
	Color  *Color
	Font   font.Face
	Drawer Drawer
}

// Drawer implements the set of methods a frame needs to draw on a draw.Image. The frame's default behavior is to use
// the native image/draw package and x/exp/font packages to satisfy this interface.
type Drawer interface {
	Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, op draw.Op)
	//DrawMask(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op)

	// StringBG draws a string to dst at point p
	StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft font.Face, s []byte, bg image.Image, bgp image.Point) int

	// Flush requests that prior calls to the draw and string methods are flushed from an underlying soft-screen. The list of rectangles provide
	// optional residency information. Implementations may refresh a superset of r, or ignore it entirely, as long as the entire region is
	// refreshed
	Flush(r ...image.Rectangle) error

	// Cache returns the set of rectangles that have been updates but not flushed. This method exists
	// temporarily and will be removed from this implementation. Frame does not use it.
	//Cache() []image.Rectangle
}
type Palette struct {
	Text, Back image.Image
}
type Color struct {
	Palette
	Hi Palette
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
