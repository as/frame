// You might want to look here if your example program has been running slowly
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


var (
	workc chan image.Rectangle
	gwind screen.Window
	gbuf  screen.Buffer
	wg    sync.WaitGroup

	winSize = image.Pt(1024, 768)
)

func main() {
	driver.Main(func(src screen.Screen) {
		// This is exp/shiny specific
		wind, _ := src.NewWindow(&screen.NewWindowOptions{Width: winSize.X, Height: winSize.Y, Title: "fast"})
		b, _ := src.NewBuffer(winSize)

		// Assign
		gbuf, gwind = b, wind

		// Draw the background color once on the entire window
		draw.Draw(b.RGBA(), b.Bounds(), frame.A.Back, image.ZP, draw.Src)

		// Create a new frame using a 14-pt GoMono font
		fr := frame.New(image.Rectangle{image.ZP, winSize}, font.NewGoMono(14), b.RGBA(), frame.A)
		fr.Refresh()

		// Prevent redundant redraws
		var dirty = true
		ck := func() {
			if dirty || fr.Dirty() {
				wind.Send(paint.Event{})
			}
			dirty = false
		}

		// Insert text into the frame

		fr.Insert([]byte("This is the fast example. It's harder to understand but it uses 4 worker\n"), fr.Len())
		fr.Insert([]byte("goroutines to update the shiny window. You can change the flags to frame.New\n"), fr.Len())
		fr.Insert([]byte("to re-create the utf8 and elastic tabstop examples\n"), fr.Len())

		fr.Insert([]byte(`
	fr := frame.New(image.Rectangle{image.ZP, winSize}, font.NewGoMono(14), b.RGBA(), frame.A) // <- you are here
		
// As you can see there are two new package scoped variables, gwind and gbuf
// these are used by the four workers created in the init function to
// concurrently upload pieces of the screen. The workc is a channel used to
// transmit rectangles, and wg tells the shiny event loop when the upload
// is completed and the window can be published
var(
	workc chan image.Rectangle
	gwind screen.Window
	gbuf screen.Buffer
	wg sync.WaitGroup
)

// The uploader goroutines are born here. The only thing they
// do is recv rectangles and call window.Upload() to draw
// that rectangular region to the window. A naive but simpler
// and more readable example would just draw the entire window
// to the screen on every change, but this gets unnecessarily
// expensive for high resolution windows
func init(){
	workc=make(chan image.Rectangle,4)
	for i := 0; i < 4; i++{
		go func(){
			for{
				select{
				case r := <- workc:
					gwind.Upload(r.Min, gbuf, r)
					wg.Done()
				}
			}
		}()
	}
}`), fr.Len())

		wind.Send(paint.Event{})

		// flush makes a copy of the rectangles that haven't been
		// drawn to the screen yet and sends them to the drawing
		// goroutines over a channel. this is an evil but necessary
		// optimization to make shiny draw changes to the window
		// more efficiently
		flush := func() {
			cache := append([]image.Rectangle{}, fr.Cache()...)
			wg.Add(len(cache))
			for _, r := range cache {
				workc <- r
			}
			wind.Publish()
			wg.Wait()
			fr.Flush()
		}

		for {
			switch e := wind.NextEvent().(type) {
			case mouse.Event:
				if e.Button == 1 && e.Direction == 1 {
					p0 := fr.IndexOf(pt(e))
					fr.Select(p0, p0)
					flush()
					fr.Sweep(wind, flush) // sweep calls the flush function internally
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
				flush()
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			}
		}
	})
}

// The uploader goroutines are born here. The only thing they
// do is recv rectangles and call window.Upload() to draw
// that rectangular region to the window. A naive but simpler
// and more readable example would just draw the entire window
// to the screen on every change, but this gets unnecessarily
// expensive for high resolution windows
func init() {
	workc = make(chan image.Rectangle, 4)
	for i := 0; i < 4; i++ {
		go func() {
			for {
				select {
				case r := <-workc:
					gwind.Upload(r.Min, gbuf, r)
					wg.Done()
				}
			}
		}()
	}
}

func pt(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}
