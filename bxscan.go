package frame

import (
	"github.com/as/frame/box"
	"image"
	//	"log"
)

func (f *Frame) bxscan(s []byte, ppt image.Point) (image.Point, image.Point) {

	const Delta = 25
	var (
		w, nb, delta,
		nl, rw int
		b   *box.Box
		tmp [256 + 3]byte
	)
	//	log.Printf("bxscan: s=%s ppt=%s\n", s, ppt)
	fr := f.fr
	fr.Reset(f.r, f.b, f.Font)
	fr.maxtab = 4 * f.Dx(" ")
	fr.Color = f.Color
	delta = Delta
	nl = 0
	min := fr.Font.Measure(' ')
	for nb = 0; len(s) > 0 && nl <= f.maxlines; nb, fr.Nbox = nb+1, fr.Nbox+1 {
		if nb == fr.Nalloc {
			fr.Grow(delta)
			if delta < 10000 {
				delta *= 2
			}
		}
		b = &fr.Box[nb]
		c := s[0]
		if c == '\t' || c == '\n' {
			b.BC = c
			b.Ptr = []byte{c}
			b.Nrune = -1
			b.Width = 5000
			if c == '\n' {
				b.Minwidth = 0
				nl++
			} else {
				b.Minwidth = min
			}
			fr.Nchars++
			s = s[1:]
		} else {
			b.BC = c
			tp := 0 // index into tmp
			nr := 0
			w = 0
			for len(s) > 0 {
				c = s[0]
				if c == '\t' || c == '\n' {
					break
				}
				// TODO: runetochar: runes can be > 1 char
				tmp[tp] = c
				rw = 1
				if tp+rw >= len(tmp) {
					break
				}
				w += f.Font.MeasureBytes([]byte(s[:1]))
				s = s[1:]
				tp += rw
				nr++
			}
			p := make([]byte, tp)
			b = &fr.Box[nb]
			b.Ptr = p
			copy(p, tmp[:tp])
			b.Width = w
			b.Nrune = nr
			fr.Nchars += int64(nr)
		}
	}
	ppt = f.lineWrap0(ppt, &fr.Box[0])
	return ppt, fr.drawAt(ppt)
}
