package main

import (
	"sync"
	//	"github.com/as/clip"
	//

	"image"

	"github.com/as/frame"
	"github.com/as/frame/font"
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
	var focused = false
	driver.Main(func(src screen.Screen) {
		var dirty = true
		wind, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y, "basic"})
		b, _ := src.NewBuffer(winSize)
		fr := frame.New(image.Rectangle{image.ZP, winSize}, font.NewGoMono(12), b.RGBA(), frame.Acme)
		ck := func() {
			if dirty || fr.Dirty() {
				wind.Send(paint.Event{})
			}
			dirty = false
		}
		for {
			switch e := wind.NextEvent().(type) {
			case mouse.Event:
				if e.Button != 1 && e.Direction != 1 {
					continue
				}
				p0 := fr.IndexOf(image.Pt(int(e.X), int(e.Y)))
				println(p0)
				fr.Select(p0, p0)
				ck()
			case key.Event:
				if e.Direction == 2 {
					continue
				}
				if e.Rune == '\r' {
					e.Rune = '\n'
				}
				p0, _ := fr.Dot()
				println(p0)
				fr.Insert([]byte{byte(e.Rune)}, p0)
				fr.Select(p0+1, p0+1)
				dirty = true
				ck()
			case size.Event:
				fr.Refresh()
				ck()
			case paint.Event:
				for _, r := range fr.Cache() {
					wind.Upload(image.ZP, b, r)
				}
				fr.Flush()
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
