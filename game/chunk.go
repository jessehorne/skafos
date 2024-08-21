package game

import "math/rand/v2"

type Chunk struct {
	X      int
	Y      int
	W      int
	H      int
	Blocks map[int]map[int][]*Block
}

func NewChunk(x, y, w, h int, chunkType string) *Chunk {
	newChunk := &Chunk{
		X:      x,
		Y:      y,
		W:      w,
		H:      h,
		Blocks: map[int]map[int][]*Block{},
	}

	for ty := 0; ty < h; ty++ {
		newChunk.Blocks[ty] = map[int][]*Block{}
		for tx := 0; tx < w; tx++ {
			var newBlock *Block

			if chunkType == "dirt" {
				newBlock = NewBlock(BlockTypeDirt, BlockTypeDirtFrameDirt)
			} else if chunkType == "grass" {
				var frame byte
				rnd := rand.IntN(100)

				if rnd < 80 {
					frame = BlockTypeGrassFrame1
				} else if rnd < 95 {
					frame = BlockTypeGrassFrame2
				} else if rnd < 99 {
					frame = BlockTypeGrassFrame3
				} else if rnd <= 100 {
					frame = BlockTypeGrassFrame4
				}
				newBlock = NewBlock(BlockTypeGrass, frame)
			}

			newChunk.Blocks[ty][tx] = append(newChunk.Blocks[ty][tx], newBlock)

			// add trees maybe
			treeRnd := rand.IntN(100)

			if treeRnd <= 2 {
				newTreeBlock := NewBlock(BlockTypeTree, BlockTypeTreeFrameGrownTop)
				newChunk.Blocks[ty][tx] = append(newChunk.Blocks[ty][tx], newTreeBlock)
			}
		}
	}

	return newChunk
}
