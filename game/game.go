package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/text"
	"golang.org/x/image/font"
)

var (
	Font  font.Face
	Atlas *text.Atlas
)

type Game struct {
	Map                   *Map
	Player                *Player
	Camera                *Camera
	Collideables          []Collideable // list of objects to check for collision
	CollideablesDrawDebug bool
	GUI                   *GUI
	Floaters              []*Floater
}

func NewGame(name string, win *opengl.Window) (*Game, error) {
	face, err := loadTTF("./assets/font/munro.ttf", 24)
	if err != nil {
		return nil, err
	}

	atlas := text.NewAtlas(face, text.ASCII)

	Font = face
	Atlas = atlas

	s, err := NewSpritesheet("./assets/tiles/all.png")
	if err != nil {
		return nil, err
	}

	tiles := map[byte]map[byte]*pixel.Sprite{
		BlockTypeDirt: {
			BlockTypeDirtFrameDirt: pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H(), 16, s.Picture.Bounds().H()-16)),
		},
		BlockTypeGrass: {
			BlockTypeGrassFrame1: pixel.NewSprite(s.Picture, pixel.R(16, s.Picture.Bounds().H(), 16*2, s.Picture.Bounds().H()-16)),
			BlockTypeGrassFrame2: pixel.NewSprite(s.Picture, pixel.R(2*16, s.Picture.Bounds().H(), 3*16, s.Picture.Bounds().H()-16)),
			BlockTypeGrassFrame3: pixel.NewSprite(s.Picture, pixel.R(3*16, s.Picture.Bounds().H(), 4*16, s.Picture.Bounds().H()-16)),
			BlockTypeGrassFrame4: pixel.NewSprite(s.Picture, pixel.R(4*16, s.Picture.Bounds().H(), 5*16, s.Picture.Bounds().H()-16)),
		},
		BlockTypeTree: {
			BlockTypeTreeFrameSapling:     pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H()-4*16, 16, s.Picture.Bounds().H()-5*16)),
			BlockTypeTreeFrameGrownTop:    pixel.NewSprite(s.Picture, pixel.R(16, s.Picture.Bounds().H()-4*16, 3*16, s.Picture.Bounds().H()-6*16)),
			BlockTypeTreeFrameGrownBottom: pixel.NewSprite(s.Picture, pixel.R(3*16, s.Picture.Bounds().H()-4*16, 5*16, s.Picture.Bounds().H()-6*16)),
		},
		BlockTypeStone: {
			BlockTypeStoneFrame1: pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H()-2*16, 16, s.Picture.Bounds().H()-3*16)),
		},
		BlockTypeCopper: {
			BlockTypeCopperFrame1: pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H()-3*16, 16, s.Picture.Bounds().H()-4*16)),
		},
	}

	p, err := NewPlayer(win)
	if err != nil {
		return nil, err
	}

	m, err := NewMap(name, s, tiles)
	if err != nil {
		return nil, err
	}

	camera := NewCamera()

	gui, err := NewGUI(win, camera, tiles)
	if err != nil {
		return nil, err
	}

	g := &Game{
		Map:          m,
		Player:       p,
		Camera:       camera,
		Collideables: []Collideable{},
		GUI:          gui,
		Floaters:     []*Floater{},
	}

	return g, nil
}

func (g *Game) Init(win *opengl.Window) {
	g.AddCollideable(g.Player)

	// add an example floater at 50, 50
	dirt := g.Map.Tiles[BlockTypeDirt][BlockTypeDirtFrameDirt]
	f := NewFloater(win, FloaterTypeDirt, pixel.V(50, 50), dirt)
	g.Floaters = append(g.Floaters, f)
	g.AddCollideable(f)
}

func (g *Game) Update(win *opengl.Window, dt float64) {
	// cleanup deleted floaters
	for i, f := range g.Floaters {
		if f.Deleted {
			g.Floaters = append(g.Floaters[:i], g.Floaters[i+1:]...)

			for x, c := range g.Collideables {
				if f == c {
					g.Collideables = append(g.Collideables[:x], g.Collideables[x+1:]...)
				}
			}
		}
	}

	g.Map.ChunkPosition = g.Player.GetChunkPosition()
	g.Map.GenerateChunksAroundPlayer(g, win)

	for _, f := range g.Floaters {
		f.Update(dt)
	}

	g.Player.Update(win, dt)
	g.Camera.Update(g.Player.Position)

	g.CheckCollisions()
}

func (g *Game) Draw(win *opengl.Window) {
	g.Camera.StartCamera(win)

	// draw map
	g.Map.FloorBatch.Clear()
	g.Map.TreeBatchBottom.Clear()
	g.Map.TreeBatchTop.Clear()

	g.Map.RefreshDrawBatch()
	g.Map.FloorBatch.Draw(win)
	g.Map.TreeBatchBottom.Draw(win)

	// draw floaters
	for _, f := range g.Floaters {
		f.Draw(win)
	}

	// draw player
	g.Player.Draw(win, g.GUI)

	g.Map.TreeBatchTop.Draw(win)

	// debug
	if g.CollideablesDrawDebug {
		for i := 0; i < len(g.Collideables); i++ {
			g.Collideables[i].DrawDebug(win)
		}
	}

	g.Camera.EndCamera(win)

	g.GUI.SetInventoryItems(g.Player.Inventory)
	g.GUI.Draw()
}

func (g *Game) ButtonCallback(btn pixel.Button, action pixel.Action) {
	g.Player.ButtonCallback(btn, action)
}

func (g *Game) CharCallback(r rune) {
	if r == ']' {
		g.CollideablesDrawDebug = !g.CollideablesDrawDebug
	} else if r == 'i' {
		g.GUI.ShouldDrawInventory = !g.GUI.ShouldDrawInventory
	}

	g.Player.CharCallback(r)
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
