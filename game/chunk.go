package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"math/rand/v2"
)

type Chunk struct {
	X      int
	Y      int
	W      int
	H      int
	Blocks map[int]map[int][]*Block
}

func NewChunk(win *opengl.Window, x, y, w, h int, chunkType string, g *Game) *Chunk {
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

			pos := pixel.V(float64(x)*256+float64(tx)*16, float64(y)*256+float64(ty)*16)

			if chunkType == "dirt" {
				newBlock = NewBlock(win, BlockTypeDirt, BlockTypeDirtFrameDirt, pos)
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
				newBlock = NewBlock(win, BlockTypeGrass, frame, pos)
			}

			newChunk.Blocks[ty][tx] = append(newChunk.Blocks[ty][tx], newBlock)

			// add trees maybe
			objRnd := rand.IntN(1000)

			if objRnd <= 20 {
				spawnSafe := 50.0
				if (pos.X > spawnSafe || pos.X < -spawnSafe) && (pos.Y > spawnSafe || pos.Y < -spawnSafe) {
					newTreeBlock := NewBlock(win, BlockTypeTree, BlockTypeTreeFrameGrownTop, pos)
					newChunk.Blocks[ty][tx] = append(newChunk.Blocks[ty][tx], newTreeBlock)
					g.AddCollideable(newTreeBlock)
				}
			} else if objRnd > 21 && objRnd < 24 {
				newStoneBlock := NewBlock(win, BlockTypeStone, BlockTypeStoneFrame1, pos)
				newChunk.Blocks[ty][tx] = append(newChunk.Blocks[ty][tx], newStoneBlock)
			} else if objRnd > 24 && objRnd < 28 {
				newCopperBlock := NewBlock(win, BlockTypeCopper, BlockTypeCopperFrame1, pos)
				newChunk.Blocks[ty][tx] = append(newChunk.Blocks[ty][tx], newCopperBlock)
			}
		}
	}

	return newChunk
}
