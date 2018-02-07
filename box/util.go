package box

// LenBox returns the lenth of box n. It panics if n is out of bounds.
func (r *Run) LenBox(n int) int {
	return r.Box[n].Len()
}

// WidthBox returns the width of box n. If the length of
// alt is different than the box, alt is measured and
// returned instead.
func (r *Run) WidthBox(n int, alt []byte) int {
	b := &(r.Box[n])
	if b.Nrune < 0 || len(alt) == b.Len() {
		return b.Width
	}
	return r.Face.Dx(alt)
}

// BoxBytes returns a trimmed pointer to box n's slice data
func (r *Run) BoxBytes(n int) (p []byte) {
	b := (&r.Box[n])
	return b.Ptr[:b.Len()]
}

// PlainBox returns true if and only if box n contains non-breaking
// characters
func (r *Run) PlainBox(n int) bool {
	return (&r.Box[n]).Nrune > 0
}
