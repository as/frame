package box

//import "log"

func (r *Run) ensure(nb int) {
	if nb == r.Nalloc {
		r.Grow(r.delta)
		if r.delta < 32768 {
			r.delta *= 2
		}
	}
}

func (r *Run) Boxscan(s []byte, ymax int) {
	r.Nchars += int64(len(s))
	i := 0
	nb := 0
	for nl := 0; nl <= ymax; nb++ {
		if nb == r.Nalloc {
			r.Grow(r.delta)
			if r.delta < 16384 {
				r.delta <<= 1
			}
		}
		if i == len(s) {
			break
		}
		i++
		c := s[i-1]
		switch c {
		default:
			for _, c = range s[i:] {
				if special(c) {
					break
				}
				i++
			}
			r.Box[nb] = Box{
				Nrune: i,
				Ptr:   s[:i],
				Width: r.MeasureBytes(s[:i]),
			}
		case '\t':
			r.Box[nb] = Box{
				Nrune:    -1,
				Ptr:      s[:i],
				Width:    r.minDx,
				Minwidth: r.minDx,
			}
		case '\n':
			r.Box[nb] = Box{
				Nrune: -1,
				Ptr:   s[:i],
				Width: r.maxDx,
			}
			nl++
		}
		s = s[i:]
		i = 0
	}
	r.Nchars -= int64(len(s))
	r.Nbox += nb
}

func special(c byte) bool {
	return c == '\t' || c == '\n'
}
