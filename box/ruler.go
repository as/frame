package box

import (
	"errors"
	"io"
	"unicode/utf8"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	//	"log"
)

var (
	ErrZeroRead     = errors.New("zero length read")
	ErrDoubleUnread = errors.New("double unread")
)

type Ruler interface {
	Advance() error
	Bytes() []byte
	Next() (size, width int, err error)
	Last() []byte
	Len() int
	Unread() (size, width int, err error)
	Width() int
	Reset(p []byte)
}

type byteRuler struct {
	font.Face
	b                   []byte
	i                   int
	w                   int
	lastSize, lastWidth int
	sizetab             [256]int
}

func NewByteRuler(b []byte, ft font.Face) Ruler {
	br := &byteRuler{Face: ft, b: b}
	for i := range br.sizetab{
		x, _ := br.Face.GlyphAdvance(rune(i))
		br.sizetab[i] = Fix(x)
	}
	return br
}

func (bs *byteRuler) MeasureByte(b byte) int {
	return bs.sizetab[b]
}

func (bs *byteRuler) measureWidth(s []byte) (width int){
	for _, c := range s{
		width+=bs.sizetab[c]
	}
	return width
}

func (bs *byteRuler) Bytes() []byte { return bs.b[:bs.i] }

func (bs *byteRuler) MeasureWidth() (width int) {
	for _, c := range bs.b[:bs.i]{
		width+=bs.sizetab[c] // bs.MeasureByte(c)
	}
	return width
}
func (bs *byteRuler) Chop() []byte {
	b := bs.b[:bs.i]
	bs.b = bs.b[bs.i:]
	bs.i=0
	return b
}
func (bs *byteRuler) Advance() (err error) {
	if bs.i == len(bs.b) {
		return io.EOF
	}
	bs.b = bs.b[bs.i:]
	bs.i = 0
	return nil
}
func (bs *byteRuler) Run(){
	for _, c := range bs.b[bs.i:] {//bs.i != len(bs.b){
		if special(c){
			return
		}
		bs.i++
	}
}

func (bs *byteRuler) Scan() bool{
	if bs.i == len(bs.b){
		return false
	}
	bs.i++ 
	return true
}
func (bs *byteRuler) Value() byte{
	return bs.b[bs.i-1]
}

func (bs *byteRuler) Next() (size, widthPx int, err error) {
	if bs.i == len(bs.b) {
		return 0, 0, io.EOF
	}
	widthPx = bs.MeasureByte(bs.b[bs.i])
	bs.lastWidth = widthPx
	bs.lastSize = 1
	bs.i += 1
	bs.w += widthPx
	return 1, widthPx, nil
}

func (bs *byteRuler) Last() []byte {
	return bs.b[: bs.i]
}
func (bs *byteRuler) Len() int   { return bs.i }
func (bs *byteRuler) Width() int { return bs.w }
func (bs *byteRuler) Reset(p []byte) {
	bs.b = p
	bs.i = 0
	bs.w = 0
	bs.lastSize = 0
	bs.lastWidth = 0
}
func (bs *byteRuler) Unread() (size, widthPx int, err error) {
	if bs.lastSize == -1 {
		return 0, 0, ErrDoubleUnread
	}
	if bs.i-bs.lastSize < 0 {
		return 0, 0, ErrZeroRead
	}
	lsz := bs.lastSize
	bs.i -= lsz
	bs.w -= bs.lastWidth
	bs.lastSize = -1
	return lsz, bs.lastWidth, nil
}

type runeRuler struct {
	font.Face
	b                   []byte
	i                   int
	w                   int
	lastSize, lastWidth int
}

func (bs *runeRuler) Advance() (err error) {
	if bs.i == len(bs.b) {
		return io.EOF
	}
	bs.b = bs.b[bs.i:]
	bs.i = 0
	bs.w = 0
	bs.lastSize = 0
	bs.lastWidth = 0

	return nil
}

func Fix(i fixed.Int26_6) int {
	return i.Ceil()
}
func (bs *runeRuler) MeasureRune(r rune) int {
	return 8
	x, _ := bs.Face.GlyphAdvance(r)
	return Fix(x)
}

func NewRuneRuler(b []byte, ft font.Face) Ruler {
	return &runeRuler{Face: ft, b: b}
}

func (bs *runeRuler) Bytes() []byte { return bs.b[:bs.i] }


func (bs *runeRuler) Next() (size, widthPx int, err error) {
	if bs.i == len(bs.b) {
		return 0, 0, io.EOF
	}
	r, size := utf8.DecodeRune(bs.b[bs.i:])
	if r == utf8.RuneError {
	}
	if size == 0 {
		panic("size = 0")
	}
	widthPx = bs.MeasureRune(r)
	bs.i += size
	bs.w += widthPx
	bs.lastSize = size
	bs.lastWidth = widthPx
	//	log.Printf("size:%d widthpx: %d err=%s\n", size, widthPx, err)
	return
}

func (bs *runeRuler) Unread() (size, widthPx int, err error) {
	if bs.lastSize == -1 {
		return 0, 0, ErrDoubleUnread
	}
	if bs.i-bs.lastSize < 0 {
		return 0, 0, ErrZeroRead
	}
	lsz := bs.lastSize
	bs.i -= lsz
	bs.w -= bs.lastWidth
	bs.lastSize = -bs.lastSize
	return lsz, bs.lastWidth, nil
}


func (bs *runeRuler) Last() []byte {
	//	log.Printf("called Last: bytes=%s\n\tbs.i=%d\n\tlastsize=%d\n\n",bs.b, bs.i, bs.lastSize)
	return bs.b[:bs.i]
}

func (bs *runeRuler) Reset(p []byte) {
	bs.b = p
	bs.i = 0
	bs.w = 0
	bs.lastSize = 0
	bs.lastWidth = 0
}
func (bs *runeRuler) Len() int   { return bs.i }
func (bs *runeRuler) Width() int { return bs.w }
