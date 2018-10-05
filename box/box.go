package box

const SLOP = 25

type Box struct {
	Ptr      []byte
	Nrune    int
	Width    int
	Minwidth int
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
	if b.Nrune > 0{
		return b.Ptr[:b.Nrune]
	}
	if b.Nrune < 0{
		return b.Ptr[:1]
	}
	return nil
}
