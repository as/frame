package box

import (
	"errors"
	"github.com/as/frame/font"
	"io"
	"unicode/utf8"
	//	"log"
)

var (
	ErrZeroRead     = errors.New("zero length read")
	ErrDoubleUnread = errors.New("double unread")
)

type Ruler interface {
	Advance() error
	Bytes() string
	Next() (size, width int, err error)
	Last() string
	Len() int
	Unread() (size, width int, err error)
	Width() int
	Reset(p string)
}

type byteRuler struct {
	*font.Font
	b                   string
	i                   int
	w                   int
	lastSize, lastWidth int
}

type runeRuler struct {
	*font.Font
	b                   string
	i                   int
	w                   int
	lastSize, lastWidth int
}

func NewByteRuler(b string, ft *font.Font) Ruler {
	return &byteRuler{Font: ft, b: b}
}

func NewRuneRuler(b string, ft *font.Font) Ruler {
	return &runeRuler{Font: ft, b: b}
}

func (bs *byteRuler) Bytes() string { return bs.b[:bs.i] }
func (bs *runeRuler) Bytes() string { return bs.b[:bs.i] }

func (bs *byteRuler) Advance() (err error) {
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

func (bs *byteRuler) Next() (size, widthPx int, err error) {
	if bs.i == len(bs.b) {
		return 0, 0, io.EOF
	}
	size = 1
	widthPx = bs.Font.MeasureByte(bs.b[bs.i])
	bs.i += size
	bs.w += widthPx
	bs.lastSize = size
	bs.lastWidth = widthPx
	return
}

func (bs *runeRuler) Next() (size, widthPx int, err error) {
	if bs.i == len(bs.b) {
		return 0, 0, io.EOF
	}
	r, size := utf8.DecodeRuneInString(bs.b[bs.i:])
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

func (bs *byteRuler) Last() string {
	return bs.b[bs.i-bs.lastSize : bs.i]
}

func (bs *runeRuler) Last() string {
	//	log.Printf("called Last: bytes=%s\n\tbs.i=%d\n\tlastsize=%d\n\n",bs.b, bs.i, bs.lastSize)
	return bs.b[bs.i-bs.lastSize : bs.i]
}

func (bs *byteRuler) Reset(p string) {
	bs.b = p
	bs.i = 0
	bs.w = 0
	bs.lastSize = 0
	bs.lastWidth = 0
}
func (bs *runeRuler) Reset(p string) {
	bs.b = p
	bs.i = 0
	bs.w = 0
	bs.lastSize = 0
	bs.lastWidth = 0
}
func (bs *byteRuler) Len() int   { return bs.i }
func (bs *byteRuler) Width() int { return bs.w }
func (bs *runeRuler) Len() int   { return bs.i }
func (bs *runeRuler) Width() int { return bs.w }
