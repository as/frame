package frame

import (
	"image"
//		"fmt"
)

// IndexOf computes the index of the glyph containing pt
func (f *Frame) IndexOf(pt image.Point) (i int) {
//	defer func() {fmt.Printf("IndexOf: pt=%v i=%d (%c)\n", pt, i, f.s[i])}()
	pt = f.downscale(pt)
	pt = f.alignY(pt)			
	qt := f.alignY(f.origin)	
	s := f.s[:f.nbytes]
	w, maxw := 0, f.Wrap
	breakline := func() {
		qt.Y += f.FontHeight()
		qt.X, w = f.origin.X, 0
	}

	for {
		if i >= len(s) || qt.Y > pt.Y {
			break
		}
		switch c := s[i]; {
		case c == '\n':
			breakline()
			i++
		case w >= maxw:
			breakline()
		default:
			qt.X += f.Advance(rune(c))
			w++
			i++	
			if qt.Y == pt.Y && qt.X-3 >= pt.X {
				return i-1
			} 
		}
	}
	if len(f.s) > i && i > 1 && f.s[i-1] == '\n'{
		return i-1
	}
	return i
}

// PointOf computes the point of origin for glyph i
func (f *Frame) PointOf(i int) (pt image.Point) {
//	defer func(){fmt.Printf("PointOf: pt=%v i=%d (%c)\n", pt, i, f.s[i])}()
	if i < 0 {
		i = 0
	}	
	pt = f.alignY(f.origin)
	s, j := f.s[:i], 0
	w, maxw := 0, f.Wrap
	breakline := func() {
		pt.Y += f.FontHeight()
		pt.X, w = f.origin.X, 0
	}
	for {
		if j >= len(s) {
			break
		}
		switch c := s[j]; {
		case c == '\n':
			if j == i {
				return pt
			}
			breakline()
			j++
		case w >= maxw:
			breakline()
		default:
			pt.X += f.Advance(rune(c))
			w++
			j++
		}
	}
	return pt
}

// PointWalk walks from index s to index e. It returns the point of
// origin for glyph at index e.
func (f *Frame) PointWalk(s, e int, sp image.Point) (ep image.Point) {
	if s > e{
		return f.PointOf(e)
	}
	if s < 0 {
		s = 0
	}
	sp = f.alignY(sp)
	data := f.s
	w, maxw := 0, f.Wrap
	breakline := func() {
		sp.Y += f.FontHeight()
		sp.X = f.origin.X	
		w = 0 	// If reverse (dir=-1), set width to max
	}
	for {
		if s >= len(data) || s >= e {
			break
		}
		switch c := data[s]; {
		case c == '\n':
			if s == e {
				return sp
			}
			breakline()
			s++
		case w >= maxw:
			breakline()
		default:
			sp.X += f.Advance(rune(c))
			w++
			s++
		}
	}
	return sp
}

// IndexWalk walks from index s at point sp to the terminus, ep.
// It returns the index of the glyph under ep.
func (f *Frame) IndexWalk(sp, ep image.Point, s int) (e int) {
	if sp.Y > ep.Y || sp.Y == ep.Y && sp.X > ep.X {
		return f.IndexOf(ep)
	}
	
	sp = f.alignY(sp)
	ep = f.alignY(ep)
	data := f.s[:f.nbytes]
	w, maxw := 0, f.Wrap
	breakline := func() {
		sp.Y += f.FontHeight()
		sp.X, w = f.origin.X, 0
	}

	for {
		if s >= len(data) || sp.Y > ep.Y {
			break
		}
		switch c := data[s]; {
		case c == '\n':
			breakline()
			s++
		case w >= maxw:
			breakline()
		default:
			sp.X += f.Advance(rune(c))
			w++
			s++	
			if sp.Y == ep.Y && sp.X-3 >= ep.X {
				return s-1
			} 
		}
	}
	if len(data) > s && s > 1 && data[s-1] == '\n'{
		return s-1
	}
	return s
}

func (f *Frame) Advance(r rune) int {
	advance, _ := f.Font.GlyphAdvance(r)
	return int(float64(advance >> 6))
}
