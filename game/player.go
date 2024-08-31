package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
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
	Window              *opengl.Window
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
}

func NewPlayer(win *opengl.Window) (*Player, error) {
	s, err := NewSpritesheet("./assets/player/character.png")
	if err != nil {
		return nil, err
	}

	//w := s.Picture.Bounds().W()
	h := s.Picture.Bounds().H()

	p := &Player{
		Window:   win,
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
		InventoryW:          8,
		InventoryH:          5,
		ShouldDrawInventory: false,
	}

	p.ClearInventory()
	p.AddInventoryItem(NewInventoryItem(win, BlockTypeDirt, 10, pixel.V(0, 0)))
	p.AddInventoryItem(NewInventoryItem(win, BlockTypeDirt, 10, pixel.V(1, 0)))
	p.AddInventoryItem(NewInventoryItem(win, BlockTypeDirt, 10, pixel.V(3, 0)))
	p.AddInventoryItem(NewInventoryItem(win, BlockTypeDirt, 10, pixel.V(4, 0)))
	p.AddInventoryItem(NewInventoryItem(win, BlockTypeDirt, 10, pixel.V(5, 3)))

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

	if win.Pressed(pixel.KeyLeftControl) {
		p.WalkingOrRunning = PlayerRunning
	} else {
		p.WalkingOrRunning = PlayerWalking
	}

	if p.CurrentFrame > p.MaxMovementFrame {
		if p.IsSwinging {
			p.IsSwinging = !p.IsSwinging
		}
		p.CurrentFrame = 0
	}
}

func (p *Player) Draw(win *opengl.Window, gui *GUI) {
	currentFrame := int(math.Floor(p.CurrentFrame))

	if p.IsSwinging {
		p.SwingFrames[p.MovementDirection][currentFrame].Draw(win, pixel.IM.Moved(p.Position))
	} else {
		p.Frames[p.MovementDirection][currentFrame].Draw(win, pixel.IM.Moved(p.Position))
	}

	gui.SetHotbarItems(p.Inventory[0], p.HotbarX)
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
			p.AddItemToInventory(f.ItemType)
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

func (p *Player) ButtonCallback(btn pixel.Button, action pixel.Action) {
	if btn == pixel.MouseButtonLeft && action == pixel.Press {
		if !p.InInventory {
			if !p.IsSwinging {
				p.CurrentFrame = 0
				p.IsSwinging = true
			}
		}
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

func (p *Player) AddItemToInventory(itemType byte) {
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
			newInvItem := NewInventoryItem(p.Window, itemType, 1, pixel.V(float64(foundX), float64(foundY)))
			p.Inventory[foundY][foundX] = newInvItem
		}
	}
}

func (p *Player) CharCallback(r rune) {
	if r >= 49 && r < 59 {
		p.HotbarX = int(r - 49)
	} else {
		if !p.InInventory {
			if r == 'q' {
				// throw currently held item
				p.ThrowInventoryItem()
			}
		}
	}
}

func (p *Player) GetHeldItem() *InventoryItem {
	return p.Inventory[0][p.HotbarX]
}

func (p *Player) ThrowInventoryItem() {
	item := p.GetHeldItem()

	if item != nil {
		if item.Amount > 0 {
			item.Amount -= 1
			mousePos := p.Window.MousePosition()
			adjusted := Cam.Matrix.Project(p.Position)

			delta := mousePos.Sub(adjusted)
			l := delta.Len()

			delta.X = (delta.X / l) * 100.0
			delta.Y = (delta.Y / l) * 100.0

			newFloater := NewFloater(p.Window, item.ItemType, p.Position, delta)
			Floaters = append(Floaters, newFloater)
			AddCollideable(newFloater)

			if item.Amount == 0 {
				p.Inventory[0][p.HotbarX] = nil
			}
		}
	}
}
