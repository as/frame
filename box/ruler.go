package box

import "bytes"

func newRuler(measureByte func(b byte) int, cap int) *ruler {
	r := &ruler{
		measureByte: measureByte,
		Buffer:      bytes.NewBuffer(make([]byte, cap)),
		widthc:      make(chan bool),
		closec:      make(chan bool),
		rstc:        make(chan bool),
		widthretc:   make(chan int),
		w:           0,
		inc:         make(chan byte),
	}
	go r.measureProc()
	return r
}

type ruler struct {
	measureByte func(b byte) int
	*bytes.Buffer
	widthc    chan bool
	closec    chan bool
	rstc      chan bool
	widthretc chan int
	inc       chan byte
	w         int
}

func (r *ruler) Close() {
	r.closec <- true
}

func (r *ruler) Reset() {
	r.rstc <- true
}

func (r *ruler) Width() (widthpx int) {
	r.widthc <- true
	return <-r.widthretc
}

func (r *ruler) WriteByte(c byte) (err error) {
	r.inc <- c
	return r.Buffer.WriteByte(c)
}

func (r *ruler) drain() {
	for len(r.inc) != 0 {
		<-r.inc
	}
}

func (r *ruler) measureProc() {
	for {
		select {
		case <-r.closec:
			return
		case <-r.rstc:
			r.w = 0
			r.drain()
			r.Buffer.Reset()
		case b := <-r.inc:
			r.w += r.measureByte(b)
		case <-r.widthc:
			for len(r.inc) != 0 {
				b := <-r.inc
				r.w += r.measureByte(b)
			}
			r.widthretc <- r.w
		}
	}
}
