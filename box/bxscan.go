package box

// MaxBytes is the largest capacity of bytes in a box
var MaxBytes = 256+3

func (r *Run) ensure(nb int) {
	if nb == r.Nalloc{
		r.Grow(r.delta)
		if r.delta < 16384{
			r.delta *= 2
		}
	}
}

func (r *Run) Bxscan(s []byte, ymax int){
	var nl int

	r.delta = 32
	r.ss.Reset(s)

	for nb := 0; r.ss.Len() > 0 && nl <= ymax; nb++ {
		r.ensure(nb)
		c, _ := r.ss.ReadByte()
		if special(c) {
			nl += r.specialbox(nb, c, r.minDx, r.maxDx)
		} else {
			r.ss.UnreadByte()
			nl += r.plainbox(nb)
		}
		r.Nbox++
	}
}

func special(c byte) bool {
	return c == '\t' || c == '\n'
}

func (r *Run) specialbox(nb int, c byte, min, max int) (nl int) {
	b := &r.Box[nb]
	if c == '\n' {
		b.Minwidth = 0
		nl++
	} else {
		b.Minwidth = min
	}
	b.BC = c
	b.Ptr = []byte{c}
	b.Nrune = -1
	b.Width = max
	r.Nchars++
	return
}

func (r *Run) plainbox(nb int) (nl int){
	r.ww.Reset()

	nr := 0
	for ; r.ss.Len() > 0; nr++ {
		c, _ := r.ss.ReadByte()
		rw := 1
		if special(c) || nr+rw >= MaxBytes{
			r.ss.UnreadByte()
			break
		}
		r.ww.WriteByte(c)
	}

	b := &r.Box[nb]
	b.Ptr = make([]byte, nr)
	copy(b.Ptr, r.ww.Bytes())
	b.Width = r.MeasureBytes(b.Ptr)
	b.Nrune = nr
	r.Nchars += int64(nr)

	return 0
}
