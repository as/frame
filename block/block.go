package block

const SLOP = 25

type Block struct {
	Nrune int
	Ptr   []byte
}

func (b *Block) Len() int {
	return b.Nrune
}

func (b *Block) Bytes() []byte {
	n := b.Len()
	if n <= 0 {
		return nil
	}
	return b.Ptr[:n]
}
