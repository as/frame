package tag

// Edit ,x,pt\(,x,pt,c,Pt,
import (
	"bytes"
	"fmt"
	"image"
	"io"

	"github.com/as/frame/win"
	"github.com/as/text"
	"golang.org/x/mobile/event/mouse"
)

var chorded = false

func Pt(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

func Visible(w *win.Win, q0, q1 int64) bool {
	if q0 < w.Origin() {
		return false
	}
	if q1 > w.Origin()+w.Frame.Nchars {
		return false
	}
	return true
}

func Paste(w text.Editor, e mouse.Event) {
	n, _ := Clip.Read(ClipBuf)
	s := fromUTF16(ClipBuf[:n])
	q0, q1 := w.Dot()
	if q0 != q1 {
		w.Delete(q0, q1)
		q1 = q0
	}
	w.Insert(s, q0)
	w.Select(q0, q0+int64(len(s)))
}

func Rdsel(w text.Editor) string {
	q0, q1 := w.Dot()
	return string(w.Bytes()[q0:q1])
}

func Snarf(w text.Editor, e mouse.Event) {
	n := copy(ClipBuf, toUTF16([]byte(Rdsel(w))))
	io.Copy(Clip, bytes.NewReader(ClipBuf[:n]))
	q0, q1 := w.Dot()
	w.Delete(q0, q1)
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
