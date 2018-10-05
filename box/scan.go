package box

import "unicode/utf8"

//import "log"

func (r *Run) ensure(nb int) {
	if nb != r.Nalloc {
		return
	}
	r.Grow(r.delta)
	if r.delta < 32768 {
		r.delta *= 2
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *Run) zRunescan(s []byte, ymax int) {
	r.Nbox = 0
	r.Nchars = 0
	r.Nchars += int64(len(s))
	i := 0
	nb := 0

	adv := 0
	for nl := 0; nl <= ymax; nb++ {
		if nb == r.Nalloc {
			r.Grow(r.delta)
			if r.delta < 32768 {
				r.delta *= 2
			}
		}
		i += adv
		if i == len(s) {
			break
		}
		c := s[i]
		switch c {
		default:
			for _, c := range string(s[i:min(len(s), MaxBytes)]) {
				if c == '\t' || c == '\n' {
					break
				}
				adv = utf8.RuneLen(c)
			}
			r.Box[nb] = Box{
				Nrune: i,
				Ptr:   s[i : i+adv],
				Width: r.Face.Dx(s[i : i+adv]),
			}
		case '\t':
			adv = 1
			r.Box[nb] = Box{
				Nrune:    -1,
				Ptr:      s[i : i+adv],
				Width:    r.minDx,
				Minwidth: r.minDx,
			}
		case '\n':
			adv = 1
			r.Box[nb] = Box{
				Nrune: -1,
				Ptr:   s[i : i+adv],
				Width: r.maxDx,
			}
			nl++
		}
	}
	r.Nchars -= int64(len(s))
	r.Nbox += nb
}

func (r *Run) Runescan(s []byte, ymax int) {
	r.Boxscan(s, ymax)
}
func (r *Run) Boxscan(s []byte, ymax int) {
	r.Nbox = 0
	r.Nchars = int64(len(s))
	i := 0
	nb := 0
	

	for nl := 0; nl <= ymax; nb++ {
		if i == len(s) {
			break
		}
		i++
		if nb == r.Nalloc {
			r.Nalloc += r.delta
			r.Box = append(r.Box, make([]Box, r.delta)...)
			// r.Grow(r.delta)
			if r.delta < 32768 {
				r.delta <<= 1
			}
		}
		if c := s[i-1]; c != '\t' && c != '\n' {
			span := len(s)
			if span > MaxBytes{
				span = MaxBytes
			}
			for _, c = range s[i:span] {
				if c == '\t' || c == '\n' {
					break
				}
				i++
			}
			r.Box[nb].Nrune = i
			r.Box[nb].Ptr = s[:i]
		} else if c == '\t' {
			r.Box[nb].Nrune = -1
			r.Box[nb].Ptr = s[:i]
			r.Box[nb].Width = r.minDx
			r.Box[nb].Minwidth = r.minDx
		} else {
			r.Box[nb].Nrune = -1
			r.Box[nb].Ptr = s[:i]
			r.Box[nb].Width = r.maxDx
			nl++
		}
		s = s[i:]
		i = 0
	}
	r.Nchars -= int64(len(s))
	for i := range r.Box[r.Nbox:r.Nbox+nb] {
		if b := &r.Box[i]; b.Nrune > 0 {
			b.Width = r.Face.Dx(b.Ptr)
		}
	}
	r.Nbox += nb
}

