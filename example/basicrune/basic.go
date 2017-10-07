package main

import (
	"github.com/as/frame"
	"github.com/as/frame/font"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"image"
	"image/draw"
	"sync"
)

var wg sync.WaitGroup
var winSize = image.Pt(1024, 768)

func pt(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

func main() {
	driver.Main(func(src screen.Screen) {
		var dirty = true
		wind, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y, "basic"})
		b, _ := src.NewBuffer(winSize)
		draw.Draw(b.RGBA(), b.Bounds(), frame.A.Back, image.ZP, draw.Src)
		fr := frame.New(image.Rectangle{image.ZP, winSize}, font.NewGoMono(14), b.RGBA(), frame.A, true)
		fr.Refresh()
		wind.Send(paint.Event{})
		ck := func() {
			if dirty || fr.Dirty() {
				wind.Send(paint.Event{})
			}
			dirty = false
		}
		utf := `Π Ρ ΢ Σ Τ Υ Φ Χ Ψ Ω Ϊ Ϋ ά έ ή ί ΰ α β γ δ ε ζ η θ ι κ λ μ ν ξ οπ ρ ς σ τ υ φ χ ψ ω ϊ ϋ ό ύ ώ Ϗ ϐ ϑ ϒ ϓ ϔ ϕ ϖ ϗ Ϙ ϙ Ϛ ϛ Ϝ ϝ Ϟ ϟϠ ϡ Ϣ ϣ Ϥ ϥ Ϧ ϧ Ϩ ϩ Ϫ ϫ Ϭ ϭ Ϯ ϯ ϰ ϱ ϲ ϳ ϴ ϵ ϶ Ϸ ϸ Ϲ Ϻ ϻ ϼ Ͻ Ͼ Ͽ`
		fr.Insert([]byte("utf8 test"), fr.Len())
		fr.Insert([]byte(utf), fr.Len())
		fr.Insert([]byte("end"), fr.Len())
		for {
			switch e := wind.NextEvent().(type) {
			case mouse.Event:
				if e.Button == 1 && e.Direction == 1 {
					p0 := fr.IndexOf(pt(e))
					fr.Select(p0, p0)
					flush := func() {
						for _, r := range fr.Cache() {
							wind.Upload(r.Min, b, r)
						}
						fr.Flush()
						wind.Publish()
					}
					flush()
					fr.Sweep(wind, flush)
					wind.Send(paint.Event{})
				}
			case key.Event:
				if e.Direction == 2 {
					continue
				}
				if e.Rune == '\r' {
					e.Rune = '\n'
				}
				if e.Rune > 0x79 || e.Rune < 0 {
					continue
				}
				p0, p1 := fr.Dot()
				if e.Rune == '\x08' {
					if p0 == p1 && p0 > 0 {
						p0--
					}
					fr.Delete(p0, p1)
				} else {
					fr.Insert([]byte{byte(e.Rune)}, p0)
					p0++
				}
				fr.Select(p0, p0)
				dirty = true
				ck()
			case size.Event:
				wind.Upload(image.ZP, b, b.Bounds())
				fr.Refresh()
				ck()
			case paint.Event:
				for _, r := range fr.Cache() {
					wind.Upload(r.Min, b, r)
				}
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
