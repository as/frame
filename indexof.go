package frame

import (
	"fmt"
	"image"
)

// IndexOf returns the chracter index under the
// point pt.
func (f *Frame) IndexOf(pt image.Point) (p int64) {
	defer func() { fmt.Printf("indexof %s is %v\n", pt, p) }()
	pt = f.grid(pt)
	qt := f.r.Min
	bn := 0
	fmt.Println()
	println := func(s string) {
		fmt.Printf("%s		p=%d bn=%d qt=%s pt=%s\n", s, p,bn, qt, pt)
	}
	for ; bn < f.Nbox && qt.Y < pt.Y; bn++ {
		println("-> 0")
		b := &f.Box[bn]
		qt = f.wrapMax(qt, b)
		if qt.Y >= pt.Y {
			println("-> 1")
			break
		}
		qt = f.advance(qt, b)
		p += int64(b.Len())
		println("-> 2")
	}

	for ; bn < f.Nbox && qt.X <= pt.X; bn++ {
		println("---> 3")
		b := &f.Box[bn]
		qt = f.wrapMax(qt, b)
		if qt.Y > pt.Y {
			println("---> 4b")
			break
		}
		println("---> 4a")
		if qt.X+b.Width > pt.X {
			println("---> 5a")
			if b.Nrune < 0 {
				println("---> 6a")
				qt = f.advance(qt, b)
			} else {
				println("---> 6b")
				left := pt.X-qt.X
				fmt.Printf("fits at %s = %v\n", pt.X-qt.X, f.Font.Fits(b.Ptr, left ))
				p += int64(f.Font.Fits(b.Ptr, left))
				qt.X += left
				
				
				//				bs := f.newRulerFunc(b.Ptr, f.Font)
				//				for {
				//					size, width, err := bs.Next()
				//					if err != nil {
				//						break
				//					}
				//					qt.X += width
				//					if qt.X > pt.X {
				//						break
				//					}
				//					p += int64(size)
				//				}
			}
		} else {
			println("---> 5b")
			p += int64(b.Len())
			qt = f.advance(qt, b)
		}
	}
	println("---> ret")
	return p
}
