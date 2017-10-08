package box

import (
	"bytes"
	"github.com/as/frame/font"
	"testing"
)

var fsize = 11

func BenchmarkScanByte(b *testing.B) {
	bb := []byte("a")
	r := NewRun(5, 5000, font.NewBasic(fsize))
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 1)
		b.StopTimer()
	}
}

func BenchmarkScanByteFont(b *testing.B) {
	bb := []byte("a")
	ft := font.NewGoRegular(8)
	r := NewRun(5, 5000, ft)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 1)
		b.StopTimer()
	}
}

func BenchmarkScan16Bytes(b *testing.B) {
	bb := []byte("The quick brown ")
	fn := font.NewBasic(fsize)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 1)
		b.StopTimer()
	}
}

func BenchmarkScan16BytesFont(b *testing.B) {
	bb := []byte("The quick brown ")
	ft := font.NewGoRegular(8)
	r := NewRun(5, 5000, ft)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 1)
		b.StopTimer()
	}
}
func BenchmarkScanHelloWorld(b *testing.B) {
	bb := []byte(`package main\nimport "fmt"\n\nfunc main(){\n\tfmt.Println("hello world")\n}\n\n`)
	fn := font.NewBasic(fsize)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 10)
		b.StopTimer()
	}
}
func BenchmarkScanHelloWorldFont(b *testing.B) {
	bb := []byte(`package main\nimport "fmt"\n\nfunc main(){\n\tfmt.Println("hello world")\n}\n\n`)
	ft := font.NewGoRegular(8)
	r := NewRun(5, 5000, ft)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 10)
		b.StopTimer()
	}
}

func roll(size int) []byte {
	b := new(bytes.Buffer)
	for i := 0; i < size; i++ {
		b.WriteByte(byte(i % 256))
	}
	return b.Bytes()
}

func BenchmarkScanBinary100(b *testing.B) {
	bb := roll(100)
	fn := font.NewBasic(fsize)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100)
		b.StopTimer()
	}

}
func BenchmarkScanBinary100Font(b *testing.B) {
	bb := roll(100)
	ft := font.NewGoRegular(8)
	r := NewRun(5, 5000, ft)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100)
		b.StopTimer()
	}

}

func BenchmarkScanBinary1000(b *testing.B) {
	bb := roll(100)
	fn := font.NewBasic(fsize)

	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100)
		b.StopTimer()
	}

}
func BenchmarkScanBinary1000Font(b *testing.B) {
	bb := roll(1000)
	ft := font.NewGoRegular(8)
	r := NewRun(5, 5000, ft)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 1000)
		b.StopTimer()
	}
}

func BenchmarkScanBinary5000(b *testing.B) {
	bb := roll(5000)
	fn := font.NewBasic(fsize)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100)
		b.StopTimer()
	}

}
func BenchmarkScanBinary5000Font(b *testing.B) {
	bb := roll(5000)
	fn := font.NewBasic(fsize)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 1000)
		b.StopTimer()
	}
}

func BenchmarkScanBinary100000(b *testing.B) {
	bb := roll(100000)
	fn := font.NewBasic(fsize)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}

}
func BenchmarkScanBinary100000Font(b *testing.B) {
	bb := roll(100000)
	fn := font.NewGoRegular(9)
	r := NewRun(5, 5000, fn)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}
}

func BenchmarkLongLine100000(b *testing.B) {
	bb := bytes.Repeat([]byte{'a'}, 100000)
	fn := font.NewGoRegular(8)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}

}
func BenchmarkLongLine100000Font(b *testing.B) {
	bb := bytes.Repeat([]byte{'a'}, 100000)
	ft := font.NewGoRegular(8)
	r := NewRun(5, 5000, ft)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}
}

func Benchmark100000Lines(b *testing.B) {
	bb := bytes.Repeat([]byte{'\n'}, 100000)
	fn := font.NewGoRegular(8)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}

}
func Benchmark100000LinesFont(b *testing.B) {
	bb := bytes.Repeat([]byte{'\n'}, 100000)
	ft := font.NewGoRegular(8)
	r := NewRun(5, 5000, ft)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}
}

func Benchmark100000Lines2Byte(b *testing.B) {
	bb := bytes.Repeat([]byte{'a', '\n'}, 100000)
	fn := font.NewGoRegular(16)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}

}
func Benchmark100000Lines2ByteFont(b *testing.B) {
	bb := bytes.Repeat([]byte{'a', '\n'}, 100000)
	ft := font.NewGoRegular(8)
	r := NewRun(5, 5000, ft)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}
}

func Benchmark100000Lines4Byte(b *testing.B) {
	bb := bytes.Repeat([]byte{'a', 'a', 'a', '\n'}, 100000)
	fn := font.NewGoRegular(8)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}

}
func Benchmark100000Lines4ByteFont(b *testing.B) {
	bb := bytes.Repeat([]byte{'a', 'a', 'a', '\n'}, 100000)
	ft := font.NewGoRegular(8)
	r := NewRun(5, 5000, ft)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}
}

func Benchmark100000Lines16Byte(b *testing.B) {
	bb := bytes.Repeat([]byte{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', '\n'}, 100000)
	fn := font.NewGoRegular(8)
	r := NewRun(5, 5000, fn)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}

}
func Benchmark100000Lines16ByteFont(b *testing.B) {
	bb := bytes.Repeat([]byte{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', '\n'}, 100000)
	ft := font.NewGoRegular(8)
	r := NewRun(5, 5000, ft)

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		r.Boxscan(bb, 100000)
		b.StopTimer()
	}
}
