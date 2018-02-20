package frame

import (
	"image"
	"strings"
	"testing"

	. "github.com/as/font"
)

func TestInsertOneChar(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert([]byte("1"), 0)
	h.Untick()
	//etch.WriteFile(t, `testdata/TestInsertOneChar.expected.png`, have)
	check(t, have, "TestInsertOneChar", testMode)
}

func TestInsert10Chars(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for i := 0; i < 10; i++ {
		h.Insert([]byte("1"), 0)
	}
	check(t, have, "TestInsert10Chars", testMode)
}

func TestInsert22Chars2Lines(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 2; j++ {
		for i := 0; i < 10; i++ {
			h.Insert([]byte("1"), h.Len())
		}
		h.Insert([]byte("\n"), h.Len())
	}
	check(t, have, "TestInsert22Chars2Lines", testMode)
}

func TestInsert1000(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 1000; j++ {
		h.Insert([]byte{byte(j)}, int64(j))
	}
	check(t, have, "TestInsert1000", testMode)
}

func TestInsertTabSpaceNewline(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 5; j++ {
		h.Insert([]byte("abc\t \n\n\t $\n"), int64(j))
	}
	check(t, have, "TestInsertTabSpaceNewline", testMode)
}

type benchOp struct {
	name string
	data string
	at   int64
}

var dst = image.NewRGBA(image.Rect(0, 0, 1024, 768))

func BenchmarkInsertGoMono(b *testing.B) { b.Helper(); bInsert(b, withFace(NewGoMono(fsize))) }
func BenchmarkInsertGoMonoCache(b *testing.B) {
	if b, ok := interface{}(b).(help); ok {
		b.Helper()
	}
	bInsert(b, withFace(NewCache(NewGoMono(fsize))))
}
func BenchmarkInsertGoMonoReplaceCache(b *testing.B) {
	if b, ok := interface{}(b).(help); ok {
		b.Helper()
	}
	bInsert(b, withFace(
		NewCache(
			Replacer(
				NewGoMono(fsize), NewHex(fsize), nil,
			),
		),
	),
	)
}
func BenchmarkInsertGoMonoCliche(b *testing.B) {
	if b, ok := interface{}(b).(help); ok {
		b.Helper()
	}
	bInsert(b, withFace(NewCliche(NewGoMono(fsize))))
}
func BenchmarkInsertGoMonoRune(b *testing.B) {
	if b, ok := interface{}(b).(help); ok {
		b.Helper()
	}
	bInsert(b, withFace(NewRune(NewGoMono(fsize))))
}
func BenchmarkInsertGoMonoRuneCache(b *testing.B) {
	if b, ok := interface{}(b).(help); ok {
		b.Helper()
	}
	bInsert(b, withFace(NewCache(NewRune(NewGoMono(fsize)))))
}

func withFace(ft Face) *Frame {
	return New(dst, dst.Bounds(), &Config{Face: ft, Color: A})
}

func bInsert(b *testing.B, f *Frame) {
	//	b.Skip("not ready")
	if b, ok := interface{}(b).(help); ok {
		b.Helper()
	}
	for _, v := range []benchOp{
		{"1", "a", 0},
		{"10", strings.Repeat("a\n", 10), 0},
		{"100", strings.Repeat("a\n", 100), 0},
		{"1000", strings.Repeat("a\n", 1000), 0},
		{"10000", strings.Repeat("a\n", 10000), 0},
		{"100000", strings.Repeat("a\n", 100000), 0},
	} {
		b.Run(v.name, func(b *testing.B) {
			bb := []byte(v.data)
			b.SetBytes(int64(len(bb)))
			at := v.at
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f.Insert(bb, at)
			}
		})
	}
}

// help is an interface that allows this code to use Go1.9's t.Helper() method
// without breaking out of data continuous integration components (CircleCI) which
// run older Go versions not supporting t.Helper().
//
//
type help interface {
	Helper()
}
