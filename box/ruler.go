package box

import (
	"errors"
	"github.com/as/frame/font"
	"io"
	"unicode/utf8"
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
}

type byteRuler struct {
	*font.Font
	b                   []byte
	i                   int
	w                   int
	lastSize, lastWidth int
}

func (bs *byteRuler) Bytes() []byte {
	return bs.b[:bs.i]
}

func NewByteRuler(b []byte, ft *font.Font) *byteRuler {
	return &byteRuler{Font: ft, b: b}
}

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

func (bs *byteRuler) Last() []byte {
	return bs.b[bs.i-bs.lastSize : bs.i]
}

func (bs *byteRuler) Len() int {
	return bs.i
}
func (bs *byteRuler) Width() int {
	return bs.w
}

type runeRuler struct {
	*font.Font
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

func NewRuneRuler(b []byte, ft *font.Font) runeRuler {
	return runeRuler{Font: ft, b: b}
}

func (bs runeRuler) Next() (size, widthPx int, err error) {
	if bs.i == len(bs.b) {
		return 0, 0, io.EOF
	}
	r, size := utf8.DecodeRune(bs.b[bs.i:])
	widthPx = bs.MeasureRune(r)
	bs.i += size
	bs.w += widthPx
	bs.lastSize = size
	bs.lastWidth = widthPx
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
	bs.lastSize = -1
	return lsz, bs.lastWidth, nil
}

func (bs *runeRuler) Last() []byte {
	return bs.b[:bs.i]
}

func (bs *runeRuler) Len() int {
	return bs.i
}
func (bs *runeRuler) Width() int {
	return bs.w
}
