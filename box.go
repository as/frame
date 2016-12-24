package frame

import "io"
import "fmt"

type Box struct{
	data []byte
	width int
}

type Boxes struct{
	Box []*Box
	measure func([]byte) int
}

func NewBoxes(measure func([]byte) int) *Boxes{
	b := &Boxes{measure: measure}
	b.Add(10)
	return b
}

func (b *Boxes) WriteAt(p []byte, off int64) (n int, err error){
	b.Insert(p, int(off))
	return len(p), nil
}

func (b *Boxes) ReadAt(p []byte, off int64) (n int, err error){
	bn, err := b.Find(0, 0, int(off))
	fmt.Println("offset %d found in box %d", off, bn)
	if err != nil{
		return n,err
	}
	bp := b.Box[bn]
	for need := len(p); need > 0; {
		m := copy(p[n:], bp.data[:min(len(bp.data), need)])
		n += m
		need -= m
		bn++
		if bn >= len(b.Box){
			return n, io.ErrUnexpectedEOF
		}
		bp = b.Box[bn]
	}
	return n, nil
}


func (b *Boxes) Insert(p []byte, off int){
	n, err := b.Find(0,0,off)
	if err != nil{
		panic("insert")
	}
	if n >= len(b.Box){
		b.Add(1)
		n = len(b.Box)-1
	}
	bp := b.Box[n]
	bp.data = append([]byte{}, p...)
	return
	if len(bp.data) < len(p){
		bp.data = append([]byte{}, p...)
	} else {
		copy(bp.data, p)
		bp.data = bp.data[:len(p)]
	}
}

func (b *Boxes) Find(n, i, j  int) (int, error){
	for _, v := range b.Box{
		w := len(v.data)
		if w+i > j{
			break
		}
		i += w
	}
	if i != j {
		b.Split(n+1, j-i)
	}
	return n+1, nil
}

func (b *Boxes) Add(n int){
	boxes := make([]*Box, n)
	for i := range boxes{
		boxes[i] = &Box{data: make([]byte, 1)}
	}
	b.Box = append(b.Box, boxes...)
}

func (b *Boxes) Dup(n int) (*Box){
	bp0 := b.Box[n]
	bp1 := &Box{
		width: bp0.width,
		data: append([]byte{}, bp0.data...),
	}
	b.Box = append(b.Box[:n+1], append([]*Box{bp1}, b.Box[n+1:]...)...)
	fmt.Printf("%#v\n", b.Box)
	return bp1
}

func (b *Boxes) Split(n int, at int){
	box := b.Dup(n)
	b.Truncate(n, len(box.data)-at)
	b.Chop(n+1, at)
}

func (b *Boxes) Truncate(n, at int){
	box := b.Box[n]
	box.data = box.data[:at]
	box.width = b.measure(box.data)
}

func (b *Boxes) Chop(n int, at int){
	box := b.Box[n]
	box.data = box.data[at:]
	box.width = b.measure(box.data)
}

func (b *Boxes) Merge(n int){
	sp := b.Box[n:]
	sp[0].data = append(sp[0].data, sp[1].data...)
	sp[0].width += sp[1].width
}

func (b *Boxes) Delete(n0, n1 int){
	dn := len(b.Box[n1:]) - len(b.Box[n0:])
	copy(b.Box[n0:], b.Box[n1:])
	b.Box = b.Box[:len(b.Box)-dn]
}

func (b *Box) Len() int{
	return len(b.data)
}

func (b *Box) Cap() int{
	return cap(b.data)
}

func min(a,b int)int {
	if a<b{
		return a
	}
	return b
}