package frame

import (
	"bufio"
	"fmt"
	"github.com/as/etch"
	"github.com/as/io/spaz"
	"image"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func TestFuzz(t *testing.T) {
	t.Skip("warning: fuzz test skipped")
	// The inverse of Insert is Delete. We can use this assumption
	// to create a graphical fuzz test.
	var (
		err error
		n   int
		B   [327 * 777]byte
	)
	buf := B[:]
	N := 128 // number of rounds
	sr := spaz.NewReader(bufio.NewReader(reader{}))
	fr, fr2, a, b := abtestbg(image.Rect(0, 0, 327, 771))
	for i := 0; i < N; i++ {
		// Sync: Start by inserting something into to both frames
		n, err = sr.Read(buf)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		in := buf[:n]
		p0 := clamp(rand.Int63(), 0, fr.Len())
		fr.Insert(in, p0)
		fr2.Insert(in, p0)

		// Insert something random into the first frame
		// and see if Deleting that gets us coherent with the
		// second frame

		prior := buf[:n] // In case it fails, we need a record of the random data
		n, err = sr.Read(buf)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		in = buf[:n]
		p0 = clamp(rand.Int63(), 0, fr.Len())
		ni := int64(fr.Insert(in, p0))
		for i := 0; i < 1000; i++ {
			fr.Select(clamp(rand.Int63(), 0, fr.Len()), clamp(rand.Int63(), 0, fr.Len()))
		}
		nd := fr.Delete(p0, p0+ni)
		delta, ok := etch.Delta(a, b)
		if !ok {
			name := fmt.Sprintf("TestFuzz.Len%dIns%dCharsAt%d", len(prior), len(in), p0)
			result := etch.Report(a, b, delta)
			info := fmt.Sprintf("bufprior = %q\n", prior)
			info += fmt.Sprintf("insert = %q\n", in)
			info += fmt.Sprintf("p0 = %d\n", p0)
			info += fmt.Sprintf("nins = %d\n", ni)
			info += fmt.Sprintf("delete = %d:%d\n", p0, p0+ni)
			info += fmt.Sprintf("ndel = %d\n", nd)
			ioutil.WriteFile(name+".info", []byte(info), 0666)
			t.Logf("see %s.png and %s.buf\n", name, name)
			etch.WriteFile(t, name+".png", result)
			t.FailNow()
		}
	}
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

func init() {
	rand.Seed(time.Now().Unix())
}

type reader struct {
}

func (reader) Read(p []byte) (n int, err error) {
	return rand.Read(p)
}
