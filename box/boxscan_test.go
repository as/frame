package box

import (
	"bytes"
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
)

var gomonoTTF, _ = truetype.Parse(gomono.TTF)

func NewGoMono(size int) font.Face {
	return truetype.NewFace(gomonoTTF, &truetype.Options{
		SubPixelsX:        64,
		SubPixelsY:        64,
		GlyphCacheEntries: 32768,
		Hinting:           font.HintingFull,
		Size:              float64(size),
	})
}

var fsize = 11

func genBench(b *testing.B, in []byte, min, max int, fn func(int) font.Face, ftsize int, bxceil int) {
	b.Helper()
	b.SetBytes(int64(len(in)))
	r := NewRun(min, max)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Boxscan(in, bxceil)
	}
}

func roll(size int) []byte {
	b := new(bytes.Buffer)
	for i := 0; i < size; i++ {
		b.WriteByte(byte(i))
	}
	return b.Bytes()
}

func BenchmarkScanByte(b *testing.B) { genBench(b, []byte("a"), 5, 5000, NewGoMono, 11, 1) }
func BenchmarkScan16Bytes(b *testing.B) {
	genBench(b, []byte("The quick brown "), 5, 5000, NewGoMono, 11, 1)
}
func BenchmarkScanHelloWorld(b *testing.B) {
	genBench(b, []byte(`package main\nimport "fmt"\n\nfunc main(){\n\tfmt.Println("hello world")\n}\n\n`), 5, 5000, NewGoMono, 11, 1)
}
func Benchmark100000Lines(b *testing.B) {
	genBench(b, bytes.Repeat([]byte{'\n'}, 100000), 5, 5000, NewGoMono, 8, 100000)
}
func Benchmark100000Lines2Byte(b *testing.B) {
	genBench(b, bytes.Repeat([]byte{'a', '\n'}, 100000), 5, 5000, NewGoMono, 16, 100000)
}
func Benchmark100000Lines4Byte(b *testing.B) {
	genBench(b, bytes.Repeat([]byte{'a', 'a', 'a', '\n'}, 100000), 5, 5000, NewGoMono, 16, 100000)
}
func Benchmark100000Lines16Byte(b *testing.B) {
	genBench(b, bytes.Repeat([]byte{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', '\n'}, 100000), 5, 5000, NewGoMono, 16, 100000)
}
func BenchmarkScanBinary100(b *testing.B)  { genBench(b, roll(100), 5, 5000, NewGoMono, 11, 100) }
func BenchmarkScanBinary1000(b *testing.B) { genBench(b, roll(1000), 5, 5000, NewGoMono, 11, 1000) }
func BenchmarkScanBinary5000(b *testing.B) { genBench(b, roll(5000), 5, 5000, NewGoMono, 11, 5000) }
func BenchmarkScanBinary10000(b *testing.B) {
	genBench(b, roll(10000), 5, 5000, NewGoMono, 11, 10000)
}
func BenchmarkScanBinary100000(b *testing.B) {
	genBench(b, roll(100000), 5, 5000, NewGoMono, 11, 100000)
}
func BenchmarkLongLine100000(b *testing.B) {
	genBench(b, bytes.Repeat([]byte{'a'}, 100000), 5, 5000, NewGoMono, 8, 100000)
}
