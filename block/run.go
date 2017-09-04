package block

import (
	"fmt"
	"io"
)

const Size = 4096

type Run struct {
	Nchars int64
	Nalloc int
	Nblock int
	Block  []Block

	cache map[int64]local
}

type local struct {
	bn int
	i  int
}

func NewRun() *Run {
	return &Run{
		Block: []Block{
			{
				Ptr: make([]byte, Size),
			},
		},
		Nalloc: 1,
		Nblock: 1,
		cache:  make(map[int64]local),
	}
}

// Count recomputes and returns the number of bytes
// stored between box nb and the last block
func (f *Run) Count(nb int) int64 {
	n := int64(0)
	for ; nb < f.Nblock; nb++ {
		n += int64((&f.Block[nb]).Len())
	}
	return n
}

// Reset resets all boxes in the run without deallocating
// their data on the heap. If widthfn is not nill, it
// becomes the new measuring function for the run. Boxes
// in the run are not remeasured upon reset.
func (f *Run) Reset(func([]byte) int) {
	f.Nblock = 0
	f.Nchars = 0
}

func (f *Run) find(bn int, p, q int64) (int, int) {
	for ; bn < f.Nblock; bn++ {
		b := &f.Block[bn]
		if p+int64(b.Len()) > q {
			break
		}
		p += int64(b.Len())
	}
	return bn, int(q - p)
}

//Find finds the box containing q starting from box bn index
// p and puts q at the start of the next box
func (f *Run) Find(bn int, p, q int64) int {
	//	fmt.Printf("find %d.%d -> %d\n",bn,p,q)
	bn, i := f.find(bn, p, q)
	if i != 0 {
		f.Split(bn, i)
		bn++
	}
	return bn
}

func (f *Run) Dump() {
	fmt.Println("dumping boxes")
	fmt.Printf("nboxes: %d\n", f.Nblock)
	fmt.Printf("nalloc: %d\n", f.Nalloc)
	for i, b := range f.Block {
		fmt.Printf("[%d] (%p) (nrune=%d): %q\n", i, &f.Block[i], b.Nrune, b.Bytes())
	}
}

// Merge merges box bn and bn+1
func (f *Run) Merge(bn int) {
	b0 := &f.Block[bn]
	b1 := &f.Block[bn+1]
	b0.Ptr = append(b0.Ptr[:b0.Len()], b1.Ptr...)
	b0.Nrune += b1.Nrune
	f.Delete(bn+1, bn+1)
}

// Split splits box bn into two boxes; bn and bn+1, at index n
func (f *Run) Split(bn, n int) {
	f.Dup(bn)
	f.Truncate(&f.Block[bn], (&f.Block[bn]).Nrune-n)
	f.Chop(&f.Block[bn+1], n)
}

// Chop drops the first n chars in box b
func (f *Run) Chop(b *Block, n int) {
	//	fmt.Printf("Chop %q at %d\n", b.ptr, n)
	if b.Nrune < 0 || b.Nrune < n {
		panic("Chop")
	}
	copy(b.Ptr, b.Ptr[n:])
	b.Nrune -= n
	b.Ptr = b.Ptr[:b.Nrune]
}

func (f *Run) Truncate(b *Block, n int) {
	if b.Nrune < 0 || b.Nrune < n {
		panic("Truncate")
	}
	b.Nrune -= n
	b.Ptr = b.Ptr[:b.Nrune]
}

// Add adds n boxes after box bn, the rest are shifted up
func (f *Run) Add(bn, n int) {
	if bn > f.Nblock {
		panic("Frame.Add")
	}
	if f.Nblock+n > f.Nalloc {
		f.Grow(n + SLOP)
	}
	for i := f.Nblock - 1; i >= bn; i-- {
		f.Block[i+n] = f.Block[i]
	}
	f.Nblock += n
}

// Delete closes and deallocates n0-n1 inclusively
func (f *Run) Delete(n0, n1 int) {
	if n0 >= f.Nblock || n1 >= f.Nblock || n1 < n0 {
		panic("Delete")
	}
	f.Free(n0, n1)
	f.Close(n0, n1)
}

// Free deallocates memory for boxes n0-n1 inclusively
func (f *Run) Free(n0, n1 int) {
	if n1 < n0 {
		return
	}
	if n0 >= f.Nblock || n1 >= f.Nblock {
		panic("Free")
	}
	for i := n0; i < n1; i++ {
		if f.Block[i].Nrune >= 0 {
			f.Block[i].Ptr = nil
		}
	}
}

// Grow allocates memory for delta more boxes
func (f *Run) Grow(delta int) {
	f.Nalloc += delta
	f.Block = append(f.Block, make([]Block, delta)...)
}

// Dup copies the contents of box bn to box bn+1
func (f *Run) Dup(bn int) {
	f.Add(bn, 1)
	if f.Block[bn].Nrune >= 0 {
		f.Block[bn+1].Ptr = append([]byte{}, f.Block[bn].Ptr...)
	}
}

// Close closess box n0-n1 inclusively. The rest are shifted down
func (f *Run) Close(n0, n1 int) {
	if n0 >= f.Nblock || n1 >= f.Nblock || n1 < n0 {
		panic("Frame.Close")
	}
	n1++
	for i := n1; i < f.Nblock; i++ {
		f.Block[i-(n1-n0)] = f.Block[i]
	}
	f.Nblock -= n1 - n0
}

func (f *Run) WriteAt(p []byte, off int64) (n int, err error) {
	bn := f.Find(0, 0, off)
	bn0 := bn
	for len(p) != 0 {
		if bn == f.Nblock {
			f.Add(bn, 1)
		}
		println("0")
		b := &f.Block[bn]
		if len(b.Ptr) != Size {
			b.Ptr = make([]byte, 4096)
		}
		m := copy(b.Ptr[:], p)
		if m <= 0 {
			break
		}
		n += m
		b.Nrune += m
		fmt.Printf("loop %q\n", p)
		p = p[m:]
		fmt.Printf("loop %q\n", p)
		bn++
	}
	for i := bn0; i+1 < bn; i++ {
		f.Merge(i)
	}
	//	fmt.Printf("find %d.%d -> %d = box %d\n",bn,p,q, bn)
	return
}

func (f *Run) ReadAt(p []byte, off int64) (n int, err error) {
	//fmt.Printf("find %d.%d -> %d\n",bn,p,q)
	bn, i := f.find(0, 0, off)

	for len(p) != 0 {
		if bn == f.Nblock {
			fmt.Printf("bn == f.Nblock")
			return 0, io.EOF
		}
		b := &f.Block[bn]
		fmt.Printf("read: b.Ptr=%q\n", b.Bytes())
		m := copy(p, b.Ptr[i:b.Nrune])
		i = 0
		if m <= 0 {
			break
		}
		n += m
		p = p[m:]
		bn++
	}
	//	fmt.Printf("find %d.%d -> %d = box %d\n",bn,p,q, bn)
	return
}

func (b Run) String() string {
	s := ""
	bn, Nbox := 0, b.Nblock
	for ; bn < Nbox; bn++ {
		b := &b.Block[bn]
		s += string(b.Ptr)
	}
	return s
}
