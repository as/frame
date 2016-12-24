package frame

import "io"
import "image"
import "image/draw"

type Tick struct {
	P0, P1 int
	Fr     *Frame
	Select
	dirty bool
}

// Selection consists of 3 rectangles
//
// r0 is non-zero when the selection begins at a non-zero X coordinate
//
// r1 is non-zero if r0 is non-zero and r2 terminates
// the selection on an x-coordinate below max X
//
// r2 always terminates the selection. 
type Select struct{
	Img *image.RGBA	// todo: use New
	
	// anchor index and index of last sweep	
	a int
	b int
	
	// cache of points for a,b
	at image.Point
	bt image.Point
}

func (t *Tick) Close() error{
	draw.Draw(t.Img, t.Fr.disp.Bounds(), image.Transparent, image.ZP, draw.Src)
	t.a, t.b = 0, 0
	t.at, t.bt = image.ZP, image.ZP
	return nil
}

func (t *Tick) Draw() error{
	if t.P1 == t.P0{
		pt := t.Fr.PointOf(t.P1)
		r := image.Rect(0, 0, 2, t.Fr.FontHeight()).Add(pt)
		draw.Draw(t.Fr.disp, r, t.Fr.Colors.Text, image.ZP, draw.Over)
	} else {
		// assuming the underlying selection is already
		// drawn on t.Img, see Sweep.
		draw.Draw(t.Fr.disp, t.Fr.Bounds(), t.Img, image.ZP, draw.Over)
	}
	return nil
}

func (t *Tick) SelectAt(i int) {
	t.a = i
	t.at = t.Fr.PointOf(i)
	t.b = i
}

func (t *Tick) Sweep(j int){
	a := t.a
	b := t.b
	c := j
	if abs(b-c) < 1{
		return
	}
	bg, erase := t.Fr.Colors.HBack, image.Transparent
	pt := t.Fr.PointOf
	at, bt, ct := pt(a), pt(b), pt(c)
	switch{
	case a <= b && b < c:
		// down
		t.drawsel(bt, ct, bg)
	case a <= c && c < b: 
		// down and up
		t.drawsel(ct, bt, erase)
	case c < a && a <= b: 
		// down and up over
		t.drawsel(at, bt, erase)
		t.drawsel(ct, at, bg)
	case c < b && b <= a: 
		// up
		t.drawsel(ct, bt, bg)
	case b < c && c < a: 
		// up and down
		t.drawsel(bt, ct, erase)
	case b < a && a < c: 
		// up and down over
		t.drawsel(bt, at, erase)
		t.drawsel(at, ct, bg)
	}
	
	t.b = c

}

func (t *Tick) Insert(p []byte) (err error) {
	if len(p) == 0 {
		return nil
	}
	if t.P1 != t.P0 {
		t.Delete()
	}
	if err = t.Fr.Insert(p, t.P1); err != nil {
		return err
	}
	t.P0 += len(p)
	t.P1 += len(p)
	return nil
}

func (t *Tick) Delete() (err error) {
	if t.P0 == t.P1 && t.P0 == 0 {
		return nil
	}
	// Either act like the delete button or erase the
	// contents of an active selectiond
	if t.P0 == t.P1 {
		t.P0--
		t.Fr.Delete(t.P0, t.P1)
		t.P1--
	} else {
		if t.P0 > t.P1 {
			t.P0, t.P1 = t.P1, t.P0
		}
		t.Fr.Delete(t.P0, t.P1)
		t.P1 = t.P0
	}
	return nil
}

func (t *Tick) WriteRune(r rune) (err error) {
	return t.Insert([]byte{byte(r)})
}

func (t *Tick) ck() {
	nb := len(t.Fr.s)
	if t.P1 >= nb {
		t.P1 = nb - 1
	}
	if t.P0 >= nb {
		t.P0 = nb - 1
	}
	if t.P1 < 0 {
		t.P1 = 0
	}
	if t.P0 < 0 {
		t.P0 = 0
	}
}

func (t *Tick) String() string {
	t.ck()
	return string(t.Fr.s[t.P0:t.P1])
}

func (t *Tick) Read(p []byte) (n int, err error) {
	if t.P0 == t.P1 {
		return 0, io.EOF
	}
	q := t.Fr.s[t.P0:t.P1]
	return copy(p, q), nil
}

func (t *Tick) Write(p []byte) (n int, err error) {
	err = t.Insert(p)
	if err != nil {
		return 0, err
	}
	t.Fr.Dirty = true
	return len(p), nil
}
