package box

import (
	"testing"

	"github.com/as/font"
)

// helpful is an interface that allows this code to use Go1.9's t.Helper() method
// without breaking out of data continuous integration components (CircleCI) which
// run older Go versions not supporting t.Helper().
//
//
type help interface {
	Helper()
}

// genrun generates a run with a pre-fixed font and dimensions
func genrun(s ...string) (r Run) {
	return Run{
		delta:  32,
		minDx:  5,
		maxDx:  5000,
		Face:   font.NewGoMono(11),
		Nalloc: len(s),
		Nbox:   len(s),
		Box:    genbox(5, 5000, 10, s...),
	}
}

func genbox(min, max, fontdx int, s ...string) (bx []Box) {
	if min == 0 {
		min = 5
	}
	if max == 0 {
		max = 5000
	}
	if fontdx == 0 {
		fontdx = 10
	}
	for _, s := range s {
		b := Box{
			Nrune: len(s),
			Ptr:   append(make([]byte, 0, MaxBytes), []byte(s)...),
		}
		if s == "\t" {
			b.Width = min
			b.Minwidth = min
		} else if s == "\n" {
			b.Nrune = -1
			b.Width = 5000
		} else {
			b.Width = fontdx * b.Nrune
		}
		bx = append(bx, b)
	}
	return bx
}

func runCk(t *testing.T, have Run) {
	if t, ok := interface{}(t).(help); ok {
		t.Helper()
	}
	boxCk(t, have.Box)
}

func boxCk(t *testing.T, have []Box) {

	if t, ok := interface{}(t).(help); ok {
		t.Helper()
	}
	for bn, h := range have {
		if h.Nrune < -1 {
			t.Logf("box %d: should never have Nrune < -1", bn)
			t.Fail()
		}
		if h.Nrune == -1 && len(h.Ptr) > 1 {
			t.Logf("box %d: Nrune < -1 but len(Ptr) > 1", bn)
			t.Fail()
		}
		if h.Nrune > MaxBytes {
			t.Logf("box %d: h.Nrune [%d] > MaxBytes [%d]", bn, h.Nrune, MaxBytes)
			t.Fail()
		}
	}
}

func runCompare(t *testing.T, strict bool, have, want Run) {

	if t, ok := interface{}(t).(help); ok {
		t.Helper()
	}
	runCk(t, have)
	runCk(t, want)
	h, w := have.Box, want.Box
	if !strict {
		h = h[:have.Nbox]
		w = w[:want.Nbox]
	}
	boxCompare(t, h, w)
}

func boxCompare(t *testing.T, have, want []Box) {
	if t, ok := interface{}(t).(help); ok {
		t.Helper()
	}

	failed := false
	fail := func() {
		t.Fail()
		failed = true
	}
	defer func() {
		if failed {
			dumpBoxes(have)
			dumpBoxes(want)
		}
	}()
	if len(want) != len(have) {
		t.Logf("box counts differ: have %d, want %d\n", len(have), len(want))
		fail()
	}
	boxCk(t, have)
	boxCk(t, want)
	for bn := 0; bn < len(have); bn++ {
		h := have[bn]
		w := want[bn]
		if h.Nrune != w.Nrune {
			t.Logf("box reported sizes differ: have %d, want %d\n", h.Nrune, w.Nrune)
			fail()
		}
		sh, wh := string(h.Bytes()), string(w.Bytes())
		if sh != wh {
			t.Logf("box contents differ: have: \n\t%q, want:\n\t%q\n", h.Nrune, w.Nrune)
			fail()
		}
	}

}
func TestCombine(t *testing.T) {
	r0 := genrun("hello")
	r1 := genrun("world")
	rWant := genrun("hello", "world")
	r0.Combine(&r1, 1)
	runCompare(t, false, r0, rWant)
}
