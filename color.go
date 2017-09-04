package frame

import (
	"image"
	"image/color"
)

var (
	Yellow = image.NewUniform(color.RGBA{255, 255, 224, 255})
	Green  = image.NewUniform(color.RGBA{0x99, 0xCC, 0x99, 255})
	Red    = image.NewUniform(color.RGBA{0xCC, 0x99, 0x99, 255})
	Gray   = image.NewUniform(color.RGBA{0x55, 0x55, 0x55, 255})
	Mauve  = image.NewUniform(color.RGBA{0x99, 0x99, 0xDD, 255})

	Ozone  = image.NewUniform(color.RGBA{216, 216, 232, 255})
	Strata = image.NewUniform(color.RGBA{248, 242, 248, 255})
	Peach  = image.NewUniform(color.RGBA{255, 248, 232, 255})
)

type Color struct {
	Palette
	Hi Palette
}

var (
	Acme = Color{
		Palette: Palette{Text: Gray, Back: Yellow},
		Hi:      Palette{Text: image.White, Back: Mauve},
	}
	Mono = Color{
		Palette: Palette{Text: image.Black, Back: image.White},
		Hi:      Palette{Text: image.White, Back: image.Black},
	}
	A = Color{
		Palette: Palette{
			Text: Gray,
			Back: Peach,
		},
		Hi: Palette{
			Back: Mauve,
			Text: image.White,
		},
	}
	ATag0 = Color{
		Palette: Palette{
			Text: Gray,
			Back: Ozone,
		},
		Hi: Palette{
			Back: Mauve,
			Text: image.White,
		},
	}
	ATag1 = Color{
		Palette: Palette{
			Text: Gray,
			Back: Strata,
		},
		Hi: Palette{
			Back: Mauve,
			Text: image.White,
		},
	}
)

type Palette struct {
	Text, Back image.Image
}
