package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

const (
	BlockTypeDirt   byte = 0
	BlockTypeGrass  byte = 1
	BlockTypeTree   byte = 2
	BlockTypeStone  byte = 3
	BlockTypeCopper byte = 4

	BlockTypeDirtFrameDirt byte = 0

	BlockTypeGrassFrame1 byte = 0
	BlockTypeGrassFrame2 byte = 1
	BlockTypeGrassFrame3 byte = 2
	BlockTypeGrassFrame4 byte = 3

	BlockTypeTreeFrameSapling     byte = 0
	BlockTypeTreeFrameGrownTop    byte = 1
	BlockTypeTreeFrameGrownBottom byte = 2

	BlockTypeStoneFrame1 byte = 0

	BlockTypeCopperFrame1 byte = 0
)

type Block struct {
	Position  pixel.Vec
	Type      byte
	Frame     byte
	DebugRect *pixel.Sprite
}

func NewBlock(win *opengl.Window, blockType, frame byte, pos pixel.Vec) *Block {
	return &Block{
		Type:      blockType,
		Frame:     frame,
		Position:  pos,
		DebugRect: MakeDebugRect(win, 16, 16),
	}
}

func (b *Block) GetPosition() pixel.Vec {
	return b.Position
}

func (b *Block) GetSize() pixel.Vec {
	return pixel.V(16, 16)
}

func (b *Block) Collide(c Collideable) {

}

func (b *Block) IsSolid() bool {
	if b.Type == BlockTypeTree {
		return true
	}

	return false
}

func (b *Block) GetType() byte {
	return CollideableTypeBlock
}

func (b *Block) GetOldPosition() pixel.Vec {
	return b.Position
}

func (b *Block) DrawDebug(win *opengl.Window) {
	b.DebugRect.Draw(win, pixel.IM.Moved(b.Position))
}
