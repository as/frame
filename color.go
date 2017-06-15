package frame

import (
	"image"
	"image/color"
)

var (
	Black  = image.NewUniform(color.RGBA{0, 0, 0, 255})
	Red    = image.NewUniform(color.RGBA{255, 0, 0, 255})
	Green  = image.NewUniform(color.RGBA{0, 255, 0, 255})
	Blue   = image.NewUniform(color.RGBA{0, 192, 192, 255})
	Cyan   = image.NewUniform(color.RGBA{234, 255, 255, 255})
	White  = image.NewUniform(color.RGBA{255, 255, 255, 255})
	Yellow = image.NewUniform(color.RGBA{255, 255, 224, 255})
	Gray   = image.NewUniform(color.RGBA{0x55, 0x55, 0x55, 255})
	Mauve  = image.NewUniform(color.RGBA{0x99, 0x99, 0xDD, 255})
)

type Color struct {
	Pallete
	Hi Pallete
}

var Acme = Color{
	Pallete: Pallete{Text: Gray, Back: Yellow},
	Hi:      Pallete{Text: White, Back: Mauve},
}

type Pallete struct {
	Text, Back image.Image
}
