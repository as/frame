package box

import "golang.org/x/image/math/fixed"

type i26 = fixed.Int26_6

type Box struct {
	Ptr      []byte
	Nrune    int
	Width    i26
	Minwidth i26
}

func (b Box) Break() byte {
	if b.Nrune == 0 {
		return 0
	}
	return b.Ptr[0]
}

func (b Box) Len() int {
	if b.Nrune < 0 {
		return 1
	}
	return b.Nrune
}

func (b Box) Bytes() []byte {
	if b.Nrune > 0 {
		return b.Ptr[:b.Nrune]
	}
	if b.Nrune < 0 {
		return b.Ptr[:1]
	}
	return nil
}
