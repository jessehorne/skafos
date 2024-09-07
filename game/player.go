package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"image"
	"math"
)

const (
	PlayerDirectionUp    byte = 0
	PlayerDirectionDown  byte = 1
	PlayerDirectionLeft  byte = 2
	PlayerDirectionRight byte = 3

	PlayerWalking byte = 0
	PlayerRunning byte = 1
)

type Player struct {
	Position            pixel.Vec
	OldPosition         pixel.Vec
	Speed               map[byte]float64 // pixels per second
	WalkingOrRunning    byte
	Spritesheet         *Spritesheet
	Frames              map[byte][]*pixel.Sprite
	SwingFrames         map[byte][]*pixel.Sprite
	FrameSpeed          map[byte]float64
	Inventory           [][]*InventoryItem
	InventoryW          int
	InventoryH          int
	HotbarX             int
	ShouldDrawInventory bool
	CurrentFrame        float64
	MaxMovementFrame    float64
	MovementDirection   byte
	MovementDirections  []byte
	Solid               bool
	DebugRect           *pixel.Sprite
	IsSwinging          bool
	SwingFrameSpeed     float64
	InInventory         bool
	MouseX              int // block position from bottom left
	MouseY              int // block position from bottom left
	MouseRectImage      *image.RGBA
	MouseRectSprite     *pixel.Sprite
}

func NewPlayer(win *opengl.Window) (*Player, error) {
	s, err := NewSpritesheet("./assets/player/character.png")
	if err != nil {
		return nil, err
	}

	//w := s.Picture.Bounds().W()
	h := s.Picture.Bounds().H()

	mSprite := MakeDebugRect(win, 16, 16)

	p := &Player{
		Position: pixel.V(0, 0),
		Speed: map[byte]float64{
			PlayerWalking: 32,
			PlayerRunning: 64,
		},
		WalkingOrRunning: PlayerWalking,
		Spritesheet:      s,
		Frames: map[byte][]*pixel.Sprite{
			PlayerDirectionLeft: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h-(3*32), 16, h-(4*32))),
				pixel.NewSprite(s.Picture, pixel.R(16, h-(3*32), 2*16, h-(4*32))),
				pixel.NewSprite(s.Picture, pixel.R(2*16, h-(3*32), 3*16, h-(4*32))),
				pixel.NewSprite(s.Picture, pixel.R(3*16, h-(3*32), 4*16, h-(4*32))),
			},
			PlayerDirectionRight: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h-(32), 16, h-(2*32))),
				pixel.NewSprite(s.Picture, pixel.R(16, h-(32), 2*16, h-(2*32))),
				pixel.NewSprite(s.Picture, pixel.R(2*16, h-(32), 3*16, h-(2*32))),
				pixel.NewSprite(s.Picture, pixel.R(3*16, h-(32), 4*16, h-(2*32))),
			},
			PlayerDirectionUp: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h-(2*32), 16, h-(3*32))),
				pixel.NewSprite(s.Picture, pixel.R(16, h-(2*32), 2*16, h-(3*32))),
				pixel.NewSprite(s.Picture, pixel.R(2*16, h-(2*32), 3*16, h-(3*32))),
				pixel.NewSprite(s.Picture, pixel.R(3*16, h-(2*32), 4*16, h-(3*32))),
			},
			PlayerDirectionDown: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h, 16, h-(32))),
				pixel.NewSprite(s.Picture, pixel.R(16, h, 2*16, h-(32))),
				pixel.NewSprite(s.Picture, pixel.R(2*16, h, 3*16, h-(32))),
				pixel.NewSprite(s.Picture, pixel.R(3*16, h, 4*16, h-(32))),
			},
		},
		SwingFrames: map[byte][]*pixel.Sprite{
			PlayerDirectionDown: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h-(4*32), 32, h-(5*32))),
				pixel.NewSprite(s.Picture, pixel.R(32, h-(4*32), 2*32, h-(5*32))),
				pixel.NewSprite(s.Picture, pixel.R(2*32, h-(4*32), 3*32, h-(5*32))),
				pixel.NewSprite(s.Picture, pixel.R(3*32, h-(4*32), 4*32, h-(5*32))),
			},
			PlayerDirectionUp: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h-(5*32), 32, h-(6*32))),
				pixel.NewSprite(s.Picture, pixel.R(32, h-(5*32), 2*32, h-(6*32))),
				pixel.NewSprite(s.Picture, pixel.R(2*32, h-(5*32), 3*32, h-(6*32))),
				pixel.NewSprite(s.Picture, pixel.R(3*32, h-(5*32), 4*32, h-(6*32))),
			},
			PlayerDirectionRight: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h-(6*32), 32, h-(7*32))),
				pixel.NewSprite(s.Picture, pixel.R(32, h-(6*32), 2*32, h-(7*32))),
				pixel.NewSprite(s.Picture, pixel.R(2*32, h-(6*32), 3*32, h-(7*32))),
				pixel.NewSprite(s.Picture, pixel.R(3*32, h-(6*32), 4*32, h-(7*32))),
			},
			PlayerDirectionLeft: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h-(7*32), 32, h-(8*32))),
				pixel.NewSprite(s.Picture, pixel.R(32, h-(7*32), 2*32, h-(8*32))),
				pixel.NewSprite(s.Picture, pixel.R(2*32, h-(7*32), 3*32, h-(8*32))),
				pixel.NewSprite(s.Picture, pixel.R(3*32, h-(7*32), 4*32, h-(8*32))),
			},
		},
		SwingFrameSpeed: 20,
		FrameSpeed: map[byte]float64{
			PlayerWalking: 4,
			PlayerRunning: 8,
		}, // change frame this many times per second
		CurrentFrame:        0,
		MaxMovementFrame:    4,
		Solid:               true,
		DebugRect:           MakeDebugRect(win, 16, 16),
		MovementDirections:  []byte{},
		Inventory:           [][]*InventoryItem{},
		InventoryW:          7,
		InventoryH:          3, // not counting hotbar
		ShouldDrawInventory: false,
		MouseRectSprite:     mSprite,
	}

	p.ClearInventory()
	p.AddInventoryItem(NewInventoryItem(UnderlyingTypePlaceableBlock, BlockTypeDirt, BlockTypeDirtFrameDirt, 10, pixel.V(0, 0)))

	return p, nil
}

func (p *Player) AddMovementDirection(d byte) {
	// only add if its not in
	for _, oldDir := range p.MovementDirections {
		if oldDir == d {
			return
		}
	}

	p.MovementDirections = append(p.MovementDirections, d)
}

func (p *Player) RemoveMovementDirection(d byte) {
	if len(p.MovementDirections) == 0 {
		return
	}

	// only add if its not in
	for i, oldDir := range p.MovementDirections {
		if oldDir == d {
			p.MovementDirections = append(p.MovementDirections[:i], p.MovementDirections[i+1:]...)
			return
		}
	}
}

func (p *Player) IsMovingInDirection(d byte) bool {
	for _, oldDir := range p.MovementDirections {
		if oldDir == d {
			return true
		}
	}

	return false
}

func (p *Player) Update(win *opengl.Window, dt float64) {
	if win.Pressed(pixel.KeyA) {
		p.AddMovementDirection(PlayerDirectionLeft)
	} else {
		p.RemoveMovementDirection(PlayerDirectionLeft)
	}

	if win.Pressed(pixel.KeyD) {
		p.AddMovementDirection(PlayerDirectionRight)
	} else {
		p.RemoveMovementDirection(PlayerDirectionRight)
	}

	if win.Pressed(pixel.KeyW) {
		p.AddMovementDirection(PlayerDirectionUp)
	} else {
		p.RemoveMovementDirection(PlayerDirectionUp)
	}

	if win.Pressed(pixel.KeyS) {
		p.AddMovementDirection(PlayerDirectionDown)
	} else {
		p.RemoveMovementDirection(PlayerDirectionDown)
	}

	p.OldPosition = p.Position

	if p.IsMovingInDirection(PlayerDirectionUp) {
		p.MovementDirection = PlayerDirectionUp

		p.Position.Y += p.Speed[p.WalkingOrRunning] * dt
		if p.IsMovingInDirection(PlayerDirectionLeft) {
			p.Position.X -= p.Speed[p.WalkingOrRunning] / 2 * dt
		} else if p.IsMovingInDirection(PlayerDirectionRight) {
			p.Position.X += p.Speed[p.WalkingOrRunning] / 2 * dt
		}
	} else if p.IsMovingInDirection(PlayerDirectionDown) {
		p.MovementDirection = PlayerDirectionDown

		p.Position.Y -= p.Speed[p.WalkingOrRunning] * dt
		if p.IsMovingInDirection(PlayerDirectionLeft) {
			p.Position.X -= p.Speed[p.WalkingOrRunning] / 2 * dt
		} else if p.IsMovingInDirection(PlayerDirectionRight) {
			p.Position.X += p.Speed[p.WalkingOrRunning] / 2 * dt
		}
	} else if p.IsMovingInDirection(PlayerDirectionLeft) {
		p.MovementDirection = PlayerDirectionLeft
		p.Position.X -= p.Speed[p.WalkingOrRunning] * dt
	} else if p.IsMovingInDirection(PlayerDirectionRight) {
		p.MovementDirection = PlayerDirectionRight
		p.Position.X += p.Speed[p.WalkingOrRunning] * dt
	}

	if len(p.MovementDirections) > 0 && !p.IsSwinging {
		p.CurrentFrame += p.FrameSpeed[p.WalkingOrRunning] * dt
	}

	if p.IsSwinging {
		p.CurrentFrame += p.SwingFrameSpeed * dt
	}

	if p.CurrentFrame > p.MaxMovementFrame {
		if p.IsSwinging {
			p.IsSwinging = !p.IsSwinging
		}
		p.CurrentFrame = 0
	}
}

func (p *Player) Draw(game *Game) {
	currentFrame := int(math.Floor(p.CurrentFrame))

	if !game.GUI.ShouldDrawInventory {
		p.MouseRectSprite.Draw(game.Window, pixel.IM.Moved(p.GetMouseMapBlockPosition(game)))
	}

	if p.IsSwinging {
		p.SwingFrames[p.MovementDirection][currentFrame].Draw(game.Window, pixel.IM.Moved(p.Position))
	} else {
		p.Frames[p.MovementDirection][currentFrame].Draw(game.Window, pixel.IM.Moved(p.Position))
	}

	game.GUI.SetHotbarItems(p.Inventory[0], p.HotbarX)
}

func (p *Player) GetChunkPosition() pixel.Vec {
	x := math.Floor(p.Position.X / 256)
	y := math.Floor(p.Position.Y / 256)

	return pixel.V(x, y)
}

func (p *Player) GetPosition() pixel.Vec {
	return p.Position
}

func (p *Player) GetSize() pixel.Vec {
	return pixel.V(16, 16)
}

func (p *Player) Collide(c Collideable) {
	if c.GetType() == CollideableTypeBlock {
		d := GetCollisionDirection(p, c)

		pos2 := c.GetPosition()
		size2 := c.GetSize()

		if d == CollisionDirectionUp {
			p.Position.Y = pos2.Y - 16
		} else if d == CollisionDirectionDown {
			p.Position.Y = pos2.Y + size2.Y
		} else if d == CollisionDirectionLeft {
			p.Position.X = pos2.X + size2.X
		} else if d == CollisionDirectionRight {
			p.Position.X = pos2.X - 16
		}
	} else if c.GetType() == CollideableTypeFloater {
		f := c.(*Floater)

		if !f.Deleted {
			p.AddItemToInventory(f.UnderlyingType, f.ItemType, f.Frame)
			f.Deleted = true
		}
	}
}

func (p *Player) IsSolid() bool {
	return p.Solid
}

func (p *Player) GetType() byte {
	return CollideableTypePlayer
}

func (p *Player) GetOldPosition() pixel.Vec {
	return p.OldPosition
}

func (p *Player) DrawDebug(win *opengl.Window) {
	p.DebugRect.Draw(win, pixel.IM.Moved(p.Position))
}

func (p *Player) ButtonCallback(game *Game, btn pixel.Button, action pixel.Action) {
	if btn == pixel.MouseButtonLeft && action == pixel.Press {
		if !p.InInventory {
			if !p.IsSwinging {
				p.CurrentFrame = 0
				p.IsSwinging = true
			}
		}
	} else if btn == pixel.MouseButtonRight && action == pixel.Press {
		p.HandleRightClick(game)
	} else if btn == pixel.KeyLeftControl && action == pixel.Press {
		if p.WalkingOrRunning == PlayerWalking {
			p.WalkingOrRunning = PlayerRunning
		} else if p.WalkingOrRunning == PlayerRunning {
			p.WalkingOrRunning = PlayerWalking
		}
	}
}

func (p *Player) HandleRightClick(game *Game) {
	if p.InInventory {
		return
	}

	held := p.GetHeldItem()

	if held == nil {
		return
	}

	if held.UnderlyingType == UnderlyingTypePlaceableBlock {
		p.PlaceBlock(game, held)
	}
}

func (p *Player) GetMouseMapCoords(game *Game) (IntVec, IntVec) {
	mousePos := p.GetMouseMapBlockPosition(game)

	chunkX := math.Floor(mousePos.X / 256)
	chunkY := math.Floor(mousePos.Y / 256)
	x := (mousePos.X - chunkX*256) / 16
	y := (mousePos.Y - chunkY*256) / 16

	return NewIntVec(int(chunkX), int(chunkY)), NewIntVec(int(x), int(y))
}

func (p *Player) PlaceBlock(game *Game, item *InventoryItem) {
	if item == nil {
		return
	}

	if item.Amount <= 0 {
		return
	}

	chunk, coords := p.GetMouseMapCoords(game)

	exists := game.Map.BlockExists(chunk, coords)

	if !exists {
		return
	}

	b := NewBlock(game.Window, item.ItemType, item.Frame, game.Map.Chunks[chunk.Y][chunk.X].Blocks[coords.Y][coords.X][0].Position)
	game.Map.Chunks[chunk.Y][chunk.X].Blocks[coords.Y][coords.X] = append(game.Map.Chunks[chunk.Y][chunk.X].Blocks[coords.Y][coords.X], b)
	item.Amount -= 1

	if item.Amount <= 0 {
		p.RemoveInventoryItem(item)
	}
}

func (p *Player) ClearInventory() {
	p.Inventory = [][]*InventoryItem{}

	for y := 0; y <= p.InventoryH; y++ {
		for x := 0; x <= p.InventoryW; x++ {
			if y == 0 {
				p.Inventory = append(p.Inventory, []*InventoryItem{})
			}

			p.Inventory[y] = append(p.Inventory[y], nil)
		}
	}
}

func (p *Player) AddInventoryItem(i *InventoryItem) {
	p.Inventory[int(i.InventoryPosition.Y)][int(i.InventoryPosition.X)] = i
}

func (p *Player) RemoveInventoryItem(i *InventoryItem) {
	p.Inventory[int(i.InventoryPosition.Y)][int(i.InventoryPosition.X)] = nil
}

func (p *Player) AddItemToInventory(underType, itemType, frame byte) {
	var foundItem *InventoryItem
	var found bool
	var foundX int
	var foundY int
	for y := 0; y < len(p.Inventory); y++ {
		for x := 0; x < len(p.Inventory[y]); x++ {
			item := p.Inventory[y][x]
			if item != nil {
				if item.ItemType == itemType {
					foundItem = item
					break
				}
			} else {
				if !found {
					foundX = x
					foundY = y
					found = true
				}
			}
		}

		if foundItem != nil {
			break
		}
	}

	if foundItem != nil {
		foundItem.Amount++
	} else {
		if found {
			newInvItem := NewInventoryItem(underType, itemType, frame, 1, pixel.V(float64(foundX), float64(foundY)))
			p.Inventory[foundY][foundX] = newInvItem
		}
	}
}

func (p *Player) CharCallback(game *Game, r rune) {
	if r >= 49 && r < 59 {
		p.HotbarX = int(r - 49)
	} else {
		if !p.InInventory {
			if r == 'q' {
				// throw currently held item
				p.ThrowInventoryItem(game)
			}
		}
	}
}

func (p *Player) GetHeldItem() *InventoryItem {
	return p.Inventory[0][p.HotbarX]
}

func (p *Player) ThrowInventoryItem(game *Game) {
	item := p.GetHeldItem()

	if item != nil {
		if item.Amount > 0 {
			item.Amount -= 1
			mousePos := game.Window.MousePosition()
			adjusted := game.Camera.Matrix.Project(p.Position)

			delta := mousePos.Sub(adjusted)
			l := delta.Len()

			delta.X = (delta.X / l) * 100.0
			delta.Y = (delta.Y / l) * 100.0

			newFloater := NewFloater(game.Window, item.UnderlyingType, item.ItemType, item.Frame, p.Position, delta)
			Floaters = append(Floaters, newFloater)
			AddCollideable(newFloater)

			if item.Amount == 0 {
				p.Inventory[0][p.HotbarX] = nil
			}
		}
	}
}

func (p *Player) GetMouseMapPosition(game *Game) pixel.Vec {
	return game.Camera.Matrix.Unproject(game.Window.MousePosition()).Add(pixel.V(8, 8))
}

// Returns the X and Y coordinates for the block the mouse is over on the map
func (p *Player) GetMouseMapBlockPosition(game *Game) pixel.Vec {
	mousePos := p.GetMouseMapPosition(game)

	x := math.Floor(mousePos.X/16) * 16
	y := math.Floor(mousePos.Y/16) * 16

	return pixel.V(x, y)
}
