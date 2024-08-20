package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

type Map struct {
	Name         string
	Chunks       map[int]map[int]*Chunk
	Spritesheets map[string]*Spritesheet
	DrawBatch    *pixel.Batch // the holder for batch drawing
	DrawRadius   int          // how many chunks around the current center chunk should be drawn
	ChunkX       int          // the current center chunks X value
	ChunkY       int          // the current center chunks Y value
	Tiles        map[string]*pixel.Sprite
}

func NewMap(name string, s *Spritesheet) *Map {
	return &Map{
		Name:   name,
		Chunks: map[int]map[int]*Chunk{},
		Spritesheets: map[string]*Spritesheet{
			"all": s,
		},
		DrawBatch: pixel.NewBatch(&pixel.TrianglesData{}, s.Picture),
		Tiles: map[string]*pixel.Sprite{
			"dirt": pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H(), 16, s.Picture.Bounds().H()-16)),
		},
		DrawRadius: 7,
		ChunkX:     0,
		ChunkY:     0,
	}
}

func (m *Map) GenerateAllDirtChunk(x, y int, force bool) {
	newChunk := NewChunk(x, y, 16, 16)

	_, yExists := m.Chunks[y]
	if !yExists {
		m.Chunks[y] = map[int]*Chunk{}
	}

	_, xExists := m.Chunks[y][x]
	if !xExists {
		m.Chunks[y][x] = newChunk
	} else {
		if force {
			m.Chunks[y][x] = newChunk
		}
	}
}

// RefreshDrawBatch loads the chunks around the maps center chunk using
func (m *Map) RefreshDrawBatch() {
	m.DrawBatch.Clear()

	// load tiles into batch around player
	for y := m.ChunkY - m.DrawRadius; y < m.ChunkY+m.DrawRadius; y++ {
		for x := m.ChunkX - m.DrawRadius; x < m.ChunkX+m.DrawRadius; x++ {
			_, yExists := m.Chunks[y]
			if !yExists {
				continue
			}

			_, xExists := m.Chunks[y][x]
			if !xExists {
				continue
			}

			for ty := 0; ty < 16; ty++ {
				for tx := 0; tx < 16; tx++ {
					// get tile
					tile, tileExists := m.Chunks[y][x].Blocks[ty][tx]
					if !tileExists {
						continue
					}

					// where the chunks are draw from top left
					chunkOffsetX := x * 256
					chunkOffsetY := y * 256

					// where the tiles are drawn relative to the chunks top-left position
					tileX := chunkOffsetX + tx*16
					tileY := chunkOffsetY + ty*16

					// add tile to batch
					if tile.Type == BlockTypeDirt && tile.Frame == BlockTypeDirtFrameDirt {
						m.Tiles["dirt"].Draw(m.DrawBatch, pixel.IM.Moved(pixel.V(float64(tileX), float64(tileY))))
					}

				}
			}
		}
	}
}

func (m *Map) Draw(win *opengl.Window) {
	m.DrawBatch.Clear()
	m.RefreshDrawBatch()
	m.DrawBatch.Draw(win)
}
