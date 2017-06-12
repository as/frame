package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"io"
	"log"
	"strings"

	"github.com/as/frame/win"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
)

func pt(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

func Expand(w *win.Win, i int64) {
	p := w.Bytes()
	q0, q1 := FindAlpha(p, i)
	Sweep(w, q0, q1)
}

// Put	active
func Sweep(w *win.Win, q0, q1 int64) {
	if q0 == w.Q0 && q1 == w.Q1 {
		println("select is same")
		return
	}
	w.Q0, w.Q1 = q0, q1
	if !Visible(w, q0, q1) {
		org := w.BackNL(q0, 2)
		w.SetOrigin(org, true)
	}
	w.P0, w.P1 = q0-w.Org, q1-w.Org
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

// Put
func press(w, wtag *win.Win, e mouse.Event) {
	if e.Direction != mouse.DirPress {
		return
	}
	Buttonsdown |= 1 << uint(e.Button)
	fmt.Printf("press %d %08x\n", e.Button, Buttonsdown)
	switch e.Button {
	case 1:
		w.Selectq = w.Org + w.IndexOf(pt(e))
		w.Sweep = true
		w.Sweep(pt(e), w, w.Upload)
		w.Sweep = false
		w.Q0, w.Q1 = w.Org+w.P0, w.Org+w.P1
		w.Selectq = w.Q0
		w.Refresh()
	case 2:
		if down(1) {
			snarf(w, wtag, e)
		}
	case 3:
		if down(1) {
			paste(w, wtag, e)
		}
	}
}

func down(but uint) bool {
	return Buttonsdown&(1<<but) != 0
}

func paste(w, wtag *win.Win, e mouse.Event) {
	fmt.Println("paste")
	//x := w.Q0
	n, _ := Clip.Read(ClipBuf)
	s := fromUTF16(ClipBuf[:n])
	fmt.Printf("paste: %s\n", s)
	q0 := w.Q0
	w.Insert(s, q0)
	Sweep(w, q0, q0+int64(len(s)))
	//w.Q1 = x
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

func release(w, wtag *win.Win, e mouse.Event) {
	fmt.Printf("release %d %08x", e.Button, Buttonsdown)
	defer func() {
		Buttonsdown &^= 1 << uint(e.Button)
	}()
	if e.Direction != mouse.DirRelease {
		return
	}
	if e.Button == 1 || down(1) {
		return
	}
	w.Selectq = w.Org + w.IndexOf(pt(e))
	Expand(w, w.Selectq)
	switch e.Button {
	case 1:
		return
	case 2:
		x := w.Rdsel()
		switch string(x) {
		case "Put":
			if wtag == nil {
				panic("window has no tag")
			}
			name, err := bufio.NewReader(bytes.NewReader(wtag.Bytes())).ReadString('\t')
			name = strings.TrimSpace(name)
			if err != nil {
				log.Printf("save: err: %s\n", err)
			}
			writefile(name, w.Bytes())
		default:
			println("Unknown command:", string(x))
		}
	case 3:
		q0, q1 := Next(w.Bytes(), w.Q0, w.Q1)
		Sweep(w, q0, q1)
		moveMouse(w.PointOf(w.P0).Add(w.Sp))
	}
	w.Refresh()
}
