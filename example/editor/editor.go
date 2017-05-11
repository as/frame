package main

import (
	"sync"
	//	"github.com/as/clip"
	//
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
)

var (
	wg      sync.WaitGroup
	winSize = image.Pt(2560, 1300)
	ClipBuf = make([]byte, 8192)
	Clip    *clip.Clip
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
	pad := image.Pt(4, 4)
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

var fsize = 15

func main() {
	driver.Main(func(src screen.Screen) {
		wind, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y})
		wind.Send(paint.Event{})
		focused := false
		focused = focused
		fsize := 15
		ft := mkfont(fsize)
		filename := "/dev/stdin"
		if len(os.Args) > 1 {
			filename = strings.Join(os.Args[1:], " ")
		}
		cols := frame.Acme
		pad := image.Pt(4, 4)

		// Make the main tag
		tagY := fsize * 2
		cols.Back = Cyan
		wmain := win.New(src, ft, wind, image.ZP, image.Pt(winSize.X, tagY), pad, cols)

		// Make tag
		wtag := win.New(src, ft, wind, image.Pt(0, tagY), image.Pt(winSize.X, tagY), pad, cols)

		// Make window
		cols.Back = Yellow
		w := win.New(src, ft, wind, image.Pt(0, tagY*2), winSize.Sub(image.Pt(0, tagY*2)), pad, cols)

		wtag.Insert(filename+"\tPut Del Exit", 0)
		wtag.Redraw()

		if len(os.Args) > 1 {
			s := string(readfile(filename))
			fmt.Printf("files size is %d\n", len(s))
			w.Insert(s, w.Q1)
		}

		// lambda to paint only rectangles changed during a sweep of the mouse

		act := w
		buttonsdown := 0x00000000
		for {
			switch e := act.NextEvent().(type) {
			case mouse.Event:
				pt := image.Pt(int(e.X), int(e.Y))
				if e.Direction == mouse.DirRelease {
					bt := uint(e.Button)
					buttonsdown &^= 1 << bt
				}
				if (e.Direction == mouse.DirNone || e.Direction == mouse.DirPress) && buttonsdown == 0 {
					apt := act.Sp.Add(pt)
					if !apt.In(image.Rectangle{act.Sp, act.Sp.Add(act.Size())}) {
						list := []*win.Win{wmain, wtag, w}
						for i, w := range list {
							r := image.Rectangle{w.Sp, w.Sp.Add(w.Size())}
							if apt.In(r) {
								if list[i] != act {
									fmt.Printf("active %d [because %s is in %s]\n", i, pt, r)
									act = list[i]
									break
								}
							}
						}
					}
				}
				if e.Button == mouse.ButtonWheelUp {
					act.FrameScroll(1)
					act.Send(paint.Event{})
					fmt.Println(e)
				}
				if e.Button == mouse.ButtonWheelDown {
					act.FrameScroll(-1)
					act.Send(paint.Event{})
					continue
				}
				if e.Direction == mouse.DirPress {
					bt := uint(e.Button)
					if e.Direction == mouse.DirPress {
						buttonsdown |= 1 << bt
					}
				}
				if e.Direction == mouse.DirPress {
					if e.Button == 2 && buttonsdown&(1<<1) != 0 {
						ClipBuf = ClipBuf[:cap(ClipBuf)]
						n, err := w.Read(ClipBuf)
						fmt.Printf("clip: read %d (err = %s)\n", n, err)
						if n > 0 {
							fmt.Printf("clip: %q\n", ClipBuf[:n])
						}
						ClipBuf = ClipBuf[:n]
						io.Copy(Clip, bytes.NewReader(toUTF16(ClipBuf)))
						act.Delete(act.Q0, act.Q1)
						act.Send(paint.Event{})
						fmt.Println("snarf")
					}
					if e.Button == 3 {
						if buttonsdown&(1<<1) != 0 {
							fmt.Println("paste")
							x := act.Q0
							act.Insert(string(toUTF16(ClipBuf)), act.Q0)
							act.Q0 = x
							act.Redraw()
							act.Send(paint.Event{})
						} else {

							//	act.P0, act.P1 = FindSpecial(NewReader(act.Frame), act.P0)
							//	act.Redraw()
							//	act.Send(paint.Event{})
						}
					}
				}
				if e.Button == 1 {
					switch e.Direction {
					case mouse.DirPress:
						act.Selectq = act.Org + act.IndexOf(pt)
						act.Select(pt, act, act.Upload)
						act.Q0 = act.Org + act.P0
						act.Q1 = act.Org + act.P1
						act.Redraw()
						act.Send(paint.Event{})
					case mouse.DirRelease:
					}
				}
				if (e.Button == 2 || e.Button == 3) && e.Direction == mouse.DirRelease {
					act.Selectq = act.Org + act.IndexOf(pt)
					fmt.Printf("index of %s is %d\n", pt, act.IndexOf(pt))
					act.Q0, act.Q1 = Expand(act.Bytes(), act.Selectq)
					act.P0, act.P1 = act.Q0-act.Org, act.Q1-act.Org
					x := act.Rdsel()
					if e.Button == 3 {
						act.Q0, act.Q1 = Next(act.Bytes(), act.Q1, x)
						act.P0, act.P1 = act.Q0-act.Org, act.Q1-act.Org
						moveMouse(act.PtOfChar(act.P0).Add(act.Sp))
					} else if e.Button == 2 {
						println(string(x))
						switch string(x) {
						case "Put":
							name, err := bufio.NewReader(bytes.NewReader(wtag.Bytes())).ReadString('\t')
							name = strings.TrimSpace(name)
							if err != nil {
								log.Printf("save: err: %s\n", err)
							}
							writefile(name, w.Bytes())
						default:
							println("Unknown command:", string(x))
						}
					}
					act.Redraw()
					act.Send(paint.Event{})
				}
			case key.Event:
				if e.Direction == key.DirRelease {
					continue
				}
				if e.Code == key.CodeEqualSign || e.Code == key.CodeHyphenMinus {
					if e.Modifiers == key.ModControl {
						if key.CodeHyphenMinus == e.Code {
							fsize--
						} else {
							fsize++
						}
						act.SetFont(mkfont(fsize))
						act.Send(paint.Event{})
						continue
					}
				}
				if e.Rune == '\r' {
					e.Rune = '\n'
				}
				if e.Code == key.CodeLeftArrow {
					if e.Modifiers&key.ModShift == 0 {
						act.Q1--
					}
					act.Q0--
					act.Redraw()
					act.Send(paint.Event{})
					continue
				}
				if e.Code == key.CodeUpArrow || e.Code == key.CodePageUp || e.Code == key.CodeDownArrow || e.Code == key.CodePageDown {
					n := 1
					if e.Code == key.CodeUpArrow || e.Code == key.CodePageUp {
						n = -1
					}
					if e.Code == key.CodePageUp || e.Code == key.CodePageDown {
						n *= 10
					}
					act.FrameScroll(n)
					act.Send(paint.Event{})
					continue
				}
				if e.Code == key.CodeRightArrow {
					if e.Modifiers&key.ModShift == 0 {
						act.Q0++
					}
					act.Q1++
					act.Redraw()
					act.Send(paint.Event{})
					continue
				}
				if e.Rune == '\x08' {
					si := act.Q0
					ei := act.Q1
					if si == ei && si != 0 {
						si--
					}
					fmt.Printf("si,ei=%d,%d\n", si, ei)
					act.Delete(si, ei)
					act.Send(paint.Event{})
					continue
				}
				if e.Rune == -1 {
					continue
				}
				if w.Q0 != w.Q1 {
					act.Delete(act.Q0, act.Q1)
				}
				fmt.Printf("outside %p: w.Nr=%d\n", act, act.Nr)
				act.Insert(string(e.Rune), act.Q1)
				act.Q0 = act.Q1
				act.Send(paint.Event{})
			case size.Event:
				pt := image.Pt(e.WidthPx, e.HeightPx)
				if pt.X < fsize || pt.Y < fsize {
					println("ignore daft size request:", pt.String())
					continue
				}
				winSize = pt
				wmain.Resize(image.Pt(winSize.X, tagY))
				wtag.Resize(image.Pt(winSize.X, tagY))
				//w.Resize(winSize.Sub(image.Pt(0, tagY*2)))
				act.Send(paint.Event{})
			case paint.Event:
				wind.Upload(wmain.Sp, wmain.Buffer(), wmain.Buffer().Bounds())
				wind.Upload(wtag.Sp, wtag.Buffer(), wtag.Buffer().Bounds())
				wind.Upload(w.Sp, w.Buffer(), w.Buffer().Bounds())
				wind.Publish()
				w.Flushcache()
				wtag.Flushcache()
				wmain.Flushcache()
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

func Expand(p []byte, i int64) (int64, int64) {
	return FindAlpha(p, i)
}

func Next(p []byte, i int64, sep []byte) (int64, int64) {
	x := int64(bytes.Index(p[i:], sep))
	if x == -1 {
		x = int64(bytes.Index(p[:i], sep))
	}
	if x != -1 {
		i += x
	}
	return i, i + int64(len(sep))
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
	Green  = image.NewUniform(color.RGBA{0, 255, 0, 255})
	Blue   = image.NewUniform(color.RGBA{0, 192, 192, 255})
	Cyan   = image.NewUniform(color.RGBA{234, 255, 255, 255})
	White  = image.NewUniform(color.RGBA{255, 255, 255, 255})
	Yellow = image.NewUniform(color.RGBA{255, 255, 224, 255})
	Gray   = image.NewUniform(color.RGBA{66, 66, 66, 255})
	Mauve  = image.NewUniform(color.RGBA{0x99, 0x99, 0xDD, 255})
)
