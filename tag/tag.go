package tag

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
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

// Put
var (
	Buttonsdown = 0
	noselect    bool
	lastclickpt image.Point
)

type Tag struct {
	sp        image.Point
	Wtag      *Invertable
	W         *Invertable
	Scrolling bool
	scrolldy  int
	dirty bool
}

func (t *Tag) Dirty() bool{
	return t.dirty || t.Wtag.Dirty() || (t.W != nil && t.W.Dirty())
}

func (t *Tag) Mark() {
	t.dirty = true
}

func (t *Tag) Loc() image.Rectangle {
	r := t.Wtag.Loc()
	if t.W != nil {
		r.Max.Y += t.W.Loc().Dy()
	}
	return r
}

// Put
func NewTag(src screen.Screen, wind screen.Window, ft frame.Font,
	sp, size, pad image.Point, cols frame.Color) *Tag {

	// Make the main tag
	tagY := ft.Dy() * 2
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
	size = size.Sub(image.Pt(0, tagY))
	if size.Y < tagY {
		return &Tag{sp: sp, Wtag: wtag, W: nil}
	}
	// Make window
	cols.Back = Yellow
	w := &Invertable{
		win.New(src, ft, wind,
			sp,
			size,
			pad, cols,
		), nil, nil, 0,
	}
	return &Tag{sp: sp, Wtag: wtag, W: w}
}

func (t *Tag) Move(pt image.Point) {
	t.Wtag.Move(pt)
	if t.W == nil {
		return
	}
	pt.Y += t.Wtag.Loc().Dy()
	t.W.Move(pt)
}

func (t *Tag) Resize(pt image.Point) {
	dy := t.Wtag.Font.Dy() * 2
	if t.W != nil && dy < t.W.Font.Dy() {
		dy = t.W.Font.Dy() * 2
	}
	if pt.X < dy || pt.Y < dy {
		println("ignore daft size request:", pt.String())
		return
	}
	t.Wtag.Resize(image.Pt(pt.X, dy))
	pt.Y -= dy
	if t.W != nil {
		t.W.Resize(pt)
	}
}

func (t *Tag) split(path string) (name string, addr string) {
	name = path
	x := strings.Index(name, ":")
	if x == -1{
		return name, ""
	}
	if x == 0{
		if len(name) == 1{
			return ":", "" // This is invalid 
		}
		return "", name[1:]
	}
	if x == 1 && strings.IndexAny(name, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") == 0{
		if isdir(name[:2]) {
			n, a := t.split(name[2:])
			return name[:2]+n, a
		}
	}
	return name[:x], name[x+1:] 
}

func (t *Tag) Open(name string) {
	name, addr := t.split(name)
	w := t.W
	wtag := t.Wtag
	wtag.InsertString(name+"\tPut Del [Edit , ]", 0)
	wtag.Refresh()
	if w == nil {
		return
	}
	s := readfile(name)
	fmt.Printf("files size is %d\n", len(s))
	w.Insert(s, w.Q1)
	if addr != "" {
		w.Send(cmdparse("#0"))
		w.Send(cmdparse(addr))
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
		if e.Direction == key.DirRelease {
			return
		}
		fsize := act.Font.Size()
		if e.Modifiers == key.ModControl {
			if key.CodeHyphenMinus == e.Code {
				fsize -= 2
			} else {
				fsize += 2
			}
			act.SetFont(mkfont(fsize))
			t.Mark()
			// act.SendFirst(paint.Event{})
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
		t.Mark()
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
		t.Mark()
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
				q0 = acceptback(act.Bytes(), q0, AlphaNum)
			}
			act.Delete(q0, q1)
			act.Select(q0, q0)
		case '\x08':
			fallthrough
		default:
			act.Delete(q0, q1)
		}
		t.Mark()
		return
	}

	ch := []byte(string(e.Rune))
	if q1 != q0 {
		act.Delete(q0, q1)
		if region(act.Org, act.Org+act.Nchars, q0) == 0 || region(act.Org, act.Org+act.Nchars, q1) == 0{
			t.Mark()
		}
		q1 = q0
	}
	q1 += act.Insert(ch, q0)

	if region(act.Org, act.Org+act.Nchars, q0) == 0 || region(act.Org, act.Org+act.Nchars, q1) == 0{
		t.Mark()
	}
	q0 = q1
	act.Select(q0, q1)
	if region(act.Org, act.Org+act.Nchars, q0) == 0 || region(act.Org, act.Org+act.Nchars, q1) == 0{
		t.Mark()
	}
}



func (t *Tag) MouseIn(act *Invertable, e mouse.Event) {
	//	defer un(trace(db, "Tag.Mousein"))
	//	      func(){db.Trace(whatsdot(t.W))}()
	//	defer func(){db.Trace(whatsdot(t.W))}()

	pt := Pt(e)
	if e.Direction == mouse.DirRelease {
		t.Scrolling = false
	}
	
	if (e.Button != 0 && pt.In(act.Scrollr.Sub(act.Sp))) || t.Scrolling {
		//fmt.Printf("mouse.Event: %s\n", e)
		if e.Direction == mouse.DirRelease {
			return
		}
		if t.Scrolling {
			act.Clicksb(pt, 0)
			t.Mark()
		} else {
			if e.Button == 2 {
				t.Scrolling = true
			}
			act.Clicksb(pt, int(e.Button)-2)
			t.Mark()
		}
		
		return
	}
	switch e.Direction {
	case mouse.DirPress:
		lastclickpt = Pt(e)
		t.press(act, e)
		t.Mark()
	case mouse.DirRelease:
		lastclickpt = image.Pt(-5, -5)
		t.release(act, e)
		t.Mark()
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

type GetEvent struct {
	Path  string
	Addr  string
	IsDir bool
}

func isdir(path string) bool{
	fi, err := os.Stat(path)
	if err != nil {
		if err == os.ErrNotExist{
			return false
		}
		fmt.Println(err)
		return false
	}
	return fi.IsDir()
}
func isfile(path string) bool{
	_, err := os.Stat(path)
	return err == nil
}

func (t *Tag) Look(w *win.Win, q0, q1 int64) bool {
	if q0 == q1 {
		q1 = accept(w.Bytes(), q1, []byte(string(AlphaNum)+`\/.:`))
		q0 = acceptback(w.Bytes(), q0, []byte(string(AlphaNum)+`\/.:`))
	}
	name, addr := t.split(string(w.Bytes()[q0:q1]))
	fmt.Printf("name=%s addr=%s\n",name,addr)
	if name == "" && addr != "" {
		w.SendFirst(cmdparse(addr))
		return true
	}
	path := t.FileName()
	if !isdir(path){
		path = filepath.Dir(path)
	}
	path = filepath.Join(path, name)
	if isdir(path) {
		if addr != "" {
			// A directory with an address doesn't make sense
			// user probably refers to a file on another system
			// with the same name as the dir, so look should fail
			return false
		}
		w.SendFirst(GetEvent{Path: path, IsDir: true})
		return true
	} else if isfile(path){
		w.SendFirst(GetEvent{Path: path, Addr: addr})
		return true
	}
	fmt.Printf("look returns false")
	return false
}

func (t *Tag) FileName() string {
	if t.Wtag == nil {
		return ""
	}
	name, err := bufio.NewReader(bytes.NewReader(t.Wtag.Bytes())).ReadString('\t')
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
	t.W.Delete(0, t.W.Nr)
	t.W.Insert(readfile(name), 0)
	return nil
}

func (t *Tag) Put() (err error) {
	name := t.FileName()
	if name == "" {
		return fmt.Errorf("no file")
	}
	writefile(name, t.W.Bytes())
	return nil
}

// Put
func (t *Tag) Handle(act *Invertable, e interface{}) {
	switch e := e.(type) {
	case string:
		if e == "Redo" {
			act.Redo()
		} else if e == "Undo" {
			act.Undo()
		} else if e == "Put" {
			t.Put()
		} else if e == "Get" {
			t.Get()
		}
		t.Mark()
	case *Command:
		fmt.Printf("command %#v\n", e)
		if e == nil {
			break
		}
		if e.fn != nil {
			e.fn(t.W) // Always execute on body for now
		}
		t.Mark()
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
	if t.Dirty(){
		act.Send(paint.Event{})
	}
}

func (t *Tag) Upload(wind screen.Window) {
	if t.W != nil {
		wind.Upload(t.W.Sp, t.W.Buffer(), t.W.Buffer().Bounds())
		t.W.Flushcache()
		t.W.SetDirty(false)
	}
	wind.Upload(t.Wtag.Sp, t.Wtag.Buffer(), t.Wtag.Buffer().Bounds())
	t.Wtag.SetDirty(false)
	t.dirty = false
}
