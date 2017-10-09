package box

import (
	"fmt"
)

// Elastic tabstop experiment section. Not ready for general use by any means
// the text/tabwriter package implements elastic tabstops, but that package
// assumes that all chars are the same width and that text needs to be rescanned.
//
// The frame already distinguishes between tabs, newlines, and plain text characters
// by encapsulating them in measured boxes. A direct copy of the tabwriter code would
// ignore the datastructures in the frame and their sublinear runtime cost.
//
func (f *Run) Stretch(nb int) (pb int) {
	if nb <= 0 {
		return 0
	}
	//	fmt.Println()
	//	fmt.Printf("\n\ntrace bn=%d\n", nb)
	nc := 0
	nl := 0
	dx := 0

	cmax := make(map[int]int)
	cbox := make(map[int][]int)

	nb = f.FindCell(nb)
	fmt.Println("\n\ncell start at box", nb)
	pb = nb - 1
Loop:
	for ; nb < f.Nbox; nb++ {
		b := &f.Box[nb]
		fmt.Printf("switch box: %#v\n", b)
		switch b.BC {
		case '\t':
			dx += b.Width
			cbox[nc] = append(cbox[nc], nb)
			max := cmax[nc]
			if dx > max {
				cmax[nc] = dx
			}
			nc++
			fmt.Printf("	tab: dx=%d ncol=%d\n", dx, nc)
			dx = 0
		case '\n':
			nl++
			dx = 0
			if nc == 0 {
				// A line with no tabs; end of cell
				fmt.Printf("	nl (no cols): dx=%d nl=%d\n", dx, nl-1)
				break Loop
			}
			fmt.Printf("	nl : dx=%d nl=%d nc=%d\n", dx, nl-1, nc)
			nc = 0
		default:
			dx += b.Width
			fmt.Printf("	plain : dx=%d wid=%d nc=%d\n", dx, b.Width, nc)
		}
	}
	for c, bns := range cbox {
		max := cmax[c]
		for _, bn := range bns {
			b := &f.Box[bn]
			b.Width = max
			if bn == 0 {
				continue
			}
			pb := f.Box[bn-1]
			if pb.BC != '\n' {
				b.Width -= f.Box[bn-1].Width
			}
			if b.Width < b.Minwidth {
				b.Width = b.Minwidth
			}
		}
	}
	return pb
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

// FindCell returns the first box in the cell
func (f *Run) FindCell(bn int) int {
	if bn == 0 {
		return 0
	}
	ncols := 0
	nrows := 0
	oldbn := bn
	bn = f.EndLine(bn)
	b := &f.Box[bn]
	for bn-1 != 0 {
		b = &f.Box[bn-1]
		switch b.BC {
		case '\n':

			if ncols == 0 {
				if nrows == 0 {
					return oldbn
				}
				return bn + 1
			}
			nrows++
			ncols = 0
		case '\t':
			ncols++
		default:
		}
		bn--
	}
	//	println("bn-1", bn-1)
	//	f.DumpBoxes()
	if bn-1 == 0 && f.Box[bn-1].BC != '\n' {
		return 0
	}
	//	println("return", bn)
	return bn
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

func (f *Run) EndLine(bn int) int {
	for bn < f.Nbox {
		b := &f.Box[bn]
		if b.BC == '\n' {
			break
		}
		bn++
	}
	return bn
}

func (f *Run) NextLine(bn int) int {
	bn = f.EndLine(bn)
	if bn < f.Nbox {
		return bn + 1
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
