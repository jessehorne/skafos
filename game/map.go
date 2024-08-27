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
	Tiles           map[byte]map[byte]*pixel.Sprite
}

func NewMap(name string, s *Spritesheet, tiles map[byte]map[byte]*pixel.Sprite) (*Map, error) {
	return &Map{
		Name:   name,
		Chunks: map[int]map[int]*Chunk{},
		Spritesheets: map[string]*Spritesheet{
			"all": s,
		},
		FloorBatch:      pixel.NewBatch(&pixel.TrianglesData{}, s.Picture),
		TreeBatchBottom: pixel.NewBatch(&pixel.TrianglesData{}, s.Picture),
		TreeBatchTop:    pixel.NewBatch(&pixel.TrianglesData{}, s.Picture),
		Tiles:           tiles,
		DrawRadius:      4,
		ChunkPosition:   pixel.V(0, 0),
	}, nil
}

func (m *Map) GenerateChunksAroundPlayer(g *Game, win *opengl.Window) {
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
			m.GenerateAllDirtChunk(win, int(x), int(y), true, g)
		}
	}
}

func (m *Map) GenerateAllDirtChunk(win *opengl.Window, x, y int, force bool, g *Game) {
	newChunk := NewChunk(win, x, y, 16, 16, "grass", g)

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

					// add tile to batch
					for i := 0; i < len(tiles); i++ {
						tile := tiles[i]

						if tile.Type == BlockTypeTree {
							treeTops = append(treeTops, tile.GetPosition())
							m.Tiles[BlockTypeTree][BlockTypeTreeFrameGrownBottom].Draw(m.TreeBatchBottom, pixel.IM.Moved(tile.GetPosition()))
						} else {
							m.Tiles[tile.Type][tile.Frame].Draw(m.FloorBatch, pixel.IM.Moved(tile.GetPosition()))
						}
					}
				}
			}
		}
	}

	// add tree tops to their batch
	for i := len(treeTops) - 1; i >= 0; i-- {
		m.Tiles[BlockTypeTree][BlockTypeTreeFrameGrownTop].Draw(m.TreeBatchTop, pixel.IM.Moved(treeTops[i]))
	}
}

func (m *Map) Draw(win *opengl.Window) {

}
