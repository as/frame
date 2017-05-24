package win

/*
	Put
	optimization for SetSelect broken
	hold mouse and move up - broken (cache upload)
	TODO - print statements in set select, why is selection not working until mouse let go?
	Delete	SetSe
 */

import (
	"fmt"
	"image"
	"image/draw"
	"sync"
	//"time"

	"github.com/as/frame"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/mouse"
	// "golang.org/x/mobile/event/paint"	
)

const (
	HiWater  = 1024 * 1024
	LoWater  = 2 * 1024
	MinWater = 1024
	MsgSize  = 64 * 1024
)

type Win struct {
	*frame.Frame
	Sp      image.Point // window offset
	size    image.Point // window size
	pad     image.Point // window text offset
	b       screen.Buffer
	scr     screen.Screen
	events  screen.Window
	Org     int64
	Qh      int64
	Q0, Q1  int64
	Nr      int64
	R       []byte
	Maxr    int64
	Mc      Mc
	Selectq int64
	Scrollr image.Rectangle
	Sweep bool
}

func (w *Win) Buffer() screen.Buffer {
	return w.b
}
func (w *Win) Size() image.Point {
	return w.size
}

func New(scr screen.Screen, ft frame.Font, events screen.Window,
	sp, size, pad image.Point, cols frame.Color) *Win {
	b, err := scr.NewBuffer(size)
	if err != nil {
		panic(err)
	}
	r := image.Rectangle{pad, size}.Inset(1)
	w := &Win{
		Frame:  frame.New(r, ft, b.RGBA(), cols),
		b:      b,
		pad:    pad,
		scr:    scr,
		Sp:     sp,
		size:   size,
		events: events,
	}
	w.Frame.Scroll = w.FrameScroll
	w.Blank()
	return w
}

func (w *Win) NextEvent() (e interface{}) {
	switch e := w.events.NextEvent().(type) {
	case mouse.Event:
		e.X -= float32(w.Sp.X)
		e.Y -= float32(w.Sp.Y)
		return e
	case interface{}:
		return e
	}
	return nil
}
func (w *Win) Send(e interface{}) {
	w.events.Send(e)
}
func (w *Win) SendFirst(e interface{}) {
	w.events.SendFirst(e)
}
func (w *Win) Blank() {
	buf := w.b
	if buf != nil {
		r := buf.RGBA().Bounds()
		draw.Draw(buf.RGBA(), r, w.Color.Back, image.ZP, draw.Src)
		if w.Sp.Y > 0 {
			r.Min.Y--
			//r.Max.Y--
		}
		//		drawBorder(buf.RGBA(), r, w.Color.Hi.Back, image.ZP, 1)
	}
}

type Mc struct {
	Buttons int
	Msec    int
	XY      image.Point
}

func (w *Win) Readmouse(e mouse.Event) {

}

/*
flush
FrameScroll
	//fmt.Printf("fs3 q0 = %d\n", q0)

	w.SetOrigin(q0, true)
	//w.Redraw()
	//w.Drawsel(w.Frame.PtOfChar(w.P0), w.P0, w.P1, true)
	r := w.Bounds()
	sp := image.Pt(1,1).Add(w.Sp)
	N  := 4
	Ny := r.Dy()/4
	for i := 0; i < N; i++{
		r0 := image.Rect(r.Min.X, r.Min.Y+(Ny*i), r.Max.X, Ny*(i+1))
		go func(i int, r0 image.Rectangle){
			w.events.Upload(sp.Add(image.Pt(0, Ny*3)), w.b, r0); 
			wg.Done()
		}(i, r0)
	}
	wg.Wait()
*/

func (w *Win) Dot() (q0, q1 int64){
	q0 = clamp(w.Q0, 0, w.Nr)
	q1 = clamp(w.Q1, 0, w.Nr)
	return
}

func (w *Win) FrameScroll(dl int) {
	if dl == 0 {
		//time.Sleep(15*time.Millisecond)
		return
	}
	q0 := int64(0)
	if dl < 0 {
		q0 = w.BackNL(w.Org, -dl)
		if w.Sweep{
			if w.Selectq > w.Org+w.P0 {
				//fmt.Printf("a\n")
				 x := w.Selectq
				//x := w.Org+w.P1
				w.SetSelect(w.Org+w.P0, x)
			} else {
				//fmt.Printf("b %d:%d\n", w.Selectq, w.Org+w.P1)	// was w.P0 for both
				// x := w.Selectq
				x := w.Org+w.P0
				w.SetSelect(x, w.Org+w.P1)
			}
		}
	} else {
		if w.Org+w.Nchars == w.Nr {
			return
		}
		r := w.Frame.Bounds()
		q0 = w.Org + w.IndexOf(image.Pt(r.Min.X, r.Min.Y+dl*w.Font.Dy()))
		if w.Sweep{
			if w.Selectq >= w.Org+w.P1 {
				w.SetSelect(w.Org+w.P1, w.Selectq)
			} else {
				w.SetSelect(w.Selectq, w.Org+w.P1)
			}
		}
	}
	//fmt.Printf("fs3 q0 = %d\n", q0)
	if w.Sweep {
		w.flush()
	}
	w.SetOrigin(q0, true)
	
}

// Put	sp.Add	Upload
func (w *Win) SetSelect(q0, q1 int64) {
	//fmt.Printf("SetSelect: [%d:%d]\n", q0, q1)


	w.Q0, w.Q1 = q0, q1
	p0 := clamp(q0-w.Org, 0, w.Nchars)
	p1 := clamp(q1-w.Org, 0, w.Nchars)
	if p0 == w.P0 && p1 == w.P1 {
		return
	}
	if w.P1 <= p0 || p1 <= w.P0 || p0 == p1 || w.P1 == w.P0 {
		w.Drawsel(w.PtOfChar(w.P0), w.P0, w.P1, false)
		w.Drawsel(w.PtOfChar(p0), p0, p1, true)
	} else {
		step := func(i, j int64) {
			//
			// fmt.Printf("step %d,%d %b\n",i,j,i<j)
			if i < j {
				w.Drawsel(w.PtOfChar(i), i, j, true)
			} else if i > j {
				w.Drawsel(w.PtOfChar(j), j, i, false)
			}
		}
		step(p0, w.P0) // trim or extend origin
		step(w.P1, p1) // trim or extend insertion
	}
	w.P0 = p0
	w.P1 = p1
}

func (w *Win) BackNL(p int64, n int) int64 {
	if n == 0 && p > 0 && w.R[p-1] != '\n' {
		n = 1
	}
	for i := n; i > 0 && p > 0; {
		i--
		p--
		if p == 0 {
			break
		}
		for j := 128; j-1 > 0 && p > 0; p-- {
			j--
			if w.R[p-1] == '\n' {
				break
			}
		}
	}
	return p
}

/*	SetS
func (w *Win) Drawscroll(){
	r := w.Scrollr
	b := w.Scrolltmp
	r1 := r
	r1.Min.X = 0
	r1.Max.X = r.Dx()
	r2 := scrpos(r1, w.Org, w.Org+w.Nchars, w.Nr)
}
*/

func (w *Win) SetOrigin(org int64, exact bool) {
	//fmt.Printf("SetOrigin: %d %v\n", org, exact)
	org = clamp(org, 0, w.Nr)
	if org > 0 && !exact {
		for i := 0; i < 256 && org < w.Nr; i++ {
			if w.R[org] == '\n' {
				org++
				break
			}
			org++
		}
	}
	//fmt.Printf("SetOrigin: found %d %v\n", org, exact) Put
	a := org - w.Org // distance to new origin
	fix := false
	if a >= 0 && a < w.Nchars {
		// a bytes to the right; intersects the frame
		// fmt.Printf("Delete(%d,%d)\n", 0,a)
		w.Frame.Delete(0, a)
		fix = true
	} else if a < 0 && -a < w.Nchars {
		// -a bytes to the left; intersects the frame		
		i := org - a
		j := org
		if i > j {
			i, j = j, i
		}
		i = max(0, i)
		j = min(w.Nr, j)
		// fmt.Printf("-a=%d to left: w.R[%d:%d]\n", -a, i,j)				
		w.Frame.Insert(w.R[i:j], 0)		
	} else {
		w.Frame.Delete(0, w.Nchars)
	}
	w.Org = org
	w.Fill()
	//w.ScrDraw(w)
	w.SetSelect(w.Q0, w.Q1)
	if fix && w.P1 > w.P0 {
		w.Drawsel(w.PtOfChar(w.P1-1), w.P1-1, w.P1, true);
	}
	//fmt.Printf("p[%d:%d]\n", w.P0, w.P1)
}

func (w *Win) filldebug(){
	// Put
	fmt.Printf("lines/maxlines = %d/%d\n", w.Line(), w.MaxLine())
}

func (w *Win) Fill() {
	//w.filldebug()
	if w.Frame.Full() {
		return
	}
	var rp [MsgSize]byte
	for  !w.Frame.Full() {
		qep := w.Org + w.Nchars
		n := min(w.Nr-qep, 2000)
		if n == 0 {
			break
		}
		//fmt.Printf("w.org=%d w.Nchars=%d\n", w.Org, w.Nchars)
		//fmt.Printf("copy(rp, w.R[%d:%d] (len=%d)\n", qep, qep+n, len(w.R))
		m := copy(rp[:], w.R[qep:qep+n])
		//fmt.Printf("copied %q\n", rp[:m])
		nl := w.MaxLine() - w.Line()
		m = 0
		i := int64(0)
		for i < n {
			if rp[i] == '\n' {
				m++
				if m >= nl {
					i++
					break
				}
			}
			i++
		}
		//fmt.Printf("w.Frame.Insert rp[:%d-%d], %d\n", w.Nchars, i, w.Nchars)
		w.Frame.Insert(rp[:i], w.Nchars)
	}
}

func (w *Win) Delete(q0, q1 int64) {
	n := q1 - q0
	if n == 0 {
		return
	}
	//fmt.Printf("copy(w.R[%d:], w.R[%d:%d])\n", q0, q1, w.Nr-q1)
	copy(w.R[q0:], w.R[q1:][:w.Nr-q1])
	w.Nr -= n
	if q0 < w.Q0 {
		w.Q0 -= min(n, w.Q0-q0)
	}
	if q0 < w.Q1 {
		w.Q1 -= min(n, w.Q1-q0)
	}
	if q1 < w.Qh {
		w.Qh = q0
	} else if q0 < w.Qh {
		w.Org -= n
	}

	if q1 <= w.Org {
		w.Org -= n
	} else if q0 < w.Org+w.Nchars {
		p1 := q1 - w.Org
		p0 := int64(0)
		if p1 > w.Nchars {
			p1 = w.Nchars
		}
		if q0 < w.Org {
			w.Org = q0
		} else {
			p0 = q0 - w.Org
		}
		fmt.Printf("w.Frame.Delete(%d,%d)\n", p0, p1)
		w.Frame.Delete(p0, p1)
		w.Fill()
	}
}

func (w *Win) InsertString(s string, q0 int64) int64 {
	return w.Insert([]byte(s), q0)
}

func (w *Win) Insert(s []byte, q0 int64) int64 {
	// invariant r = p - origin
	//           5 = 5 - 0
	//           4 = 5 - 1
	//fmt.Printf("%p: w.Nr=%d\n", w, w.Nr)
	//fmt.Printf("%p: Insert: len(s)=%d\n", w, len(s))
	n := int64(len(s))
	if n == 0 {
		return q0
	}
	if w.Nr+n > HiWater && q0 >= w.Org && q0 >= w.Qh {
		m := min(HiWater-LoWater, min(w.Org, w.Qh))
		w.Org -= m
		w.Qh -= m
		if w.Q0 > m {
			w.Q0 -= m
		} else {
			w.Q0 = 0
		}
		if w.Q1 > m {
			w.Q1 -= m
		} else {
			w.Q1 = 0
		}
		//fmt.Printf("%p: A w.Nr bef: %d\n", w.Nr)
		w.Nr -= m
		//fmt.Printf("%p: A w.Nr after: %d\n", w.Nr)
		copy(w.R, w.R[m:][:w.Nr])
		q0 -= m
	}
	if w.Nr+n > w.Maxr {
		println("insert if D")
		m := max(min(2*(w.Nr+n), HiWater), w.Nr+n) + MinWater
		if m > HiWater {
			m = max(HiWater+MinWater, w.Nr+n)
		}
		if m > w.Maxr {
			extra := int64(m) - int64(len(w.R))
			w.R = append(w.R, make([]byte, extra)...)
			w.Maxr = m
		}
	}
	copy(w.R[q0+n:], w.R[q0:][:w.Nr-q0])
	copy(w.R[q0:], s[:n])
	//fmt.Printf("%p: B w.Nr bef: %d\n", w, w.Nr)
	w.Nr += n
	//fmt.Printf("%p: B w.Nr after: %d\n", w, w.Nr)
	//fmt.Printf("w.Nr = %d\n", w.Nr)
	if q0 <= w.Q1 {
		w.Q1 += n
	}
	if q0 <= w.Q0 {
		w.Q0 += n
	}
	if q0 < w.Qh {
		w.Qh += n
	}
	if q0 < w.Org {
		w.Org += n
	} else if q0 <= w.Org+w.Nchars {
		//fmt.Printf("w.Frame.Insert: @ %d -> %q\n", n, s)
		n--
		if n < 0 {
			n++
		}
		w.Frame.Insert(s, q0-w.Org)
	}
	//	fmt.Printf("buf: %q\n", w.R)
	return q0
}

func (w *Win) SetFont(ft frame.Font) {
	b := w.Bytes()
	w.Clear(true)
	w.Blank()
	w.Frame = frame.New(w.Frame.Bounds(), ft, w.b.RGBA(), w.Color)
	w.Insert(b, 0)
	w.Redraw()
}

func (w *Win) Resize(size image.Point) {
	w2 := New(w.scr, w.Font, w.events, w.Sp, size, w.pad, w.Color)
	bb := w.Bytes()
	w.b.Release()
	w.b = w2.b
	w.Frame = w2.Frame
	*w = *w2
	w.Sp = w2.Sp
	w.size = w2.size
	w.pad = w2.pad
	w.scr = w2.scr
	w.Blank()
	w.Insert(bb, 0)
}

func (w *Win) upload() {
	w.events.Upload(w.Sp.Add(image.Pt(5,5)), w.b, w.Bounds())
}
func (w *Win) flush(){
	r := w.Bounds()
	sp := image.Pt(1,1).Add(w.Sp).Add(w.pad)
	Ny := r.Dy()/4
	r0 := image.Rect(r.Min.X, r.Min.Y,        r.Max.X, r.Min.Y+Ny)
	r1 := image.Rect(r.Min.X, r.Min.Y+(Ny),   r.Max.X, r.Min.Y+Ny*2)
	r2 := image.Rect(r.Min.X, r.Min.Y+(Ny*2), r.Max.X, r.Min.Y+Ny*3)
	r3 := image.Rect(r.Min.X, r.Min.Y+(Ny*3), r.Max.X, r.Min.Y+Ny*4)
	var wg sync.WaitGroup
	wg.Add(4)
	go func(){w.events.Upload(sp, w.b, r0); wg.Done()}()
	go func(){w.events.Upload(sp.Add(image.Pt(0, Ny)), w.b, r1); wg.Done()}()
	go func(){w.events.Upload(sp.Add(image.Pt(0, Ny*2)), w.b, r2); wg.Done()}()
	go func(){w.events.Upload(sp.Add(image.Pt(0, Ny*3)), w.b, r3); wg.Done()}()	
	w.Flushcache()
	wg.Wait()
}
// Put
func (w *Win) Upload() {
///	w.upload()
//	return
	var wg sync.WaitGroup
	wg.Add(len(w.Cache))
	for _, r := range w.Frame.Cache {
		//fmt.Printf("upload %s\n", r)
		go func(r image.Rectangle) {
			w.events.Upload(w.Sp.Add(r.Min), w.b, r)
			wg.Done()
		}(r)
	}
	//w.events.Upload(w.Sp.Add(image.Pt(0,0)), w.b, image.Rect(0,0,500,500))
	wg.Wait()
	w.Flushcache()
}

func (w *Win) ReadAt(off int64, p []byte) (n int, err error) {
	if off > w.Nr {
		return
	}
	return copy(p, w.R[off:w.Nr]), err

}

func (w *Win) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (w *Win) Bytes() []byte {
	return w.R[:w.Nr]
}

func (w *Win) Rdsel() []byte {
	i := w.Q0
	j := w.Q1
	return w.R[i:j]
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func clamp(v, l, h int64) int64 {
	if v < l {
		return l
	}
	if v > h {
		return h
	}
	return v
}
// Put