package font

import "image"

func deconvolveGlyph(g *Glyph) *Glyph {
	if g.mask == nil {
		panic("deconvolveGlyph: nil g.mask")
	}
	g.mask = deconvolve(g.mask)
	return g
}

func ampGlyph(g *Glyph) *Glyph {
	if g.mask == nil {
		panic("ampGlyph: nil g.mask")
	}
	g.mask = amp(g.mask)
	return g
}
func deconvolve(m *image.Alpha) *image.Alpha {
	av := func(x, y int) uint8 {
		a := m.AlphaAt(x, y)
		return a.A
	}
	for y := 0; y < m.Bounds().Dy(); y++ {
		for x := 0; x < m.Bounds().Dx(); x++ {
			a := m.AlphaAt(x, y)
			mean := int(
				av(x-1, y-1)+av(x, y-1)+av(x+1, y-1)+
					av(x-1, y-0)+av(x, y-0)+av(x+1, y-0)+
					av(x-1, y+1)+av(x, y+1)+av(x+1, y+1)) / 8
			if a.A-uint8(mean) > a.A {
				a.A = 0
			} else {
				a.A -= uint8(mean)
			}
			defer m.SetAlpha(x, y, a)
		}
	}
	return m
}

func amp(m *image.Alpha) *image.Alpha {
	for y := 0; y < m.Bounds().Dy(); y++ {
		for x := 0; x < m.Bounds().Dx(); x++ {
			a := m.AlphaAt(x, y)
			if a.A < 64 {
				continue
			}
			if a.A+64 < a.A {
				a.A = 255
			} else {
				a.A += 64
			}
			defer m.SetAlpha(x, y, a)
		}
	}
	return m
}
