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
func min(a,b int) int{
	if a<b{
		return a
	}
	return b
}

func (r *Run) Boxscan(s []byte, ymax int) {
	r.Nbox = 0
	r.Nchars = 0
	r.Nchars += int64(len(s))
	i := 0
	nb := 0
	for nl := 0; nl <= ymax; nb++ {
		if nb == r.Nalloc {
			r.Grow(r.delta)
			if r.delta < 32768 {
				r.delta *= 2
			}
		}
		if i == len(s) {
			break
		}
		i++
		c := s[i-1]
		switch c {
		default:
			for _, c = range s[i:min(len(s),MaxBytes)] {
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
				Ptr:      []byte("\t"),
				Width:    r.minDx,
				Minwidth: r.minDx,
			}
		case '\n':
			b := &r.Box[nb]
			b.Nrune=-1
			b.Ptr=b.Ptr[:1]
			b.Ptr[0]='\n'
			b.Width=r.maxDx
//			r.Box[nb] = Box{
//				Nrune: -1,
//				Ptr:   []byte("\n"),
//				Width: r.maxDx,
//			}
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
