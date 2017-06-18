package tag

// Edit ,x,pt\(,x,pt,c,Pt,
import (
	"bytes"
	"fmt"
	"image"
	"io"
	"log"
	"strings"
	"time"

	"github.com/as/frame/win"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
)

var chorded = false

func Pt(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

func Expand(w *win.Win, i int64) {
	p := w.Bytes()
	q0, q1 := FindAlpha(p, i)
	Sweep(w, q0, q1)
}

// Sweep selects q0:q1 and moves the origin into view
func Sweep(w *win.Win, q0, q1 int64) {
	w.Select(q0, q1)
	if !Visible(w, q0, q1) {
		org := w.BackNL(q0, 2)
		w.SetOrigin(org, true)
	}
	return
}

func Visible(w *win.Win, q0, q1 int64) bool {
	if q0 < w.Org {
		return false
	}
	if q1 > w.Org+w.Frame.Nchars {
		return false
	}
	return true
}

//
// press
// set selectq
// 	frame.Sweep |
// w.FrameScroll <--
//    \--> w.Select
func press(w, wtag *win.Win, e mouse.Event) {
	//	defer un(trace(db, "press(Tag)"))	// Debug
	//	      func(){db.Trace(whatsdot(w))}()	// Debug
	//	defer func(){db.Trace(whatsdot(w))}()	// Debug
	if e.Direction != mouse.DirPress {
		return
	}
	defer func() {
		Buttonsdown |= 1 << uint(e.Button)
	}()
	switch e.Button {
	case 1:
		if down(1) {
			println("impossible condition")
		}
		t := time.Now()
		dt := t.Sub(w.Lastclick)
		org := w.Org
		w.Selectq = w.Org + w.IndexOf(Pt(e))
		w.Lastclick = t
		if dt < 250*time.Millisecond && w.Q0 == w.Q1 && w.Q0 == w.Selectq {
			Expand(w, w.Selectq)
			P0, P1 := w.Frame.Dot()
			w.Select(w.Org+P0, w.Org+P1)
		} else {
			w.Sweeping = true
			w.Sweep(Pt(e), w, w.Upload)
			w.Sweeping = false
		}
		P0, P1 := w.Frame.Dot()
		q0, q1 := w.Org+P0, w.Org+P1
		if org > w.Org {
			q1 = w.Selectq
		} else if org < w.Org {
			q0 = w.Selectq
		}
		w.Select(q0, q1)
		//w.Selectq = w.Q0
		w.Refresh()
	case 2:
		if down(1) {
			noselect = true
			snarf(w, wtag, e)
		}
	case 3:
		if down(1) {
			noselect = true
			paste(w, wtag, e)
		}
	}
}

func release(w, wtag *win.Win, e mouse.Event) {
	func() { fmt.Printf("release (enter) %d %08x\n", e.Button, Buttonsdown) }()
	defer func() { fmt.Printf("release (leave) %d %08x\n", e.Button, Buttonsdown) }()
	defer func() {
		Buttonsdown &^= 1 << uint(e.Button)
		if Buttonsdown == 0 {
			noselect = false
		}
	}()
	if e.Direction != mouse.DirRelease || e.Button == 1 || down(1) {
		return
	}
	w.Selectq = w.Org + w.IndexOf(Pt(e))
	if region(w.Q0, w.Q1, w.Selectq) != 0 {
		Expand(w, w.Selectq)
	}
	P := strings.HasPrefix
	switch e.Button {
	case 1:
		return
	case 2:
		x := strings.TrimSpace(string(w.Rdsel()))
		switch {
		case P(x, "Edit"):
			if x == "Edit" {
				log.Printf("Edit: empty command\n")
				break
			}
			w.SendFirst(cmdparse(x[4:]))
		case P(x, "|"), P(x, "<"), P(x, ">"):
			w.SendFirst(cmdparse(x))
		case P(x, ":"):
			w.SendFirst(cmdparse(x[1:]))
		default:
			w.SendFirst(x)
		}
	case 3:
		//w.SendFirst(cmdparse(fmt.Sprintf("/%s/", w.Rdsel())))
		//break
		q0, q1 := Next(w.Bytes(), w.Q0, w.Q1)
		Sweep(w, q0, q1)
		P0, _ := w.Frame.Dot()
		moveMouse(w.PointOf(P0).Add(w.Sp))
	}
	w.Refresh()
}

// Put
func ones(n int) (c int) {
	for n != 0 {
		n &= (n - 1)
		c++
	}
	return c
}

func down(but uint) bool {
	return Buttonsdown&(1<<but) != 0
}

func paste(w, wtag *win.Win, e mouse.Event) {
	n, _ := Clip.Read(ClipBuf)
	s := fromUTF16(ClipBuf[:n])
	q0, q1 := w.Dot()
	if q0 != q1 {
		w.Delete(q0, q1)
		q1 = q0
	}
	w.Insert(s, q0)
	Sweep(w, q0, q0+int64(len(s)))
	w.Refresh()
	w.Send(paint.Event{})
}

func snarf(w, wtag *win.Win, e mouse.Event) {
	fmt.Println("snarf")
	n := copy(ClipBuf, toUTF16([]byte(w.Rdsel())))
	io.Copy(Clip, bytes.NewReader(ClipBuf[:n]))
	w.Delete(w.Q0, w.Q1)
	w.Refresh()
	w.Send(paint.Event{})
}

func region(q0, q1, x int64) int {
	if x < q0 {
		return -1
	}
	if x > q1 {
		return 1
	}
	return 0
}

func whatsdown() {
	fmt.Printf("down: %08x\n", Buttonsdown)
}
