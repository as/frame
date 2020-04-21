package spaz

import (
	"io"
	"math/rand"
	"time"
)

type Reader struct{
	ur io.Reader
	rs *rand.Rand
}

func NewReader(r io.Reader) *Reader{
	return &Reader{ur: r, rs: rand.New(rand.NewSource(time.Now().Unix()))}
}
func (r *Reader) intn(n int)int{
	return int(r.rs.Int31n(int32(n)))
}

func (r *Reader) Read(p []byte) (n int, err error) {
	rn := r.intn(len(p))
	return r.ur.Read(p[:rn])
}
