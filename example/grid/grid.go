package main

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"sync"

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
	pad     = image.Pt(15, 5)
	fsize   = 11
)

type Plane interface {
	Loc() image.Rectangle
	Move(image.Point)
	Resize(image.Point)
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

type Col struct {
	sp   image.Point
	size image.Point
	wind screen.Window
	Tag  *tag.Tag
	tdy  int
	List []Plane
}

var cols = frame.Acme

func NewCol(src screen.Screen, wind screen.Window, ft frame.Font, sp, size image.Point, files ...string) *Col {
	N := len(files)
	tdy := ft.Dy() * 2
	T := tag.NewTag(src, wind, ft, image.Pt(sp.X, sp.Y), image.Pt(size.X, tdy), pad, cols)
	//T.Open("tag")
	T.Wtag.InsertString("New Delcol Sort", 0)
	T.Wtag.Scroll = nil
	size.Y -= tdy
	sp.Y += tdy
	dy := image.Pt(size.X, size.Y/N)
	col := &Col{sp: sp, size: size, wind: wind, Tag: T, tdy: tdy, List: make([]Plane, len(files))}
	for i, v := range files {
		fmt.Printf("tag %d sp=%s size=%s\n", i, sp, dy)
		t := tag.NewTag(src, wind, ft, sp, dy, pad, cols)
		t.Open(v)
		col.List[i] = t
		sp.Y += dy.Y
	}
	col.List = append([]Plane{T}, col.List...)
	return col
}

func (co *Col) Upload(wind screen.Window) {
	for _, t := range co.List {
			t.(*tag.Tag).Upload(wind)
	}
}

func (co *Col) Resize(size image.Point) {
	co.size = size
	co.Tag.Resize(image.Pt(size.X, co.tdy))
	size.Y = co.size.Y-co.tdy
	sp := co.sp
	sp.Y += co.tdy
	N := len(co.List) - 1
	dy := image.Pt(size.X, size.Y/N)
	for _, t := range co.List[1:] {
		t.Move(sp)
		t.Resize(dy)
		sp.Y += dy.Y
	}
}

func (co *Col) Loc() image.Rectangle {
	return image.Rectangle{co.sp, co.sp.Add(co.size)}
}
func (co *Col) Move(sp image.Point) {
	
	co.sp = sp
	co.Tag.Move(sp)
	sp = sp.Add(image.Pt(0, co.tdy))
	for _, t := range co.List[1:] {
		t.Move(sp)
		sp.Y += t.Loc().Dy()
	}
}

func (co *Col) Handle(act *tag.Invertable, e interface{}) {
	for i := range co.List {
		t := co.List[i].(*tag.Tag)
		t.Handle(t.W, e)
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

		list := []string{"/dev/stdin"}
		if len(os.Args) > 1 {
			list = os.Args[1:]
		}
		if len(list) < 2 {
			list = append(list, "guide")
		}
		sp := image.ZP
		dx := winSize.X / 2
		x := dx
		co0 := NewCol(src, wind, ft, sp, image.Pt(sp.X+x, winSize.Y), list[len(list)/2:]...)
		sp.X += dx
		co1 := NewCol(src, wind, ft, sp, image.Pt(sp.X+x, winSize.Y), list[:len(list)/2]...)
		actCol := co0
		actTag := actCol.List[1]
		act := actTag.(*tag.Tag).W

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
/*
				pt := image.Pt(int(e.X), int(e.Y))
				if pt.In(act.Loc()) {
					actTag.(*tag.Tag).Handle(act, e)
					break
				}
*/
				actCol = active(e, actCol, co0, co1).(*Col)
				actTag = active(e, actTag, actCol.List...).(*tag.Tag)
				act = active(e, act, actTag.(*tag.Tag).W, actTag.(*tag.Tag).Wtag).(*tag.Invertable)
				actTag.(*tag.Tag).Handle(act, e)
			case string, *tag.Command, tag.ScrollEvent, key.Event:
				actTag.(*tag.Tag).Handle(act, e)
			case size.Event:
				winSize = image.Pt(e.WidthPx, e.HeightPx)
				sp := image.ZP
				dx := winSize.X / 2
				x := dx
				co0.Move(sp)
				co0.Resize(image.Pt(sp.X+x, winSize.Y))
				sp.X += dx
				co1.Move(sp)
				co1.Resize(image.Pt(sp.X+x, winSize.Y))
				act.SendFirst(paint.Event{})
			case paint.Event:
				var wg sync.WaitGroup
				wg.Add(2)
				go func(){ co0.Upload(wind); wg.Done()}()
				go func(){co1.Upload(wind); wg.Done()}()
				wg.Wait()
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
