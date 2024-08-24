package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

const (
	BlockTypeDirt          byte = 0
	BlockTypeDirtFrameDirt byte = 0

	BlockTypeGrass       byte = 1
	BlockTypeGrassFrame1 byte = 1
	BlockTypeGrassFrame2 byte = 2
	BlockTypeGrassFrame3 byte = 3
	BlockTypeGrassFrame4 byte = 4

	BlockTypeTree                 byte = 2
	BlockTypeTreeFrameSapling     byte = 5
	BlockTypeTreeFrameGrownTop    byte = 6
	BlockTypeTreeFrameGrownBottom byte = 7
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
