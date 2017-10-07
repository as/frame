package box

import (
//	"fmt"
)

// Elastic tabstop experiment section. Not ready for general use by any means
// the text/tabwriter package implements elastic tabstops, but this package
// assumes that all chars are the same width and that text needs to be rescanned.
//
// The frame already distinguishes between tabs, newlines, and plain text characters
// by encapsulating them in measured boxes. A direct copy of the tabwriter code would
// ignore the datastructures in the frame and their sublinear runtime cost.
//
// A few cases need to be handles for elastic tabstops:
//
// Insertion:
// 	An insert always creates a new box, this box can take on the form '\n', '\t'
//  or be part of the plain text character set which for this purpose is the set
//  of characters not equal to '\t' or '\n'. A successful insert operation runs
//  to completion and results in a net gain of runes (not taking into consideration
//  the runes which run off the frame and are deleted). An insert never subtracts
//  of boxes assuming that all characters stay on the frame. There are now a few
//  cases:
//
//  1. The insert contains plain text characters
//		This is already a non-trivial case. Plaintext characters can move a column of
//		text far enough to extend the range of the column's tab stop. The procedure then
//		is to seek forth and determine if a tab follows the insertion point. This can be
//		completed in sub-linear time, as plain text boxes are merged into one before the
//		start of a hard newline or soft wrap (when a box width exceeds frame.r.Max). The
//		soft line wrap property of boxes may aggregate the runtime performance of this
//		implementation by containing sparse newlines. In this case we consider an algorithm
//		that gives up and favors performance over elasticity. Elastic tabstops have little
//		use in binary files anyway and we probably want to give users the opportunity to turn
//		them off. The seek forth operation scans forward and locates the first newline, while
//		counting the occurrence of \t boxes and measuring their width. The maximum width column
//		is tracked as well. The seek forth operation terminates whenever a column-less line is
//		located (formally defined as a run of boxes separated by a newline containing no \t boxes
//


func (f *Run) Stretch(nb int) {
	//	fmt.Println()
	//	fmt.Printf("\n\ntrace bn=%d\n", nb)
	nc := 0
	nl := 0
	dx := 0

	cmax := make(map[int]int)
	cbox := make(map[int][]int)
	
	nb = f.FindCell(nb)
	Loop:
	for ; nb < f.Nbox; nb++ {
		switch b := &f.Box[nb]; b.BC{
		case '\t':
			dx += b.Width
			cbox[nc] = append(cbox[nc], nb)
			max := cmax[nc]
			if dx > max {
				cmax[nc] = dx
			}
			nc++
			dx = 0
		case '\n': 
			nl++
			dx = 0
			if nc == 0 {
				// A line with no tabs; end of cell
				break Loop
			}
			nc = 0
		default:
			dx += b.Width
		}
	}
	for c, bns := range cbox {
		max := cmax[c]
		for _, bn := range bns {
			//dx := f.Box[bn].Width
			if c == 0{
				f.Box[bn].Width = max
			} else {
				f.Box[bn].Width = max - f.Box[bn-1].Width
			}
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

// FindCell returns the first box in the cell
func (f *Run) FindCell(bn int) int{
	if bn == 0{
		return 0
	}
	ncols := 0
	nrows := 0
	oldbn := bn
	bn = f.EndLine(bn)
	b := &f.Box[bn]
	for bn-1 != 0 {
		b = &f.Box[bn-1]
		switch ; b.BC{
		case '\n':

			if ncols == 0{
						if nrows == 0{
				return oldbn
			}
				return bn+1
			}
			nrows++
			ncols=0
		case '\t':
			ncols++
		default:
		}
		bn--
	}
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

func (f *Run) EndLine(bn int) int{
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
	if bn < f.Nbox{
		return bn+1
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
