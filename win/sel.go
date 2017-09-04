package win

import "github.com/as/text"

// Select selects the range [q0:q1] inclusive
func (w *Win) Select(q0, q1 int64) {
	if q0 > q1 {
		q0, q1 = q1, q0
	}
	q00, q11 := w.Dot()
	w.Editor.Select(q0, q1)
	reg := text.Region3(q0, w.org-1, w.org+w.Frame.Len())
	if q00 == q0 && q11 == q1 {
		//return
	}
	w.dirty = true
	p0, p1 := q0-w.org, q1-w.org
	pp0, pp1 := w.Frame.Dot()
	if pp1 <= p0 || p1 <= pp0 || p0 == p1 || pp1 == pp0 {
		w.Redraw(w.PointOf(pp0), pp0, pp1, false)
		w.Redraw(w.PointOf(p0), p0, p1, true)
	} else {
		if p0 < pp0 {
			w.Redraw(w.PointOf(p0), p0, pp0, true)
		} else if p0 > pp0 {
			w.Redraw(w.PointOf(pp0), pp0, p0, false)
		}
		if pp1 < p1 {
			w.Redraw(w.PointOf(pp1), pp1, p1, true)
		} else if pp1 > p1 {
			w.Redraw(w.PointOf(p1), p1, pp1, false)
		}
	}
	w.Frame.Select(p0, p1)
	if q0 == q1 && reg != 0 {
		w.Untick()
	}
}
