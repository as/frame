package tag

import (
	"github.com/as/clip"
	"github.com/as/text"
	"io"
	"bytes"
)

var (
	ClipBuf = make([]byte, 1024*1024*16)
	Clip    *clip.Clip
)

func toUTF16(p []byte) (q []byte) {
	i := 0
	q = make([]byte, len(p)*2)
	for j := 0; j < len(p); j++ {
		q[i] = p[j]
		i += 2
	}
	return q
}
func fromUTF16(p []byte) (q []byte) {
	j := 0
	q = make([]byte, len(p)/2)
	for i := 0; i < len(q); i++ {
		q[i] = p[j]
		j += 2
	}
	return q
}

func snarf(w text.Editor) {
	n := copy(ClipBuf, toUTF16([]byte(Rdsel(w))))
	io.Copy(Clip, bytes.NewReader(ClipBuf[:n]))
	q0,q1:=w.Dot()
	w.Delete(q0,q1)
}

func paste(w text.Editor) {
	n, _ := Clip.Read(ClipBuf)
	s := fromUTF16(ClipBuf[:n])
	q0, q1 := w.Dot()
	if q0 != q1 && q1 > q0{
		w.Delete(q0, q1)
		w.Select(q0, q0)
		q1 = q0
	}
	
	//w.Insert(s, q0)
	w.Select(q0,q0+int64(w.Insert(s, q0)))
}
