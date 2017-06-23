package tag

import (
	"image"

	"github.com/as/frame/win"
)

type Invertable struct {
	*win.Win
	undo []func()
	do   []func()
	p    int
}

func (in *Invertable) Reset() {
	in.undo = nil
	in.do = nil
	in.p = 0
	return
}

func (in *Invertable) Insert(p []byte, at int64) int64 {
	d0, d1 := at, at+int64(len(p))
	in.do = in.do[:in.p]
	in.undo = in.undo[:in.p]
	in.do = append(in.do, func() { in.Win.Insert(p, at) })
	in.undo = append(in.undo, func() { in.Win.Delete(d0, d1) })
	in.p++
	in.do[in.p-1]()
	return 1
}
func (in *Invertable) Delete(q0, q1 int64) {
	data := append([]byte{}, in.Win.Bytes()[q0:q1]...)

	in.do = in.do[:in.p]
	in.undo = in.undo[:in.p]
	in.do = append(in.do, func() { in.Win.Delete(q0, q1) })
	in.undo = append(in.undo, func() { in.Win.Insert(data, q0) })
	in.p++
	in.do[in.p-1]()
	//in.Win.Delete(q0, q1)
}
func (in *Invertable) Undo() {
	l := len(in.undo)
	if l == 0 {
		return
	}
	in.p--
	in.undo[in.p]()
	if in.p < 0 {
		in.p = 0
	}
}
func (in *Invertable) Redo() {
	l := len(in.undo)
	if l == 0 || in.p >= l {
		return
	}
	in.p++
	in.do[in.p-1]()
}
func (in *Invertable) Select(q0, q1 int64) {
	w := in.Win
	w.Select(q0, q1)
	if Visible(w, q0, q1) {
		return
	}
	org := w.BackNL(q0, 2)
	w.SetOrigin(org, true)
	w.Frame.Select(q0-w.Org, q1-w.Org)
}
func (in *Invertable) Loc() image.Rectangle {
	if in == nil || in.Win == nil {
		return image.ZR
	}
	sp, size := in.Win.Sp, in.Win.Size()
	return image.Rectangle{sp, sp.Add(size)}
}
