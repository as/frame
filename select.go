package frame

// Select selects the region [p0:p1). The operation highlights
// the range of text under that region. If p0 = p1, a tick is
// drawn to indicate a null selection.
func (f *Frame) Select(p0, p1 int64) {
	pp0, pp1 := f.Dot()
	if pp1 <= p0 || p1 <= pp0 || p0 == p1 || pp1 == pp0 {
		f.redraw(f.point0(pp0), pp0, pp1, false)
		f.redraw(f.point0(p0), p0, p1, true)
	} else {
		if p0 < pp0 {
			f.redraw(f.point0(p0), p0, pp0, true)
		} else if p0 > pp0 {
			f.redraw(f.point0(pp0), pp0, p0, false)
		}
		if p1 > pp1 {
			f.redraw(f.point0(pp1), pp1, p1, true)
		} else if p1 < pp1 {
			f.redraw(f.point0(p1), p1, pp1, false)
		}
	}
	f.modified = true
	f.p0, f.p1 = p0, p1
}
