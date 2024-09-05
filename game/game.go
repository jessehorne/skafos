package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image"
)

var (
	Font  font.Face
	Atlas *text.Atlas

	Tiles        map[byte]map[byte]*pixel.Sprite
	Floaters     []*Floater
	Collideables []Collideable // list of objects to check for collision
	Cam          *Camera

	FloaterBorderImage  *image.RGBA
	FloaterBorderSprite *pixel.Sprite
)

type Game struct {
	Map                   *Map
	Player                *Player
	CollideablesDrawDebug bool
	GUI                   *GUI
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

	Tiles = map[byte]map[byte]*pixel.Sprite{
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

	m, err := NewMap(name, s)
	if err != nil {
		return nil, err
	}

	Cam = NewCamera()

	gui, err := NewGUI(win)
	if err != nil {
		return nil, err
	}

	g := &Game{
		Map:    m,
		Player: p,
		GUI:    gui,
	}

	return g, nil
}

func (g *Game) Init(win *opengl.Window) {
	AddCollideable(g.Player)

	FloaterBorderImage, FloaterBorderSprite = MakeRect(18, 18, colornames.Black)

	// add an example floater at 50, 50
	f := NewFloater(win, FloaterTypeDirt, pixel.V(50, 50), pixel.V(0, 0))
	Floaters = append(Floaters, f)
	AddCollideable(f)
}

func (g *Game) Update(win *opengl.Window, dt float64) {
	newFloaters := []*Floater{}
	// cleanup deleted floaters
	for _, f := range Floaters {
		if !f.Deleted {
			newFloaters = append(newFloaters, f)
		} else {
			for x, c := range Collideables {
				if f == c {
					Collideables = append(Collideables[:x], Collideables[x+1:]...)
				}
			}
		}
	}
	Floaters = newFloaters

	g.Map.ChunkPosition = g.Player.GetChunkPosition()
	g.Map.GenerateChunksAroundPlayer(g, win)

	for _, f := range Floaters {
		f.Update(dt)
	}

	g.Player.Update(win, dt)
	g.GUI.Update(dt)
	Cam.Update(g.Player.Position)

	g.CheckCollisions()
}

func (g *Game) Draw(win *opengl.Window) {
	Cam.StartCamera(win)

	g.Player.GetMouseMapBlockPosition(win)

	// draw map
	g.Map.FloorBatch.Clear()
	g.Map.TreeBatchBottom.Clear()
	g.Map.TreeBatchTop.Clear()

	g.Map.RefreshDrawBatch()
	g.Map.FloorBatch.Draw(win)
	g.Map.TreeBatchBottom.Draw(win)

	// draw floaters
	for _, f := range Floaters {
		f.Draw(win)
	}

	// draw player
	g.Player.Draw(win, g.GUI)

	g.Map.TreeBatchTop.Draw(win)

	// debug
	if g.CollideablesDrawDebug {
		for i := 0; i < len(Collideables); i++ {
			Collideables[i].DrawDebug(win)
		}
	}

	Cam.EndCamera(win)

	g.GUI.SetInventoryItems(g.Player.Inventory)
	g.GUI.Draw()
}

func (g *Game) ButtonCallback(btn pixel.Button, action pixel.Action) {
	g.Player.ButtonCallback(btn, action)
	g.GUI.ButtonCallback(btn, action)
}

func (g *Game) CharCallback(r rune) {
	if r == ']' {
		g.CollideablesDrawDebug = !g.CollideablesDrawDebug
	} else if r == 'i' {
		g.GUI.ShouldDrawInventory = !g.GUI.ShouldDrawInventory

		if g.GUI.HoldingInvItem != nil {
			g.GUI.HoldingInvItem.ShouldUseDrawPosition = false
			g.GUI.HoldingInvItem.Count.Orig = g.GUI.HoldingInvItem.GetDrawPosition(g.GUI.Window)
			i := g.GUI.HoldingInvItem

			g.GUI.Inventory[int(i.InventoryPosition.Y)][int(i.InventoryPosition.X)] = g.GUI.HoldingInvItem

			g.GUI.HoldingInvItem = nil
		}

		g.Player.InInventory = g.GUI.ShouldDrawInventory
	}

	g.Player.CharCallback(r)
}

func AddCollideable(c Collideable) {
	Collideables = append(Collideables, c)
}

func (g *Game) CheckCollisions() {
	for i := 0; i < len(Collideables); i++ {
		for x := 0; x < len(Collideables); x++ {
			if x == i {
				continue
			}

			first := Collideables[i]
			second := Collideables[x]

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
