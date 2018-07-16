package frame

import (
	"image"
	"image/color"
)

var (
	// Common colors found in Acme
	Black  = image.Black
	White  = image.White
	Blue   = &image.Uniform{color.RGBA{0x99, 0x99, 0xCC, 0xFF}}
	Yellow = &image.Uniform{color.RGBA{0xFF, 0xFF, 0xFD, 0xFF}}
	Gray   = &image.Uniform{color.RGBA{0x90, 0x0C, 0xC0, 0xFF}}
	Red    = &image.Uniform{color.RGBA{0xCC, 0x99, 0x99, 0xFF}}
	Green  = &image.Uniform{color.RGBA{0x99, 0xCC, 0x99, 0xFF}}

	// Other colors
	Peach = &image.Uniform{color.RGBA{0xFF, 0xF8, 0xE8, 0xFF}}
	Mauve = &image.Uniform{color.RGBA{0x90, 0x90, 0xC0, 0xFF}}
)

var (
	// Acme is the color scheme found in the Acme text editor
	Acme = NewColor(Gray, Yellow, White, Blue)
	Mono = NewColor(Black, White, White, Black)
	A    = NewColor(Gray, Peach, White, Mauve)
)

// NewColor returns a Color for the given foreground and background
// images. Two extra colors may be provided to set the highlighted
// foreground and background image palette.
func NewColor(fg, bg image.Image, hi ...image.Image) Color {
	c := Color{Palette: Palette{Text: fg, Back: bg}}
	if len(hi) > 0 {
		c.Hi.Text = hi[0]
	}
	if len(hi) > 1 {
		c.Hi.Back = hi[1]
	}
	return c
}

// NewUniform is like NewColor, only it accepts arguments
// as color.Color values instead of images.
func NewUniform(fg, bg color.Color, hi ...color.Color) Color {
	var img = make([]image.Image, 0, len(hi))
	if len(hi) > 0 {
		img = append(img, image.NewUniform(hi[0]))
	}
	if len(hi) > 1 {
		img = append(img, image.NewUniform(hi[1]))
	}
	return NewColor(image.NewUniform(fg), image.NewUniform(bg), img...)
}

// Color is constructed from a Palette pair. The Hi Palette describes
// the appearance of highlighted text.
type Color struct {
	Palette
	Hi Palette
}

// Pallete contains two images used to paint text and backgrounds
// on the frame.
type Palette struct {
	Text, Back image.Image
}
