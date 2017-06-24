package tag

import (
	"bytes"
	"fmt"
)

var (
	Lefts    = [...]byte{'(', '{', '[', '<', '"', '\'', '`'}
	Rights   = [...]byte{')', '}', ']', '>', '"', '\'', '`'}
	Free     = [...]byte{'"', '\'', '`'}
	AlphaNum = []byte("*&!%-_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func isany(b byte, s []byte) bool {
	for _, v := range s {
		if b == v {
			return true
		}
	}
	return false
}

func findback(p []byte, i int64, sep []byte) int64 {
	for ; i-1 >= 0 && isany(p[i-1], sep); i-- {
	}
	return i
}
func find(p []byte, j int64, sep []byte) int64 {
	for ; j != int64(len(p)) && isany(p[j], sep); j++ {
	}
	return j
}

func FindAlpha(p []byte, i int64) (int64, int64) {
	j := find(p, i, AlphaNum)
	i = findback(p, i, AlphaNum)
	return i, j
}

func FindNext(f File, text []byte) (q0, q1 int64){
	i, j := f.Dot()
	p := f.Bytes()
	x := text
	q0 = int64(bytes.Index(p[j:], x))
	if q0 == -1 {
		println("a")
		q0 = int64(bytes.Index(p[:i], x))
		if q0 < 0 {
			println("b")
			return i, j
		}
	} else {
		println("c")
		q0 += j
	}
	q1 = q0 + int64(len(x))
	println("d")
	return q0, q1
}

func Next(p []byte, i, j int64) (q0 int64, q1 int64) {
	defer func(r0, r1 int64) {
		fmt.Printf("Next: [%d:%d]->[%d:%d]\n", r0, r1, q0, q1)
	}(i, j)
	x := p[i:j]
	q0 = int64(bytes.Index(p[j:], x))
	if q0 == -1 {
		println("a")
		q0 = int64(bytes.Index(p[:i], x))
		if q0 < 0 {
			println("b")
			return i, j
		}
	} else {
		println("c")
		q0 += j
	}
	q1 = q0 + int64(len(x))
	println("d")
	return q0, q1
}

/*
func (t *Tick) Find(p []byte, back bool) int {
	return t.find(p, t.P1, len(t.Fr.s), back)
}

func (t *Tick) find(p []byte, i, j int, back bool) int {
	if back {
		panic("unimplemented")
	}
	//fmt.Printf("debug: find: %q check frame[%d:]\n", p, t.P1)
	x := bytes.Index(t.Fr.s[i:j], p)
	if x == -1 {
		return -1
	}
	println("found at index", i, ":", x+i)
	return x + i

}

func (t *Tick) FindSpecial(i int) (int, int) {
	fmt.Println("NUMBER", i)
	if i == 0 {
		return i, t.FindOrEOF([]byte{'\n'})
	}
	t.Open(i - 1)
	t.Sweep(i)
	t.Commit()
	if t.ReadByte() == '\n' {
		return i, t.FindOrEOF([]byte{'\n'})
	}
	if x := t.FindQuote(); x != -1 {
		return i, x
	}
	if x := t.FindParity(); x != -1 {
		return i, x
	}
	if isany(t.ReadByte(), AlphaNum) {
		return t.FindAlpha(i)
	}
	return i, -1
}

func (t *Tick) FindOrEOF(p []byte) int {
	i := t.Find(p, false)
	if i == -1 {
		return t.Fr.nbytes
	}
	return i
}

func (t *Tick) FindQuote() int {
	b := t.ReadByte()
	for _, v := range Free {
		if b != v {
			continue
		}
		return t.Find([]byte{v}, false)
	}
	return -1
}

func (t *Tick) FindParity() int {
	for i := range Lefts {
		j := t.findParity(Lefts[i], Rights[i], false)
		if j != -1 {
			return j
		}
	}
	return -1
}

func (t *Tick) findParity(l byte, r byte, back bool) int {
	if back {
		panic("unimplemented")
	}
	b := t.ReadByte()
	if b != l {
		return -1
	}
	push := 1
	//j := -1
	for i, v := range t.Fr.s[t.P1:] {
		if v == l {
			println("\n\n++\n\n")
			push++
		}
		if v == r {
			println("\n\n--\n\n")
			push--
			if push == 0 {
				return i + t.P1
			}
		}
	}
	return -1
}
*/
