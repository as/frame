package frame

import (
	"bytes"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
)

// Redraw redraws the entire frame. The caller should check
// that the frame is Dirty before calling this in a tight
// loop
func (f *Frame) Redraw(selecting bool, mouse image.Point ) {
	dy := f.alignY(f.origin).Y
	draw.Draw(f.disp, f.Bounds(), f.Colors.Back, image.ZP, draw.Src)
	for s := f.s[:f.nbytes]; ; {
		i := len(s)
		if i == 0 {
			break
		}
		if i > f.Wrap {
			i = f.Wrap
		}
		j := bytes.Index(s[:i], []byte("\n"))
		if j >= 0 {
			i = j
		}
		f.drawtext(image.Pt(f.origin.X, dy), s[:i])
		dy += f.FontHeight()
		if j >= 0 {
			i++
		}
		if i == len(s) {
			break
		}
		s = s[i:]
	}
	if selecting{
		e := f.IndexOf(mouse)
		f.Tick.Sweep(e)
	}
	f.Tick.Draw()
	f.Dirty = false
}

// drawtext draws the slice s at position p and returns
// the horizontal displacement dx without line wrapping
func (f *Frame) drawtext(pt image.Point, s []byte) (dx int) {
	return f.stringbg(f.disp, pt, f.Colors.Text, image.ZP, f.Font, s, f.Colors.Text, image.ZP)
}

func (f *Frame) stringbg(dst draw.Image, p image.Point, src image.Image, sp image.Point, font font.Face, s []byte, bg image.Image, bgp image.Point) int {
	for _, v := range s {
		fp := fixed.P(p.X,p.Y)
		dr, mask, maskp, advance, ok := font.Glyph(fp, rune(v))
		if !ok {
			break
		}
		h := f.FontHeight()
		dr.Min.Y += int(float64(h) - float64(h)/float64(5))
		dr.Max.Y += int(float64(h) - float64(h)/float64(5))
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over) 
		p.X += int((advance + f.Font.Kern(f.last, rune(v))) >> 6)
		f.last = rune(v)
	}
	return int(p.X)
}



// drawsel draws a highlight over points p through q. A highlight
// is a rectanguloid over three intersecting rectangles representing
// the highlight bounds.
func (t *Tick) drawsel(p, q image.Point, bg image.Image){
	h := t.Fr.FontHeight()
	m := t.Fr.r.Max
	o := t.Fr.origin

	// selection spans the same line
	if p.Y == q.Y{
		t.draw(p.X, p.Y, q.X, p.Y+h, bg)
		return
	}
	
	// draw up to three rectangles for the
	// selection
	
	t.draw(p.X, p.Y, m.X, p.Y+h, bg)
	p.Y += h
	if p.Y != q.Y {
		t.draw(o.X, p.Y, m.X, q.Y, bg)
	}
	t.draw(o.X, p.Y, q.X, q.Y+h, bg)
}

func (t *Tick) fill(x, y, xx, yy int){
	t.drawrect(image.Pt(x,y), image.Pt(xx,yy))
}
func (t *Tick) unfill(x, y, xx, yy int){
	t.deleterect(image.Pt(x,y), image.Pt(xx,yy))
}

func abs(x int) int{
	if x <0{
		return -x
	}
	return x
}


// drawrect draws a rectangle over the glyphs p0:p1
func (t *Tick) draw(x, y, xx, yy int, bg image.Image){
	r := image.Rect(x, y, xx, yy)
	draw.Draw(t.Img, r, bg, image.ZP, draw.Src)
}

// drawrect draws a rectangle over the glyphs p0:p1
func (t *Tick) drawrect(pt0, pt1 image.Point){
	r := image.Rect(pt0.X, pt0.Y, pt1.X, pt1.Y)
	draw.Draw(t.Img, r, t.Fr.Colors.HBack, image.ZP, draw.Src)
}

// delete draws a rectangle over the glyphs p0:p1
func (t *Tick) deleterect(pt0, pt1 image.Point){
	r := image.Rect(pt0.X, pt0.Y, pt1.X, pt1.Y)
	draw.Draw(t.Img, r, image.Transparent, image.ZP, draw.Src)
}
