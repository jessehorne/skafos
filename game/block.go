package game

const (
	BlockTypeDirt          byte = 0
	BlockTypeDirtFrameDirt byte = 0
)

type Block struct {
	Type  byte
	Frame byte
}

func NewBlock(blockType, frame byte) *Block {
	return &Block{
		Type:  blockType,
		Frame: frame,
	}
}
