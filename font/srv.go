package font

import (
	"hash/crc32"
	"image"

	"github.com/golang/freetype/truetype"
)

var fontIRQ chan fontPKT

type fontPKT struct {
	id    string
	data  []byte
	reply chan interface{}
}

func fontsrv() {
	parsedTTF := make(map[string]*truetype.Font)
	for {
		select {
		case in := <-fontIRQ:
			func() {
				defer close(in.reply)

				if v, ok := parsedTTF[in.id]; ok {
					in.reply <- v
					return
				}

				f, err := truetype.Parse(in.data)
				if err != nil {
					in.reply <- err
					return
				}

				parsedTTF[in.id] = f
				in.reply <- f
			}()
		}
	}
}
func makefont(data []byte, size int) *Font {
	if fontIRQ == nil {
		fontIRQ = make(chan fontPKT)
		go fontsrv()
	}
	reply := make(chan interface{})
	fontIRQ <- fontPKT{
		id:    string(crc32.NewIEEE().Sum(data)),
		reply: reply,
		data:  data,
	}
	rx := <-reply
	switch rx := rx.(type) {
	case error:
		println(rx)
		return nil
	case *truetype.Font:
		return &Font{
			Face: truetype.NewFace(rx,
				&truetype.Options{
					Size: float64(size),
				}),
			size:     size,
			ascent:   2,
			descent:  +(size / 3),
			stride:   0,
			data:     data,
			imgCache: make(map[signature]*image.RGBA),
		}
	}
	panic("makefont")
}
