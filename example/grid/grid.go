package main

import (
	"bufio"
	"image"
	"os"
	"strings"

	"github.com/as/frame"
	"github.com/as/frame/tag"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

func mkfont(size int) frame.Font {
	return frame.NewTTF(gomono.TTF, size)
}

// Put
var (
	winSize = image.Pt(1900, 1000)
	pad     = image.Pt(25, 5)
	fsize   = 11
)

type Plane interface {
	Loc() image.Rectangle
}

// Put
func active(e mouse.Event, act Plane, list ...Plane) (x Plane) {
	if tag.Buttonsdown != 0 {
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

		N := 2
		dy := winSize.Y / N
		n := 0

		sp := image.Pt(0, n*dy)
		dp := image.Pt(winSize.X, dy)
		n++
		wn := tag.NewTag(src, wind, ft, sp, dp, pad, cols)
		sp = sp.Add(image.Pt(0, dy))

		wn2 := tag.NewTag(src, wind, ft, sp, dp, pad, cols)
		wn.Open(filename)
		wn2.Open(`C:\windows\system32\drivers\etc\hosts`)

		// lambda to paint only rectangles changed during a sweep of the mouse
		// Put
		act := wn.W
		actTag := wn

		go func() {
			sc := bufio.NewScanner(os.Stdin)
			for sc.Scan() {
				if x := sc.Text(); x == "u" || x == "r" {
					act.SendFirst(x)
					continue
				}
				act.SendFirst(tag.Cmdparse(sc.Text()))
			}
		}()
		for {
			// Put
			switch e := act.NextEvent().(type) {
			case mouse.Event:
				actTag = active(e, actTag, wn, wn2).(*tag.Tag)
				act = active(e, act, actTag.W, actTag.Wtag).(*tag.Invertable)
				actTag.Handle(act, e)
			case string, *tag.Command, tag.ScrollEvent, key.Event:
				actTag.Handle(act, e)
			case size.Event:
		winSize = image.Pt(e.WidthPx, e.HeightPx)
		N := 2
		dy := winSize.Y / N
		n := 0
		sp := image.Pt(0, n*dy)
		dp := image.Pt(winSize.X, dy)
		n++
		wn.Move(sp)
		wn.Resize(dp)
		sp = sp.Add(image.Pt(0, dy))
		wn2.Move(sp)
		wn2.Resize(dp)
				act.SendFirst(paint.Event{})
			case paint.Event:
				wn.Upload(wind)
				wn2.Upload(wind)
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
