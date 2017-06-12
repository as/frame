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
	fmt.Printf("Delete %d:%d len=%d\n", q0, q1, len(in.Bytes()))
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
func (in *Invertable) Loc() image.Rectangle {
	sp, size := in.Win.Sp, in.Win.Size()
	return image.Rectangle{sp, sp.Add(size)}
}

var (
	wg      sync.WaitGroup
	winSize = image.Pt(1900, 1000)
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

type Mouse struct {
	down      int
	mask      int
	lastclick time.Time
} // Put

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

// Put
func active(e mouse.Event, act Plane, list ...Plane) (x interface{}) {
	defer func() {
		//fmt.Printf("active; %#v\n", x)
	}()
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

// Put
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
			case mouse.Event:
				act = active(e, act, wn.w, wn.wtag).(*Invertable)
				wn.Handle(act, e)
			case string, *command, scrollEvent, key.Event:
				wn.Handle(act, e)
			case size.Event:
				wn.Resize(image.Pt(e.WidthPx, e.HeightPx))
				act.SendFirst(paint.Event{})
			case paint.Event:
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
