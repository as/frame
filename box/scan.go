package box

//import "log"

func (r *Run) ensure(nb int) {
	if nb == r.Nalloc {
		r.Grow(r.delta)
		if r.delta < 16384 {
			r.delta *= 2
		}
	}
}

func (r *Run) Boxscan(s string, ymax int) {
	//r.delta = 32
	r.br.Reset(s)

	for nb, nl := 0, 0; nl <= ymax; nb++ {
		if _, _, err := r.br.Next(); err != nil {
			break
		}
		r.ensure(nb)
		if special(r.br.Last()[0]) {
			nl += r.specialbox(nb, r.minDx, r.maxDx)
		} else {
			nl += r.plainbox(nb)
		}
		r.Nbox++
	}
}

func special(c byte) bool {
	return c == '\t' || c == '\n'
}

func (r *Run) specialbox(nb int, min, max int) (nl int) {
	c := r.br.Last()[0]
	b := &r.Box[nb]
	if c == '\n' {
		b.Minwidth = 0
		nl++
	} else {
		b.Minwidth = min
	}
	b.Ptr = string(c)
	b.Nrune = -1
	b.Width = max

	if c == '\t' {
		b.Width = min
	}

	r.Nchars++
	r.br.Advance()
	return
}

func (r *Run) plainbox(nb int) (nl int) {
	for {
		_, _, err := r.br.Next()
		if err != nil {
			break
		}
		if special(r.br.Last()[0]) || r.br.Len() >= MaxBytes {
			r.br.Unread()
			break
		}
	}
	b := &r.Box[nb]
	b.Ptr = r.br.Bytes()
	b.Width = r.br.Width()
	b.Nrune = r.br.Len()
	r.Nchars += int64(r.br.Len())
	r.br.Advance()
	return 0
}
