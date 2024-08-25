package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

type Game struct {
	Map                   *Map
	Player                *Player
	Camera                *Camera
	Collideables          []Collideable // list of objects to check for collision
	CollideablesDrawDebug bool
	GUI                   *GUI
}

func NewGame(name string, win *opengl.Window) (*Game, error) {
	p, err := NewPlayer(win)
	if err != nil {
		return nil, err
	}

	m, err := NewMap(name)
	if err != nil {
		return nil, err
	}

	gui, err := NewGUI(win)
	if err != nil {
		return nil, err
	}

	return &Game{
		Map:          m,
		Player:       p,
		Camera:       NewCamera(),
		Collideables: []Collideable{},
		GUI:          gui,
	}, nil
}

func (g *Game) Init() {
	g.AddCollideable(g.Player)
}

func (g *Game) Update(win *opengl.Window, dt float64) {
	g.Map.ChunkPosition = g.Player.GetChunkPosition()
	g.Map.GenerateChunksAroundPlayer(g, win)

	g.Player.Update(win, dt)
	g.Camera.Update(g.Player.Position)

	g.CheckCollisions()
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

	// debug
	if g.CollideablesDrawDebug {
		for i := 0; i < len(g.Collideables); i++ {
			g.Collideables[i].DrawDebug(win)
		}
	}

	win.SetMatrix(pixel.IM)

	g.GUI.Draw()
}

func (g *Game) ButtonCallback(btn pixel.Button, action pixel.Action) {

}

func (g *Game) CharCallback(r rune) {
	if r == ']' {
		g.CollideablesDrawDebug = !g.CollideablesDrawDebug
	}
}

func (g *Game) AddCollideable(c Collideable) {
	g.Collideables = append(g.Collideables, c)
}

func (g *Game) CheckCollisions() {
	for i := 0; i < len(g.Collideables); i++ {
		for x := 0; x < len(g.Collideables); x++ {
			if x == i {
				continue
			}

			first := g.Collideables[i]
			second := g.Collideables[x]

			if !first.IsSolid() || !second.IsSolid() {
				continue
			}

			if CollisionBBox(first.GetPosition(), first.GetSize(), second.GetPosition(), second.GetSize()) {
				first.Collide(second)
				second.Collide(first)
			}
		}
	}
}
