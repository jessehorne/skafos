package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

type Game struct {
	Map    *Map
	Player *Player
	Camera *Camera
}

func NewGame(name string) (*Game, error) {
	p, err := NewPlayer()
	if err != nil {
		return nil, err
	}

	m, err := NewMap(name)
	if err != nil {
		return nil, err
	}

	return &Game{
		Map:    m,
		Player: p,
		Camera: NewCamera(),
	}, nil
}

func (g *Game) Update(win *opengl.Window, dt float64) {
	// generate chunks as player walks around
	g.Map.ChunkPosition = g.Player.GetChunkPosition()
	g.Map.GenerateChunksAroundPlayer()

	g.Player.Update(win, dt)
	g.Camera.Update(g.Player.Position)
}

func (g *Game) Draw(win *opengl.Window) {
	cam := pixel.IM.Scaled(g.Camera.Position, g.Camera.Zoom).Moved(win.Bounds().Center().Sub(g.Camera.Position))
	win.SetMatrix(cam)

	// draw map
	g.Map.FloorBatch.Clear()
	g.Map.TreeBatchBottom.Clear()
	g.Map.TreeBatchTop.Clear()

	g.Map.RefreshDrawBatch()
	g.Map.FloorBatch.Draw(win)
	g.Map.TreeBatchBottom.Draw(win)

	// draw player
	g.Player.Draw(win)

	g.Map.TreeBatchTop.Draw(win)
}
