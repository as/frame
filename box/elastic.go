package box

import (
//	"fmt"
)

// Elastic tabstop experiment section. Not ready for general use by any means

func (f *Run) Stretch(nb int) {
//	fmt.Println()
//	fmt.Printf("\n\ntrace bn=%d\n", nb)
	nc := 0
	nl := 0
	dx := 0

	cmax := make(map[int]int)
	cbox := make(map[int][]int)

	for ; nb < f.Nbox; nb++ {
		b := &f.Box[nb]
		if b.BC == '\n' {
			nl++
			nc = 0
			dx = 0
			continue
		}
		dx += b.Width
		if b.BC == '\t' {
			cbox[nc] = append(cbox[nc], nb)
			max := cmax[nc]
			if dx > max {
				cmax[nc] = dx
			}
//			fmt.Printf("line %d col %d width %d\n", nl, nc, dx)
			nc++
			dx = 0
		}
	}
	for c, bns := range cbox {
		max := cmax[c]
		for _, bn := range bns {
			//dx := f.Box[bn].Width
			f.Box[bn].Width = max - f.Box[bn-1].Width
		}
	}
}

func (f *Run) Findcol(bn int, coln int) (cbn int, xmax int) {
	c := 0
	for ; bn < f.Nbox; bn++ {
		b := &f.Box[bn]
		if b.BC == '\t' {
			c++
		}
		if b.BC != '\n' {
			xmax += b.Width
		}
		if c == coln {
			break
		}
		bn++
	}
	if c != coln {
		return -1, 0
	}
	return bn, xmax

}

func (f *Run) Colof(bn int) (coln, xmax int) {
	if bn == 0 {
		return 0, 0
	}
	bs := f.StartLine(bn)
	for {
		b := &f.Box[bs]
		if b.BC == '\t' {
			coln++
		}
		if b.BC != '\n' {
			xmax += b.Width
		}
		if bn == bs {
			break
		}
		bs++
	}
	if xmax != 0 {
		coln++
	}
	return coln, xmax
}

func (f *Run) StartLine(bn int) int {
	for ; bn-1 >= 0; bn-- {
		b := &f.Box[bn-1]
		if b.BC == '\n' {
			break
		}
	}
	return bn
}

func (f *Run) PrevLine(bn int) int {
	for ; bn >= 0; bn-- {
		b := &f.Box[bn]
		if b.BC == '\n' {
			break
		}
	}
	if bn == -1 && f.Box[0].BC == '\n' {
		return 0
	}
	for bn-1 >= 0 {
		b := &f.Box[bn-1]
		if b.BC == '\n' {
			break
		}
		bn--
	}
	return bn
}

func (f *Run) NextLine(bn int) int {
	for ; bn < f.Nbox; bn++ {
		b := &f.Box[bn]
		if b.BC == '\n' {
			return bn + 1
		}
		bn++
	}
	return bn
}
