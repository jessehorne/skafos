package game

const (
	BlockTypeDirt          byte = 0
	BlockTypeDirtFrameDirt byte = 0

	BlockTypeGrass       byte = 1
	BlockTypeGrassFrame1 byte = 1
	BlockTypeGrassFrame2 byte = 2
	BlockTypeGrassFrame3 byte = 3
	BlockTypeGrassFrame4 byte = 4
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
