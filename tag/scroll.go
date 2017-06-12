package tag

import (
	"golang.org/x/mobile/event/mouse"
	"time"
)

var (
	scrolldy int
	ticking  bool
)

type ScrollEvent struct {
	dy        int
	wind      *Invertable
	flushwith func(e interface{})
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
				act.SendFirst(ScrollEvent{dy: scrolldy, wind: act, flushwith: act.SendFirst})
				scrolldy = 0
			})
		} else {
			scrolldy += dy
		}
	}
}
