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
	s, err := NewSpritesheet("./assets/tiles/all.jpg")
	if err != nil {
		return nil, err
	}

	p, err := NewPlayer()
	if err != nil {
		return nil, err
	}

	return &Game{
		Map:    NewMap(name, s),
		Player: p,
		Camera: NewCamera(),
	}, nil
}

func (g *Game) Update(win *opengl.Window, dt float64) {
	g.Player.Update(win, dt)
	g.Camera.Update(g.Player.Position)
}

func (g *Game) Draw(win *opengl.Window) {
	cam := pixel.IM.Scaled(g.Camera.Position, g.Camera.Zoom).Moved(win.Bounds().Center().Sub(g.Camera.Position))
	win.SetMatrix(cam)

	// draw map
	g.Map.Draw(win)

	// draw player
	g.Player.Draw(win)
}
