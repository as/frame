package win

func (w *Win) fixEnd() {
	fr := w.Frame.Bounds()
	if pt := w.PointOf(w.Frame.Len()); pt.Y != fr.Max.Y {
		w.Paint(pt, fr.Max, w.Frame.Color.Palette.Back)
	}
}

func (w *Win) Fill() {
	if w.Frame.Full() {
		return
	}
	var rp [MsgSize]byte
	for !w.Frame.Full() {
		qep := w.org + w.Nchars
		n := max(0, min(w.Len()-qep, 2000))
		if n == 0 {
			break
		}
		m := copy(rp[:], w.Bytes()[qep:qep+n])
		nl := w.MaxLine() - w.Line()
		m = 0
		i := int64(0)
		for i < n {
			if rp[i] == '\n' {
				m++
				if m >= nl {
					i++
					break
				}
			}
			i++
		}
		w.Frame.Insert(rp[:i], w.Nchars)
		w.dirty = true
	}
}
