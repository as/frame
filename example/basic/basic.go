package main

import (
	"image"
	"image/draw"
	"sync"

	"github.com/as/frame"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var wg sync.WaitGroup
var winSize = image.Pt(1024, 768)

func main() {
	driver.Main(func(src screen.Screen) {
		var dirty = true
		wind, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y, "basic"})
		b, _ := src.NewBuffer(winSize)
		draw.Draw(b.RGBA(), b.Bounds(), frame.A.Back, image.ZP, draw.Src)
		fr := frame.New(b.RGBA(), b.Bounds(), nil)
		fr.Refresh()
		wind.Send(paint.Event{})
		ck := func() {
			if dirty || fr.Dirty() {
				wind.Send(paint.Event{})
			}
			dirty = false
		}

		flush := func() {
			wind.Upload(fr.Bounds().Min, b, fr.Bounds())
			wind.Publish()
		}
		for {
			switch e := wind.NextEvent().(type) {
			case mouse.Event:
				if e.Button == 1 && e.Direction == 1 {
					p0 := fr.IndexOf(pt(e))
					fr.Select(p0, p0)
					flush()
					frameSweep(fr, wind, flush)
					wind.Send(paint.Event{})
				}
			case key.Event:
				if e.Direction == 2 {
					continue
				}
				if e.Rune == '\r' {
					e.Rune = '\n'
				}
				p0, p1 := fr.Dot()
				if e.Rune > 0x79 || e.Rune < 0 {
					continue
				}
				if e.Rune == '\x08' {
					if p0 == p1 && p0 > 0 {
						p0--
					}
					fr.Delete(p0, p1)
				} else {
					if p0 != p1 {
						fr.Delete(p0, p1)
					}
					fr.Insert([]byte{byte(e.Rune)}, p0)
				}
				dirty = true
				ck()
			case size.Event:
				wind.Upload(image.ZP, b, b.Bounds())
				fr.Refresh()
				ck()
			case paint.Event:
				wind.Upload(fr.Bounds().Min, b, fr.Bounds())
				fr.Flush()
				wind.Publish()
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			}
		}
	})
}

func pt(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

func frameSweep(f *frame.Frame, ep screen.Window, flush func()) {
	p0, p1 := f.Dot()
	for {
		switch e := ep.NextEvent().(type) {
		case mouse.Event:
			if e.Direction != 0 {
				ep.SendFirst(e)
				return
			}
			if p1 = f.IndexOf(pt(e)); p0 > p1 {
				f.Select(p1, p0)
			} else {
				f.Select(p0, p1)
			}
			flush()
		case interface{}:
			ep.SendFirst(e)
			return
		}
	}
}
