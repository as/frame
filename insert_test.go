package frame

import (
	. "github.com/as/font"
	"image"
	"strings"
	"testing"
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

func BenchmarkInsertDelete(b *testing.B) {
	dst := image.NewRGBA(image.Rect(0, 0, 1024, 768))
	for _, v := range []benchOp{
		{"1", "a", 0},
		{"10", strings.Repeat("a", 10), 0},
		{"100", strings.Repeat("a", 100), 0},
		{"1000", strings.Repeat("a", 1000), 0},
		{"10000", strings.Repeat("a", 10000), 0},
		{"100000", strings.Repeat("a", 100000), 0},
	} {
		b.Run(v.name, func(b *testing.B) {
			h := New(dst, dst.Bounds(), tconf())
			bb := []byte(v.data)
			b.SetBytes(int64(len(bb)))
			at := v.at
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Insert(bb, at)
				h.Delete(0, at)
			}
		})
	}
}

func BenchmarkInsertDeleteStaticFace(b *testing.B) {
	dst := image.NewRGBA(image.Rect(0, 0, 1024, 768))
	for _, v := range []benchOp{
		{"1", "a", 0},
		{"10", strings.Repeat("a", 10), 0},
		{"100", strings.Repeat("a", 100), 0},
		{"1000", strings.Repeat("a", 1000), 0},
		{"10000", strings.Repeat("a", 10000), 0},
		{"100000", strings.Repeat("a", 100000), 0},
	} {
		b.Run(v.name, func(b *testing.B) {
			h := New(dst, dst.Bounds(), &Config{
				Font:  NewCache(NewGoMono(fsize)),
				Color: A,
			})
			bb := []byte(v.data)
			b.SetBytes(int64(len(bb)))
			at := v.at
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Insert(bb, at)
				h.Delete(0, at)
			}
		})
	}
}

func BenchmarkInsertDeleteBoxCache(b *testing.B) {
	dst := image.NewRGBA(image.Rect(0, 0, 1024, 768))
	for _, v := range []benchOp{
		{"1", "a", 0},
		{"10", strings.Repeat("a", 10), 0},
		{"100", strings.Repeat("a", 100), 0},
		{"1000", strings.Repeat("a", 1000), 0},
		{"10000", strings.Repeat("a", 10000), 0},
		{"100000", strings.Repeat("a", 100000), 0},
	} {
		b.Run(v.name, func(b *testing.B) {
			h := New(dst, dst.Bounds(), &Config{
				Font:  NewCliche(NewGoMono(fsize)),
				Color: A,
			})
			bb := []byte(v.data)
			b.SetBytes(int64(len(bb)))
			at := v.at
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Insert(bb, at)
				h.Delete(0, at)
			}
		})
	}
}
