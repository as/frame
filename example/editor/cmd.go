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
	Select(w, q0, q1)
}

func Select(w *win.Win, q0, q1 int64) {
	if q0 == w.Q0 && q1 == w.Q1 {
		return
	}
	w.Q0, w.Q1 = q0, q1
	w.P0, w.P1 = q0-w.Org, q1-w.Org
}

func press(w, wtag *win.Win, e mouse.Event) {
	if e.Direction != mouse.DirPress {
		return
	}
	switch e.Button {
	case 1:
		w.Selectq = w.Org + w.IndexOf(pt(e))
		w.Select(pt(e), w, w.Upload)
		w.Q0, w.Q1 = w.Org+w.P0, w.Org+w.P1
		w.Redraw()
	}
}

func paste(w, wtag *win.Win, e mouse.Event) {
	x := w.Q0
	w.Insert(toUTF16(ClipBuf), w.Q0)
	w.Q0 = x
	w.Redraw()
	w.Send(paint.Event{})
}

func snarf(w, wtag *win.Win, e mouse.Event) {
	ClipBuf = ClipBuf[:cap(ClipBuf)]
	n, err := w.Read(ClipBuf)
	if n > 0 {
		fmt.Printf("clip: %q (%s)\n", ClipBuf[:n], err)
	}
	ClipBuf = ClipBuf[:n]
	io.Copy(Clip, bytes.NewReader(toUTF16(ClipBuf)))
	w.Delete(w.Q0, w.Q1)
	w.Redraw()
	w.Send(paint.Event{})
}

func release(w, wtag *win.Win, e mouse.Event) {
	if e.Direction != mouse.DirRelease {
		return
	}
	if e.Button == 1{
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
		q0, q1 := Next(w.Bytes(), w.Q1, w.Rdsel())
		Select(w, q0, q1)
		moveMouse(w.PtOfChar(w.P0).Add(w.Sp))
	}
	w.Redraw()
}
