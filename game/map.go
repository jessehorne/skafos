package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

type Map struct {
	Name            string
	Chunks          map[int]map[int]*Chunk
	Spritesheets    map[string]*Spritesheet
	FloorBatch      *pixel.Batch // the holder for batch drawing
	TreeBatchBottom *pixel.Batch
	TreeBatchTop    *pixel.Batch
	DrawRadius      float64   // how many chunks around the current center chunk should be drawn
	ChunkPosition   pixel.Vec // the current center chunk
	Tiles           map[byte]*pixel.Sprite
}

func NewMap(name string) (*Map, error) {
	s, err := NewSpritesheet("./assets/tiles/all.png")
	if err != nil {
		return nil, err
	}

	return &Map{
		Name:   name,
		Chunks: map[int]map[int]*Chunk{},
		Spritesheets: map[string]*Spritesheet{
			"all": s,
		},
		FloorBatch:      pixel.NewBatch(&pixel.TrianglesData{}, s.Picture),
		TreeBatchBottom: pixel.NewBatch(&pixel.TrianglesData{}, s.Picture),
		TreeBatchTop:    pixel.NewBatch(&pixel.TrianglesData{}, s.Picture),
		Tiles: map[byte]*pixel.Sprite{
			BlockTypeDirtFrameDirt:        pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H(), 16, s.Picture.Bounds().H()-16)),
			BlockTypeGrassFrame1:          pixel.NewSprite(s.Picture, pixel.R(16, s.Picture.Bounds().H(), 16*2, s.Picture.Bounds().H()-16)),
			BlockTypeGrassFrame2:          pixel.NewSprite(s.Picture, pixel.R(2*16, s.Picture.Bounds().H(), 3*16, s.Picture.Bounds().H()-16)),
			BlockTypeGrassFrame3:          pixel.NewSprite(s.Picture, pixel.R(3*16, s.Picture.Bounds().H(), 4*16, s.Picture.Bounds().H()-16)),
			BlockTypeGrassFrame4:          pixel.NewSprite(s.Picture, pixel.R(4*16, s.Picture.Bounds().H(), 5*16, s.Picture.Bounds().H()-16)),
			BlockTypeTreeFrameSapling:     pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H()-4*16, 16, s.Picture.Bounds().H()-5*16)),
			BlockTypeTreeFrameGrownTop:    pixel.NewSprite(s.Picture, pixel.R(16, s.Picture.Bounds().H()-4*16, 3*16, s.Picture.Bounds().H()-6*16)),
			BlockTypeTreeFrameGrownBottom: pixel.NewSprite(s.Picture, pixel.R(3*16, s.Picture.Bounds().H()-4*16, 5*16, s.Picture.Bounds().H()-6*16)),
		},
		DrawRadius:    4,
		ChunkPosition: pixel.V(0, 0),
	}, nil
}

func (m *Map) GenerateChunksAroundPlayer() {
	for y := m.ChunkPosition.Y - m.DrawRadius; y < m.ChunkPosition.Y+m.DrawRadius; y++ {
		for x := m.ChunkPosition.X - m.DrawRadius; x < m.ChunkPosition.X+m.DrawRadius; x++ {
			// check if chunk exists
			_, yExists := m.Chunks[int(y)]

			if yExists {
				_, xExists := m.Chunks[int(y)][int(x)]

				if xExists {
					continue
				}
			}

			// generate chunk
			m.GenerateAllDirtChunk(int(x), int(y), true)
		}
	}
}

func (m *Map) GenerateAllDirtChunk(x, y int, force bool) {
	newChunk := NewChunk(x, y, 16, 16, "grass")

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
	m.FloorBatch.Clear()

	treeTops := []pixel.Vec{} // so we can redraw in reverse later because drawing from top to bottom causes overlapping issue

	// load tiles into batch around player
	for y := m.ChunkPosition.Y - m.DrawRadius; y < m.ChunkPosition.Y+m.DrawRadius; y++ {
		for x := m.ChunkPosition.X - m.DrawRadius; x < m.ChunkPosition.X+m.DrawRadius; x++ {
			_, yExists := m.Chunks[int(y)]
			if !yExists {
				continue
			}

			_, xExists := m.Chunks[int(y)][int(x)]
			if !xExists {
				continue
			}

			for ty := 0; ty < 16; ty++ {
				for tx := 0; tx < 16; tx++ {
					// get tile
					tiles, tileExists := m.Chunks[int(y)][int(x)].Blocks[ty][tx]
					if !tileExists {
						continue
					}

					// where the chunks are draw from top left
					chunkOffsetX := x * 256
					chunkOffsetY := y * 256

					// where the tiles are drawn relative to the chunks top-left position
					tileX := chunkOffsetX + float64(tx*16)
					tileY := chunkOffsetY + float64(ty*16)

					tilePosition := pixel.V(float64(tileX), float64(tileY))

					// add tile to batch
					for i := 0; i < len(tiles); i++ {
						tile := tiles[i]
						if tile.Type == BlockTypeTree {
							treeTops = append(treeTops, tilePosition)
							m.Tiles[BlockTypeTreeFrameGrownBottom].Draw(m.TreeBatchBottom, pixel.IM.Moved(tilePosition))
						} else {
							m.Tiles[tile.Frame].Draw(m.FloorBatch, pixel.IM.Moved(tilePosition))
						}
					}
				}
			}
		}
	}

	// add tree tops to their batch
	for i := len(treeTops) - 1; i >= 0; i-- {
		m.Tiles[BlockTypeTreeFrameGrownTop].Draw(m.TreeBatchTop, pixel.IM.Moved(treeTops[i]))
	}
}

func (m *Map) Draw(win *opengl.Window) {

}
