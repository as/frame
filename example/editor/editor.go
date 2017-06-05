package main

import (
	"sync"
	//	"github.com/as/clip"
	//
	"bufio"
	"bytes"
	"fmt"
	"github.com/as/clip"
	"github.com/as/cursor"
	"github.com/as/frame"
	"github.com/as/frame/win"
	window "github.com/as/ms/win"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"image"
	"image/color"
	"image/draw"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// amink 5

// amink Put

type Action []func()

func (in *Invertable) Commit() {
}

type Scroller struct {
	*win.Win
	bar     draw.Image
	scrollr image.Rectangle
}

type scrollEvent struct {
	dy        int
	wind      *Invertable
	flushwith func(e interface{})
}

type Invertable struct {
	*win.Win
	undo []func()
	do   []func()
	p    int
}

type Plane interface {
	Loc() image.Rectangle
}
type Window interface {
	Buffer()
	Dot() (int64, int64)
	Insert([]byte, int64) int64
	Delete(q0, q1 int64) int64
	Select(int64, int64)
	Sp()
	Size()
	Send(e interface{})
	SendFirst(e interface{})
}

func (in *Invertable) Reset() {
	in.undo = nil
	in.do = nil
	in.p = 0
	return
}

func (in *Invertable) Insert(p []byte, at int64) int64 {
	d0, d1 := at, at+int64(len(p))
	fmt.Printf("Delete(%d, %d)\n", d0, d1)
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
	fmt.Printf("Insert(%q, %d)\n", data, q0)
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
	fmt.Printf("inverable: SetSelect: %s\n", q0, q1)
	w := in.Win
	w.Select(q0, q1)
	if Visible(w, q0, q1) {
		return
	}
	org := w.BackNL(q0, 2)
	w.SetOrigin(org, true)
	w.Frame.Select(q0-w.Org, q1-w.Org)
}

var (
	wg      sync.WaitGroup
	winSize = image.Pt(1000, 1000)
	pad     = image.Pt(25, 5)
	ClipBuf = make([]byte, 8192)
	Clip    *clip.Clip

	fsize    = 11
	ticking  = false
	scrolldy = 0
)

func readfile(s string) []byte {
	p, err := ioutil.ReadFile(s)
	if err != nil {
		log.Println(err)
	}
	return p
}
func writefile(s string, p []byte) {
	fd, err := os.Create(s)
	if err != nil {
		log.Println(err)
	}
	n, err := io.Copy(fd, bytes.NewReader(p))
	if err != nil {
		log.Fatalln(err)
	}
	println("wrote", n, "bytes")
}
func mkfont(size int) frame.Font {
	return frame.NewTTF(gomono.TTF, size)
}

func init() {
	var err error
	Clip, err = clip.New()
	if err != nil {
		panic(err)
	}
}
func moveMouse(pt image.Point) {
	cursor.MoveTo(window.ClientAbs().Min.Add(pt))
}

type Cell struct {
	*win.Win
	sp    image.Point
	size  image.Point
	in    chan string
	b     screen.Buffer
	wind  screen.Window
	src   screen.Screen
	dirty bool
}

func NewCell(src screen.Screen, wind screen.Window,
	sp, size image.Point, cols frame.Color) *Cell {
	b, _ := src.NewBuffer(size)
	w := win.New(src, mkfont(fsize), wind, sp, size, pad, cols)
	return &Cell{w, sp, size, make(chan string), b, wind, src, true}
}
func (c *Cell) Upload() {
	b := c.b
	c.wind.Upload(c.sp, b, b.Bounds())
	c.dirty = false
}
func (c *Cell) Mark() {
	c.dirty = true
}
func (c *Cell) Dirty() bool {
	return c.dirty
}

type Mouse struct {
	down      int
	mask      int
	lastclick time.Time
}

func (m *Mouse) clear() {
	m.down = 0
	m.mask = 0
}

// Put
var (
	buttonsdown = 0
	noselect    bool
	lastclickpt image.Point
)

func (in *Invertable) Loc() image.Rectangle {
	sp, size := in.Win.Sp, in.Win.Size()
	return image.Rectangle{sp, sp.Add(size)}
}

func active(e mouse.Event, act Plane, list ...Plane) interface{} {
	if buttonsdown != 0 {
		return act
	}
	pt := image.Pt(int(e.X), int(e.Y))
	if act != nil {
		pt = pt.Add(act.Loc().Min)
		list = append([]Plane{act}, list...)
	}
	for i, w := range list {
		r := w.Loc()
		if pt.In(r) {
			return list[i]
		}
	}
	return act
}

func wheeldir(e mouse.Event) int {
	if !e.Button.IsWheel() || (e.Button != mouse.ButtonWheelUp && e.Button != mouse.ButtonWheelDown) {
		return 0
	}
	if e.Button == mouse.ButtonWheelUp {
		return -1
	}
	return 1
}

// This fucking around will surely get us nowhere
//
//
var clock15hz = time.NewTicker(time.Millisecond * 5).C

func scroll2(act *Invertable, e mouse.Event) {
	dy := wheeldir(e)
	if dy == 0 {
		return
	}
	act.Send(scrollEvent{dy: dy, wind: act, flushwith: act.SendFirst})
}

func scroll(act *Invertable, e mouse.Event) {
	if e.Button == mouse.ButtonWheelUp || e.Button == mouse.ButtonWheelDown {
		dy := 1
		if e.Button == mouse.ButtonWheelUp {
			dy = -1
		}
		if !ticking {
			act := act
			//act.SendFirst(scrollEvent{dy: dy, wind: act, flushwith: act.SendFirst})
			ticking = true
			time.AfterFunc(time.Millisecond*15, func() { // Put
				ticking = false
				act.SendFirst(scrollEvent{dy: scrolldy, wind: act, flushwith: act.SendFirst})
				scrolldy = 0
			})
		} else {
			scrolldy += dy
		}
	}
}

func keyboard(act *Invertable, e key.Event) {
	if e.Direction == key.DirRelease {
		return
	}
	if e.Rune == '\r' {
		e.Rune = '\n'
	}
	q0, q1 := act.Dot()
	fmt.Println(q0, q1)
	switch e.Code {
	case key.CodeEqualSign, key.CodeHyphenMinus:
		if e.Modifiers == key.ModControl {
			if key.CodeHyphenMinus == e.Code {
				fsize--
			} else {
				fsize++
			}
			act.Reset()
			act.SetFont(mkfont(fsize))
			act.Fill()
			act.Send(paint.Event{})
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
	case '\x08', '\x17':
		if q0 == 0 && q1 == 0 {
			return
		}
		if e.Rune == '\x08' {
			q0--
		} else {
			if isany(act.Bytes()[q0], AlphaNum) {
				q0 = findback(act.Bytes(), q0, AlphaNum)
			}
		}
		act.Delete(q0, q1)
		act.Send(paint.Event{})
		return
	}
	if q0 != q1 {
		act.Delete(q0, q1)
	}
	act.Insert([]byte(string(e.Rune)), q1)
	q1++
	act.Select(q1, q1)
	act.Send(paint.Event{})
}

type Tag struct {
	sp        image.Point
	wtag      *Invertable
	w         *Invertable
	Scrolling bool
}

func NewTag(src screen.Screen, wind screen.Window, ft frame.Font,
	sp, size, pad image.Point, cols frame.Color) *Tag {

	// Make the main tag
	tagY := fsize * 2
	cols.Back = Cyan

	// Make tag
	wtag := &Invertable{
		win.New(src, ft, wind,
			image.Pt(0, tagY),
			image.Pt(size.X, tagY),
			pad, cols,
		), nil, nil, 0,
	}

	// Make window
	cols.Back = Yellow
	w := &Invertable{
		win.New(src, ft, wind,
			image.Pt(0, tagY*2),
			size.Sub(image.Pt(0, tagY*2)),
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
	tagY := t.wtag.Loc().Max.Y
	t.wtag.Resize(image.Pt(winSize.X, tagY))
	t.w.Resize(winSize.Sub(image.Pt(0, tagY*2)))
}

func (t *Tag) Open(filename string) {
	x := strings.Index(filename, ":")
	lineexpr := ""
	if x > 0 {
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

func (t *Tag) MouseIn(act *Invertable, e mouse.Event) {
	fmt.Printf("act.Scrollr: %s\n", act.Scrollr)
	fmt.Println(t.w.Dot())
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
		if !noselect && down(1) && ones(buttonsdown) == 1 {
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
			act.Selectq = act.Org+P0
			act.Refresh()
		}
	}
}

/*
func (t *Tag) CmdExec(cmd string){
	act := t.act
	w := t.w
	wtag := t.wtag
	if e == "r" {
		act.Redo()
	} else if e == "u" {
		act.Undo()
	} else if e == "Put" {
		Put(wtag.Win, w.Win)
	} else if e == "Get" {
		Get(wtag.Win, w.Win)
	}
	act.Send(paint.Event{})
}
*/
//func NewTag(scr screen.Screen, ft frame.Font, events screen.Window,
//       sp, size, pad image.Point, cols frame.Color) *Tag{
func main() {
	driver.Main(func(src screen.Screen) {
		wind, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y})
		wind.Send(paint.Event{})
		focused := false
		focused = focused
		ft := mkfont(fsize)
		filename := "/dev/stdin"
		if len(os.Args) > 1 {
			filename = strings.Join(os.Args[1:], " ")
		}
		cols := frame.Acme

		wn := NewTag(src, wind, ft, image.ZP, winSize, pad, cols)
		wn.Open(filename)

		w := wn.w
		wtag := wn.wtag

		// lambda to paint only rectangles changed during a sweep of the mouse
		// Put
		act := w
		go func() {
			sc := bufio.NewScanner(os.Stdin)
			for sc.Scan() {
				if x := sc.Text(); x == "u" || x == "r" {
					act.SendFirst(x)
					continue
				}
				act.SendFirst(cmdparse(sc.Text()))
			}
		}()
		for {
			// Put
			switch e := act.NextEvent().(type) {
			case string:
				if e == "r" {
					act.Redo()
				} else if e == "u" {
					act.Undo()
				} else if e == "Put" {
					wn.Put()
				} else if e == "Get" {
					wn.Get()
				}
				act.Send(paint.Event{})
			case *command:
				fmt.Printf("command %#v\n", e)
				if e == nil {
					panic("command is nil")
				}
				if e.fn != nil {
					e.fn(w) // Always execute on body for now
				}
				act.Send(paint.Event{})
			case scrollEvent:
				e.wind.FrameScroll(e.dy)
				e.flushwith(paint.Event{})
			case mouse.Event:
				act = active(e, act, w, wtag).(*Invertable)
				if e.Button.IsWheel() {
					scroll(act, e)
				} else {
					wn.MouseIn(act, e)
				}
			case key.Event:
				keyboard(act, e)
			case size.Event:
				wn.Resize(image.Pt(e.WidthPx, e.HeightPx))
				act.Send(paint.Event{})
			case paint.Event:
				if false {
					drawBorder(wtag.Buffer().RGBA(), wtag.Buffer().Bounds(), LtGray, image.ZP, 1)
					wind.Upload(wtag.Sp, wtag.Buffer(), wtag.Buffer().Bounds())
					w := w.Win
					pad := pad.Sub(image.Pt(5+3, 0))
					scrollr := image.Rect(w.Sp.X, 0, w.Sp.X+pad.X, w.Sp.Y+w.Buffer().Bounds().Max.Y)
					maxy := w.Buffer().Bounds().Max.Y
					rat0 := float64(w.Org) / float64(w.Nr)          // % scrolled
					rat1 := float64(w.Org+w.Nchars) / float64(w.Nr) // % covered by screen
					spY := int(float64(w.Sp.Y+maxy) * rat0)
					epY := int(float64(w.Sp.Y+maxy) * rat1)
					bar := image.Rect(w.Sp.X, spY, w.Sp.X+pad.X, epY)
					draw.Draw(w.Buffer().RGBA(), scrollr, X, image.ZP, draw.Src)
					draw.Draw(w.Buffer().RGBA(), bar, LtGray, image.ZP, draw.Src)
					drawBorder(w.Buffer().RGBA(), w.Buffer().Bounds(), LtGray, image.ZP, 1)
					drawBorder(w.Buffer().RGBA(), scrollr, LtGray, image.ZP, 1)
				}
				wn.Upload(wind)
				wind.Publish()
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
				// NT doesn't repaint the window if another window covers it
				if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOff {
					focused = false
				} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
					focused = true
				}
			}
		}
	})
}

func (t *Tag) Upload(wind screen.Window) {
	w := t.w
	wtag := t.wtag
	wind.Upload(w.Sp, w.Buffer(), w.Buffer().Bounds())
	wind.Upload(wtag.Sp, wtag.Buffer(), wtag.Buffer().Bounds())
	w.Flushcache()
	wtag.Flushcache()
}

func Next(p []byte, i, j int64) (q0 int64, q1 int64) {
	defer func(r0, r1 int64) {
		fmt.Printf("Next: [%d:%d]->[%d:%d]\n", r0, r1, q0, q1)
	}(i, j)
	x := p[i:j]
	q0 = int64(bytes.Index(p[j:], x))
	if q0 == -1 {
		println("a")
		q0 = int64(bytes.Index(p[:i], x))
		if q0 < 0 {
			println("b")
			return i, j
		}
	} else {
		println("c")
		q0 += j
	}
	q1 = q0 + int64(len(x))
	println("d")
	return q0, q1
}

func drawBorder(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, thick int) {
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+thick), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-thick, r.Max.X, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X+thick, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Max.X-thick, r.Min.Y, r.Max.X, r.Max.Y), src, sp, draw.Src)
}

func toUTF16(p []byte) (q []byte) {
	i := 0
	q = make([]byte, len(p)*2)
	for j := 0; j < len(p); j++ {
		q[i] = p[j]
		i += 2
	}
	return q
}
func fromUTF16(p []byte) (q []byte) {
	j := 0
	q = make([]byte, len(p)/2)
	for i := 0; i < len(q); i++ {
		q[i] = p[j]
		j += 2
	}
	return q
}

var (
	Red    = image.NewUniform(color.RGBA{255, 0, 0, 255})
	Green  = image.NewUniform(color.RGBA{255, 255, 192, 25})
	Blue   = image.NewUniform(color.RGBA{0, 192, 192, 255})
	Cyan   = image.NewUniform(color.RGBA{234, 255, 255, 255})
	White  = image.NewUniform(color.RGBA{255, 255, 255, 255})
	Yellow = image.NewUniform(color.RGBA{255, 255, 224, 255})
	X      = image.NewUniform(color.RGBA{255 - 32, 255 - 32, 224 - 32, 255})

	LtGray = image.NewUniform(color.RGBA{66 * 2, 66 * 2, 66*2 + 35, 255})
	Gray   = image.NewUniform(color.RGBA{66, 66, 66, 255})
	Mauve  = image.NewUniform(color.RGBA{0x99, 0x99, 0xDD, 255})
)
