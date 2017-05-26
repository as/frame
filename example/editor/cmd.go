package main
// Edit ,x,pt\(,x,pt,c,Pt,
import (
	"bufio"
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
	Select(w, q0, q1)
}
// Put	active
func Select(w *win.Win, q0, q1 int64) {
	if q0 == w.Q0 && q1 == w.Q1 {
		println("select is same")
		return
	}
	w.Q0, w.Q1 = q0, q1
	if !Visible(w, q0, q1){
		org := w.BackNL(q0, 2)
		w.SetOrigin(org, true)
	}
	w.P0, w.P1 = q0-w.Org, q1-w.Org
}

func Visible(w *win.Win, q0, q1 int64) bool{
	if q0 < w.Org {
		return false
	}
	if q1 > w.Org+w.Frame.Nchars{
		return false
	}
	return true
}

func press(w, wtag *win.Win, e mouse.Event) {
	fmt.Printf("press (enter) %d %08x\n", e.Button, buttonsdown)
	defer func(){fmt.Printf("press (leave) %d %08x\n", e.Button, buttonsdown)}()
	if e.Direction != mouse.DirPress {
		return
	}
	defer func(){
		buttonsdown |= 1<<uint(e.Button)
	}()
	switch e.Button {
	case 1:
		if down(1){
			println("impossible condition")
		}
		t := time.Now()
		dt := t.Sub(w.Lastclick)
		w.Selectq = w.Org + w.IndexOf(Pt(e))
		w.Lastclick = t
		if dt < 250*time.Millisecond && w.Q0 == w.Q1 && w.Q0 == w.Selectq{
			Expand(w, w.Selectq)
		} else {
			w.Sweep = true
			w.Select(Pt(e), w, w.Upload)
			w.Sweep = false
		}
		w.Q0, w.Q1 = w.Org+w.P0, w.Org+w.P1
		w.Selectq = w.Q0
		w.Redraw()
	case 2:
		if down(1) {
			noselect=true
			snarf(w, wtag, e)
		}
	case 3:
		if down(1){
			noselect=true
			paste(w, wtag, e)
		}
	}
}

func ones(n int) (c int){
	for n != 0{
		n &= (n-1)
		c++
	}
	return c
}

func down(but uint) bool{
	return buttonsdown & (1<<but) != 0
}

func paste(w, wtag *win.Win, e mouse.Event) {
	fmt.Println("paste")
	//x := w.Q0
	n, _ := Clip.Read(ClipBuf)
	s := fromUTF16(ClipBuf[:n])
	fmt.Printf("paste: %s\n", s)
	q0 := w.Q0
	w.Insert(s, q0)
	Select(w, q0, q0+int64(len(s)))
	//w.Q1 = x
	w.Redraw()
	w.Send(paint.Event{})
}

func snarf(w, wtag *win.Win, e mouse.Event) {
	fmt.Println("snarf")
	n := copy(ClipBuf, toUTF16([]byte(w.Rdsel())))
	io.Copy(Clip, bytes.NewReader(ClipBuf[:n]))
	w.Delete(w.Q0, w.Q1)
	w.Redraw()
	w.Send(paint.Event{})
}

func region(q0, q1, x int64) int{
	if x < q0{
		return -1
	}
	if x > q1{
		return 1
	}
	return 0
}

func whatsdown(){
	fmt.Printf("down: %08x\n", buttonsdown)
}

func release(w, wtag *win.Win, e mouse.Event) {
	 func(){fmt.Printf("release (enter) %d %08x\n", e.Button, buttonsdown)}()
	defer func(){fmt.Printf("release (leave) %d %08x\n", e.Button, buttonsdown)}()
	defer func(){
			buttonsdown &^= 1<<uint(e.Button)
			if buttonsdown == 0{
				noselect=false
			}
	}()
	if e.Direction != mouse.DirRelease {
		return
	}
	if e.Button == 1 || down(1) {
		return
	}

	whatsdown()

	w.Selectq = w.Org + w.IndexOf(Pt(e))
	if region(w.Q0, w.Q1, w.Selectq) != 0{
		Expand(w, w.Selectq)
	}
	switch e.Button {
	case 1:
		return
	case 2:
		x := strings.TrimSpace(string(w.Rdsel()))
		switch {
		case strings.HasPrefix(x, "Edit"):
			if x == "Edit"{
				log.Printf("Edit: empty command\n")
				break
			}
			w.SendFirst(cmdparse(x[4:]))
		default:
			w.SendFirst(x)
		}
	case 3:
		//w.SendFirst(cmdparse(fmt.Sprintf("/%s/", w.Rdsel())))
		//break
			q0, q1 := Next(w.Bytes(), w.Q0, w.Q1)
			Select(w, q0, q1)
			moveMouse(w.PtOfChar(w.P0).Add(w.Sp))
	}
	w.Redraw()
}

func filename(wtag *win.Win) (string, error){
	if wtag == nil {
		return "", fmt.Errorf("window has no tag")
	}
	name, err := bufio.NewReader(bytes.NewReader(wtag.Bytes())).ReadString('\t')
	if err != nil{
		return "", err
	}
	return strings.TrimSpace(name), nil
}

func Get(wtag, w *win.Win) (err error){
	name, err := filename(wtag)
	if err != nil{
		return err
	}
	w.Delete(0, w.Nr)
	w.Insert(readfile(name), 0)
	return nil
}

func Put(wtag, w *win.Win) (err error){
	log.Println("Put")
	name, err := filename(wtag)
	if err != nil {
		return err
	}
	writefile(name, w.Bytes())
	return nil
}
