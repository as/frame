package block

import "testing"

func newBoxesFixed() *Run {
	return NewRun()
}
func TestReadAt(t *testing.T) {
	b := newBoxesFixed()
	ck := func(bn int, want string) {
		have := "<nil>"
		if b.Block[bn].Ptr != nil {
			have = string(b.Block[bn].Ptr[:b.Nchars])
		}
		if have != want {
			t.Logf("box #%d: want %q have %q\n", bn, want, have)
			t.FailNow()
		}
	}
	_ = ck
	_ = b
	if false {
		b.Dump()
		b.WriteAt([]byte("mink"), 0)
		ck(0, "mink")
		b.WriteAt([]byte("s"), 0)
		b.WriteAt([]byte("u"), 1)
		b.WriteAt([]byte("r"), 2)
		b.WriteAt([]byte("p"), 2)
		b.WriteAt([]byte("e"), 3)
		b.Dump()
		ck(1, "supe")
		b.WriteAt([]byte(" mink!!!!!!"), 5)
		want := "super mink"
		have := make([]byte, len(want))
		b.ReadAt(have, 2)
		if string(have) != want {
			t.Logf("readat: want %q have %q\n", want, have)
			t.Fail()
		}
	}
}

func TestWriteAt(t *testing.T) {
	b := newBoxesFixed()
	ck := func(want string, at int) {
		have := make([]byte, len(want))
		b.Dump()
		b.ReadAt(have, int64(at))
		b.Dump()
		if have, want := string(have), string(want); have != want {
			t.Logf("writeat %d: want %q have %q\n", at, want, have)
			t.FailNow()
		}
	}
	want := "0123456789......"
	b.WriteAt([]byte(want), 0)
	b.Dump()
	ck(want, 0)

	at := len("0123456789")
	b.WriteAt([]byte("abcdef"), int64(at))
	ck("0123456789abcdef", 0)

	at = len("0123456789")
	b.Delete(at, at+len("abcdef"))

	//b.WriteAt([]byte("or so you think..."), int64(at))
}

func xTestCommon(t *testing.T) {
	b := newBoxesFixed()
	ck := func(want string, at int) {
		have := make([]byte, len(want))
		b.ReadAt(have, int64(at))
		if have, want := string(have), string(want); have != want {
			t.Logf("writeat %d: want %q have %q\n", at, want, have)
			t.FailNow()
		}
	}
	s := "The quick brown fox jumped over the lazy dog."
	for i, v := range []byte(s) {
		b.WriteAt([]byte{v}, int64(i))
	}
	ck(s, 0)
}

func xTestDelete(t *testing.T) {
	b := newBoxesFixed()
	ck := func(want string, at int) {
		have := make([]byte, len(want))
		b.ReadAt(have, int64(at))
		if have, want := string(have), string(want); have != want {
			t.Logf("writeat %d: want %q have %q\n", at, want, have)
			t.FailNow()
		}
	}
	for i, v := range []byte("The quick brown fox jumped over the lazy dog.") {
		b.WriteAt([]byte{v}, int64(i))
	}
	b.Delete(len("The "), len("The ")+len("quick brown fox jumped over the"))
	//	b.Dump()
	ck = ck
}

func xTestWriteAtZero(t *testing.T) {
	b := newBoxesFixed()
	ck := func(want string, at int) {
		have := make([]byte, len(want))
		b.ReadAt(have, int64(at))
		if have, want := string(have), string(want); have != want {
			t.Logf("writeat %d: \n\twant %q \n\thave %q\n", at, want, have)
			t.FailNow()
		}
	}
	for _, v := range []byte("The quick brown fox jumped over the lazy dog.") {
		b.WriteAt([]byte{v}, 0)
	}
	ck(".god yzal eht revo depmuj xof nworb kciuq ehT", 0)
}

func xTestWriteAtPlusOne(t *testing.T) {
	b := newBoxesFixed()
	ck := func(want string, at int) {
		have := make([]byte, len(want))
		b.ReadAt(have, int64(at))
		if have, want := string(have), string(want); have != want {
			t.Logf("writeat %d: \n\twant %q \n\thave %q\n", at, want, have)
			t.FailNow()
		}
	}
	for i, v := range []byte("The quick brown fox jumped over the lazy dog.") {
		b.WriteAt([]byte{v}, int64(i+1))
	}
	//b.Dump()
	ck = ck
}
