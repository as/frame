package frame
// TODO(as): seperate this into its own package for testing graphics

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"

	"github.com/as/frame/font"
)

var (
	R = image.Rect(0, 0, 232, 232)
	fsize = 16
	ft = font.NewBasic(fsize)
	enc = png.Encoder{CompressionLevel: png.NoCompression}

	red = image.NewUniform(color.RGBA{255, 0, 0, 255})
	blue = image.NewUniform(color.RGBA{0, 0, 255, 255})

	black = image.NewUniform(color.RGBA{0, 0, 0, 255})
	white = image.NewUniform(color.RGBA{255, 255, 255, 255})
	gray = image.NewUniform(color.RGBA{33, 33, 33, 255})
)

func tofile(t *testing.T, file string, img *image.RGBA) {
	fd, err := os.Create(file)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer fd.Close()
	err = enc.Encode(fd, img)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func equal(c0, c1 color.Color) bool {
	r0, g0, b0, _ := c0.RGBA()
	r1, g1, b1, _ := c1.RGBA()
	return r0 == r1 && g0 == g1 && b0 == b1
}

func bg(c color.Color) bool {
	return c.(color.RGBA) == color.RGBA{255, 248, 232, 255} // Peach
}

func Eval(have, want image.Image) (ok bool, delta *image.RGBA) {
	delta = image.NewRGBA(have.Bounds())
	for y := have.Bounds().Min.Y; y < have.Bounds().Max.Y; y++ {
		for x := have.Bounds().Min.X; x < have.Bounds().Max.X; x++ {
			h := have.At(x, y)
			w := want.At(x, y)
			if equal(h, w) {
				delta.Set(x, y, white)
				continue
			}
			if bg(h) {
				delta.Set(x, y, red)
			} else {
				delta.Set(x, y, blue)
			}
		}
	}
	return delta.Opaque(), delta
}

func report(have, want, delta image.Image) *image.RGBA {
	r := have.Bounds()
	r.Max.X = r.Min.X + r.Dx()*3 + 5*4
	rep := image.NewRGBA(r)
	draw.Draw(rep, r, gray, rep.Bounds().Min, draw.Src)
	r.Min.X += 5
	for _, src := range []image.Image{have, want, delta} {
		drawBorder(rep, r.Inset(-1), black, image.ZP, 2)
		draw.Draw(rep, r, src, src.Bounds().Min, draw.Src)
		r.Min.X += want.Bounds().Dx() + 5
	}

	return rep
}

func abtestPad16(r image.Rectangle) (fr0, fr1 *Frame, a, b *image.RGBA) {
	a = image.NewRGBA(r)
	b = image.NewRGBA(r)
	fr0 = New(r.Inset(fsize), font.NewBasic(fsize), a, A)
	fr1 = New(r.Inset(fsize), font.NewBasic(fsize), b, A)
	return fr0, fr1, a, b
}
func abtest(r image.Rectangle) (fr0, fr1 *Frame, a, b *image.RGBA) {
	a = image.NewRGBA(r)
	b = image.NewRGBA(r)
	fr0 = New(r, font.NewBasic(fsize), a, A)
	fr1 = New(r, font.NewBasic(fsize), b, A)
	return fr0, fr1, a, b
}

func failimg(t *testing.T, file string, have, want, delta image.Image) {
	tofile(t, file, report(have, want, delta))
	t.Log("images differ")
	t.Fail()
}

func TestDeleteLastLineNoNL(t *testing.T) {
	w, h, want, have := abtest(R)
	draw.Draw(want, want.Bounds(), w.Color.Back, image.ZP, draw.Src)
	draw.Draw(have, have.Bounds(), h.Color.Back, image.ZP, draw.Src)
	w.Insert([]byte("1234\ncccc\ndddd\n"), 0)
	h.Insert([]byte("1234\ncccc\ndddd"), 0)
	
	h.Delete(5, 10)
	w.Delete(5, 10)
	// We can untick because have has an extra newline
		h.Untick()
		w.Untick()	
	ok, delta := Eval(have, want)
	if !ok {
		failimg(t, "delta.png", want, have, delta)
		h.DumpBoxes()
		w.DumpBoxes()
	}
}
