package frame

import (
	"image"
	"image/color"
)

func solid(r, g, b byte) *image.Uniform {
	return uniform(r, g, b, 255)
}
func uniform(r, g, b, a byte) *image.Uniform {
	return image.NewUniform(color.RGBA{r, g, b, a})
}

var (
	Yellow = solid(255, 255, 224)
	Green  = solid(0x99, 0xCC, 0x99)
	Red    = solid(0xCC, 0x99, 0x99)
	Gray   = solid(0x12, 0x12, 0x12)
	Mauve  = solid(0x99, 0x99, 0xDD)

	Ozone     = solid(216, 216, 232)
	Strata    = solid(248, 242, 248)
	AntiPeach = solid(0, 12, 24)
	Peach     = solid(255, 248, 232)

	BBody = solid(243, 248, 254)
	BTag  = solid(214, 230, 252)
	GBody = solid(226, 235, 232)
	GTag  = solid(226, 225, 232)
	PTag  = solid(222, 207, 236)
	PBody = solid(252, 232, 252)

	MMauve = solid(0x66, 0x55, 0x88)
	MTagG  = solid(28-10, 31-13, 38-15)
	MTagC  = MTagW
	MTagW  = solid(28, 31, 38)
	MBodyW = solid(43, 50, 59)
	MTextW = solid(255-59, 255-50, 255-43)
)

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
			Text: AntiPeach,
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
