package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/as/frame"
	"github.com/as/frame/win"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
)

var db = win.Db
var un = win.Un
var trace = win.Trace

type doter interface {
	Dot() (int64, int64)
}

func whatsdot(d doter) string {
	q0, q1 := d.Dot()
	return fmt.Sprintf("Dot: [%d:%d]", q0, q1)
}

type Tag struct {
	sp        image.Point
	wtag      *Invertable
	w         *Invertable
	Scrolling bool
}

// Put
func NewTag(src screen.Screen, wind screen.Window, ft font.Font,
	sp, size, pad image.Point, cols frame.Color) *Tag {

	// Make the main tag
	tagY := fsize * 2
	cols.Back = Cyan

	// Make tag
	wtag := &Invertable{
		win.New(src, ft, wind,
			sp,
			image.Pt(size.X, tagY),
			pad, cols,
		), nil, nil, 0,
	}

	sp = sp.Add(image.Pt(0, tagY))

	// Make window
	cols.Back = Yellow
	w := &Invertable{
		win.New(src, ft, wind,
			sp,
			size.Sub(image.Pt(0, sp.Y)),
			pad, cols),
		nil, nil, 0,
	}
	return &Tag{sp: sp, wtag: wtag, w: w}
}

func (t *Tag) Resize(pt image.Point) {
	if pt.X < fsize || pt.Y < fsize {
		println("ignore daft size request:", pt.String())
		return
	}
	winSize = pt
	tagY := t.wtag.Loc().Dy()
	t.wtag.Resize(image.Pt(winSize.X, tagY))
	t.w.Resize(winSize.Sub(image.Pt(0, tagY)))
}

func (t *Tag) Open(filename string) {
	x := strings.Index(filename, ":")
	lineexpr := ""
	if x > 1 {
		lineexpr = filename[x+1:]
		filename = filename[:x]
	}

	w := t.w
	wtag := t.wtag

	wtag.InsertString(filename+"\tPut Del Exit", 0)
	wtag.Refresh()
	if len(os.Args) > 1 {
		s := readfile(filename)
		fmt.Printf("files size is %d\n", len(s))
		w.Insert(s, w.Q1)
		if lineexpr != "" {
			w.Send(cmdparse("#0"))
			w.Send(cmdparse(lineexpr))
		}
	}
}
func (t *Tag) Kbdin(act *Invertable, e key.Event) {
	if e.Direction == key.DirRelease {
		return
	}
	if e.Rune == '\r' {
		e.Rune = '\n'
	}
	q0, q1 := act.Dot()
	switch e.Code {
	case key.CodeEqualSign, key.CodeHyphenMinus:
		if e.Modifiers == key.ModControl {
			if key.CodeHyphenMinus == e.Code {
				fsize--
			} else {
				fsize++
			}
			act.SetFont(mkfont(fsize))
			act.SendFirst(paint.Event{})
			return
		}
	case key.CodeUpArrow, key.CodePageUp, key.CodeDownArrow, key.CodePageDown:
		org := act.Org
		n := act.MaxLine() / 7
		if e.Code == key.CodePageUp || e.Code == key.CodePageDown {
			n *= 10
		}
		if e.Code == key.CodeUpArrow || e.Code == key.CodePageUp {
			org = act.BackNL(org, n)
		}
		if e.Code == key.CodeDownArrow || e.Code == key.CodePageDown {
			r := act.Bounds()
			org += act.IndexOf(image.Pt(r.Min.X, r.Min.Y+n*act.Frame.Dy()))
		}
		act.SetOrigin(org, true)
		act.Send(paint.Event{})
		return
	case key.CodeLeftArrow, key.CodeRightArrow:
		if e.Code == key.CodeLeftArrow {
			if e.Modifiers&key.ModShift == 0 {
				q1--
			}
			q0--
		} else {
			if e.Modifiers&key.ModShift == 0 {
				q0++
			}
			q1++
		}
		act.Select(q0, q1)
		act.Send(paint.Event{})
		return
	}
	switch e.Rune {
	case -1:
		return

	case '\x01', '\x05', '\x08', '\x15', '\x17':
		if q0 == 0 && q1 == 0 {
			return
		}
		if q0 == q1 && q0 != 0 {
			q0--
		}
		switch e.Rune {
		case '\x15', '\x01': // ^U, ^A
			p := act.Bytes()
			if q0 < int64(len(p))-1 {
				q0++
			}
			n0, n1 := findlinerev(act.Bytes(), q0, 0)
			if e.Rune == '\x15' {
				act.Delete(n0, n1)
			}
			act.Select(n0, n0)
		case '\x05': // ^E
			_, n1 := findline3(act.Bytes(), q1, 1)
			if n1 > 0 {
				n1--
			}
			act.Select(n1, n1)
		case '\x17':
			if isany(act.Bytes()[q0], AlphaNum) {
				q0 = findback(act.Bytes(), q0, AlphaNum)
			}
			act.Delete(q0, q1)
			act.Select(q0, q0)
		case '\x08':
			fallthrough
		default:
			act.Delete(q0, q1)
		}
		act.Send(paint.Event{})
		return
	}
	if q0 != q1 {
		act.Delete(q0, q1)
	}
	q0 += act.Insert([]byte(string(e.Rune)), q0)
	q1 = q0
	act.Select(q0, q1)
	act.Send(paint.Event{})
}

func (t *Tag) MouseIn(act *Invertable, e mouse.Event) {
	//	defer un(trace(db, "Tag.Mousein"))
	//	      func(){db.Trace(whatsdot(t.w))}()
	//	defer func(){db.Trace(whatsdot(t.w))}()

	pt := Pt(e)
	if e.Direction == mouse.DirRelease {
		t.Scrolling = false
	}
	if pt.In(act.Scrollr.Sub(act.Sp)) || t.Scrolling {
		fmt.Printf("mouse.Event: %s\n", e)
		if e.Direction == mouse.DirRelease {
			return
		}
		if t.Scrolling {
			act.Clicksb(pt, 0)
		} else {
			if e.Button == 2 {
				t.Scrolling = true
			}
			act.Clicksb(pt, int(e.Button)-2)
		}
		act.Send(paint.Event{})
		return
	}
	switch e.Direction {
	case mouse.DirPress:
		lastclickpt = Pt(e)
		press(act.Win, t.wtag.Win, e)
		act.Send(paint.Event{})
	case mouse.DirRelease:
		lastclickpt = image.Pt(-5, -5)
		release(act.Win, t.wtag.Win, e)
		act.Send(paint.Event{})
	case mouse.DirNone:
		if !noselect && down(1) && ones(Buttonsdown) == 1 {
			r := image.Rect(0, 0, 5, 5).Add(lastclickpt)
			pt := Pt(e)
			if pt.In(r) {
				return
			}
			// Double click happened so select function
			// never fired.
			act.Sweeping = true
			act.Sweep(lastclickpt, act, act.Upload)
			act.Sweeping = false
			P0, P1 := act.Frame.Dot()
			act.Select(act.Org+P0, act.Org+P1)
			act.Selectq = act.Org + P0
			act.Refresh()
		}
	}
}

func (t *Tag) FileName() string {
	if t.wtag == nil {
		return ""
	}
	name, err := bufio.NewReader(bytes.NewReader(t.wtag.Bytes())).ReadString('\t')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(name)
}

func (t *Tag) Get() (err error) {
	name := t.FileName()
	if name == "" {
		return fmt.Errorf("no file")
	}
	t.w.Delete(0, t.w.Nr)
	t.w.Insert(readfile(name), 0)
	return nil
}

func (t *Tag) Put() (err error) {
	name := t.FileName()
	if name == "" {
		return fmt.Errorf("no file")
	}
	writefile(name, t.w.Bytes())
	return nil
}

// Put
func (t *Tag) Handle(act *Invertable, e interface{}) {
	switch e := e.(type) {
	case string:
		if e == "r" {
			act.Redo()
		} else if e == "u" {
			act.Undo()
		} else if e == "Put" {
			t.Put()
		} else if e == "Get" {
			t.Get()
		}
		act.Send(paint.Event{})
	case *command:
		fmt.Printf("command %#v\n", e)
		if e == nil {
			panic("command is nil")
		}
		if e.fn != nil {
			e.fn(t.w) // Always execute on body for now
		}
		act.Send(paint.Event{})
	case ScrollEvent:
		e.wind.FrameScroll(e.dy)
		e.flushwith(paint.Event{})
	case mouse.Event:
		if e.Button.IsWheel() {
			scroll(act, e)
		} else {
			t.MouseIn(act, e)
		}
	case key.Event:
		t.Kbdin(act, e)
	}
}

func (t *Tag) Upload(wind screen.Window) {

	w := t.w
	wtag := t.wtag
	wind.Upload(w.Sp, w.Buffer(), w.Buffer().Bounds())
	wind.Upload(wtag.Sp, wtag.Buffer(), wtag.Buffer().Bounds())
}
