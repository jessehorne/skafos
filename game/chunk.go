package game

type Chunk struct {
	X      int
	Y      int
	W      int
	H      int
	Blocks map[int]map[int]*Block
}

func NewChunk(x, y, w, h int) *Chunk {
	newChunk := &Chunk{
		X:      x,
		Y:      y,
		W:      w,
		H:      h,
		Blocks: map[int]map[int]*Block{},
	}

	for ty := 0; ty < h; ty++ {
		newChunk.Blocks[ty] = map[int]*Block{}
		for tx := 0; tx < w; tx++ {
			newBlock := NewBlock(BlockTypeDirt, BlockTypeDirtFrameDirt)
			newChunk.Blocks[ty][tx] = newBlock
		}
	}

	return newChunk
}
