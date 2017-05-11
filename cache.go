package frame

import "image"

func (f *Frame) Cacheinit() {
	f.Cache = make([]image.Rectangle, 0, 1024)
}
func (f *Frame) Flushcache() {
	f.Cache = f.Cache[:0]
}
