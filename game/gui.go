package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"
	"image"
	"math"
	"strconv"
)

type GUI struct {
	Window      *opengl.Window
	Camera      *Camera
	Spritesheet *Spritesheet
	BarSprite   *pixel.Sprite

	OffsetX float64
	OffsetY float64
	Scale   float64

	Health            float64
	HealthBarPosition pixel.Vec
	HealthBarImage    *image.RGBA
	HealthBarSprite   *pixel.Sprite

	Hunger            float64
	HungerBarPosition pixel.Vec
	HungerBarImage    *image.RGBA
	HungerBarSprite   *pixel.Sprite

	Thirst            float64
	ThirstBarPosition pixel.Vec
	ThirstBarImage    *image.RGBA
	ThirstBarSprite   *pixel.Sprite

	BigSprite *pixel.Sprite
	BigOffset pixel.Vec

	ItemSprite *pixel.Sprite

	NeedsRedraw bool

	HotbarItems           []*InventoryItem
	HotbarX               int
	HotbarSelectionSprite *pixel.Sprite
	Inventory             [][]*InventoryItem
	ShouldDrawInventory   bool

	Tiles map[byte]map[byte]*pixel.Sprite

	HoldingInvItem *InventoryItem
}

func NewGUI(win *opengl.Window) (*GUI, error) {
	s, err := NewSpritesheet("./assets/gui.png")
	if err != nil {
		return nil, err
	}

	barSprite := pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H(), 4*16, s.Picture.Bounds().H()-16))

	healthBarImage, healthBarSprite := MakeRect(46, 4, colornames.Red)
	hungerBarImage, hungerBarSprite := MakeRect(46, 4, colornames.Red)
	thirstBarImage, thirstBarSprite := MakeRect(46, 4, colornames.Red)

	itemSprite := pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H()-16, 16, s.Picture.Bounds().H()-2*16))

	bigSprite := pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H()-2*16, s.Picture.Bounds().W(), s.Picture.Bounds().H()-7*16))

	g := &GUI{
		Window:      win,
		Spritesheet: s,
		BarSprite:   barSprite,
		OffsetX:     8 * 16,
		OffsetY:     2 * 16,
		Scale:       4.0,

		Health:          100,
		HealthBarImage:  healthBarImage,
		HealthBarSprite: healthBarSprite,

		Hunger:          100,
		HungerBarImage:  hungerBarImage,
		HungerBarSprite: hungerBarSprite,

		Thirst:          100,
		ThirstBarImage:  thirstBarImage,
		ThirstBarSprite: thirstBarSprite,

		ItemSprite: itemSprite,

		Inventory: [][]*InventoryItem{},

		HotbarSelectionSprite: pixel.NewSprite(s.Picture, pixel.R(16, s.Picture.Bounds().H()-1*16, 2*16, s.Picture.Bounds().H()-2*16)),

		BigSprite: bigSprite,
		BigOffset: pixel.V(win.Bounds().W()/2+32, 360),
	}

	//g.HealthBarPosition = pixel.V(g.OffsetX, g.Window.Bounds().H()-g.OffsetY)
	g.HealthBarPosition = pixel.V(16, g.Window.Bounds().H())
	g.HungerBarPosition = pixel.V(16, g.Window.Bounds().H()-2*16)
	g.ThirstBarPosition = pixel.V(16, g.Window.Bounds().H()-4*16)

	return g, nil
}

func (g *GUI) Draw() {
	Cam.EndCamera(g.Window)
	g.RedrawBars()
	g.DrawHotbar()

	if g.ShouldDrawInventory {
		g.DrawInventory()
	}

	if g.HoldingInvItem != nil {
		g.HoldingInvItem.Draw(g.Window)
		g.HoldingInvItem.Count.Draw(g.Window, pixel.IM)
	}

	Cam.StartCamera(g.Window)
}

func (g *GUI) RedrawBars() {
	// health
	g.UpdateHealth(50)
	g.BarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.HealthBarPosition.Add(pixel.V(g.OffsetX, -g.OffsetY))))

	// hunger
	g.UpdateHunger(50)
	g.BarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.HungerBarPosition.Add(pixel.V(g.OffsetX, -g.OffsetY))))

	// thirst
	g.UpdateThirst(50)
	g.BarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.ThirstBarPosition.Add(pixel.V(g.OffsetX, -g.OffsetY))))
}

func (g *GUI) Update(dt float64) {
	if g.HoldingInvItem != nil {
		g.HoldingInvItem.DrawPosition = g.Window.MousePosition()
		g.HoldingInvItem.Count.Orig = g.Window.MousePosition()
		g.HoldingInvItem.Count.Dot = g.Window.MousePosition()
	}
}

func (g *GUI) UpdateHealth(v float64) {
	amt := 0.48 * v

	for y := 0; y < 4; y++ {
		for x := 0; x < 46; x++ {
			if x <= int(amt) {
				g.HealthBarImage.Set(x, y, colornames.Red)
			} else {
				g.HealthBarImage.Set(x, y, pixel.RGBA{R: 0, G: 0, B: 0, A: 1})
			}
		}
	}

	g.HealthBarSprite = pixel.NewSprite(pixel.PictureDataFromImage(g.HealthBarImage), pixel.R(0, 0, 46, 4))

	g.HealthBarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.HealthBarPosition.Add(pixel.V(g.OffsetX-32, -g.OffsetY))))
}

func (g *GUI) UpdateHunger(v float64) {
	amt := 0.48 * v

	for y := 0; y < 4; y++ {
		for x := 0; x < 46; x++ {
			if x <= int(amt) {
				g.HungerBarImage.Set(x, y, colornames.Green)
			} else {
				g.HungerBarImage.Set(x, y, pixel.RGBA{R: 0, G: 0, B: 0, A: 1})
			}
		}
	}

	g.HungerBarSprite = pixel.NewSprite(pixel.PictureDataFromImage(g.HungerBarImage), pixel.R(0, 0, 46, 4))

	g.HungerBarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.HungerBarPosition.Add(pixel.V(g.OffsetX-32, -g.OffsetY))))
}

func (g *GUI) UpdateThirst(v float64) {
	amt := 0.48 * v

	for y := 0; y < 4; y++ {
		for x := 0; x < 46; x++ {
			if x <= int(amt) {
				g.ThirstBarImage.Set(x, y, colornames.Blue)
			} else {
				g.ThirstBarImage.Set(x, y, pixel.RGBA{R: 0, G: 0, B: 0, A: 1})
			}
		}
	}

	g.ThirstBarSprite = pixel.NewSprite(pixel.PictureDataFromImage(g.ThirstBarImage), pixel.R(0, 0, 46, 4))

	g.ThirstBarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.ThirstBarPosition.Add(pixel.V(g.OffsetX-32, -g.OffsetY))))
}

func (g *GUI) DrawHotbar() {
	for x, i := range g.HotbarItems {
		offsetX := g.Window.Bounds().W()/2 - (8 * 16) - 16*g.Scale
		// draw behind box
		g.ItemSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(pixel.V(offsetX+float64(x)*16*g.Scale, 40)))

		// draw items in players inventory if exists

		if i != nil {
			i.Draw(g.Window)
		}
	}

	// draw hotbar selection
	drawPos := GetInventoryItemDrawPosition(g.Window, g.HotbarX, 0)
	g.HotbarSelectionSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(drawPos))
}

func (g *GUI) SetHotbarItems(items []*InventoryItem, hotbarX int) {
	g.HotbarItems = items
	g.HotbarX = hotbarX
}

func (g *GUI) DrawInventory() {
	items := g.Inventory

	for y := 0; y < len(items); y++ {
		if y == 0 {
			continue
		}

		for x := 0; x < len(items[y]); x++ {
			posY := float64(y*16) * g.Scale

			offsetX := g.Window.Bounds().W()/2 - (8 * 16) - 16*g.Scale

			// draw behind box
			g.ItemSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(pixel.V(offsetX+float64(x)*16*g.Scale, 40+posY)))

			// draw items in players inventory if exists
			invItem := items[y][x]

			if invItem != nil {
				invItem.Draw(g.Window)
			}
		}
	}

	// draw big sprite
	g.BigSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.BigOffset))
}

func (g *GUI) SetInventoryItems(items [][]*InventoryItem) {
	g.Inventory = items
}

func (g *GUI) ButtonCallback(btn pixel.Button, action pixel.Action) {
	if btn == pixel.MouseButtonLeft && action == pixel.Press {
		if g.ShouldDrawInventory {
			g.HandleInventoryLeftClick()
		}
	} else if btn == pixel.MouseButtonRight && action == pixel.Press {
		if g.ShouldDrawInventory {
			g.HandleInventoryRightClick()
		}
	}
}

func (g *GUI) HandleInventoryLeftClick() {
	// check if mouse pressed inventory item
	mousePos := g.Window.MousePosition()

	offsetX := g.Window.Bounds().W()/2 - (8 * 16) - 8*g.Scale
	offsetY := 0.0

	clickedX := int(math.Floor((mousePos.X-offsetX)/(16*g.Scale)) + 1)
	clickedY := int(math.Floor((mousePos.Y - offsetY) / (16 * g.Scale)))

	invItem := g.Inventory[clickedY][clickedX]

	// if invItem is nil it means we're clicking into an inventory spot with nothing in it
	if invItem == nil {
		// if we're holding an item, we should place it there and stop holding it
		if g.HoldingInvItem != nil {
			g.HoldingInvItem.InventoryPosition = pixel.V(float64(clickedX), float64(clickedY))
			g.HoldingInvItem.ShouldUseDrawPosition = false
			g.HoldingInvItem.Count.Orig = g.HoldingInvItem.GetDrawPosition(g.Window)
			g.Inventory[clickedY][clickedX] = g.HoldingInvItem
			g.HoldingInvItem = nil
		}
	} else {
		if g.HoldingInvItem != nil {
			// if invItem isn't nil, it means we're trying to either merge stacks or toggle between holding what is under the mouse cursor
			if invItem.ItemType == g.HoldingInvItem.ItemType {
				// merge stacks
				invItem.Amount += g.HoldingInvItem.Amount
				g.HoldingInvItem = nil
			} else {
				toDrop := g.HoldingInvItem
				toDrop.InventoryPosition = invItem.InventoryPosition
				toDrop.Count.Orig = g.HoldingInvItem.GetDrawPosition(g.Window)

				toPickup := invItem
				toPickup.ShouldUseDrawPosition = true
				g.HoldingInvItem = toPickup

				toDrop.ShouldUseDrawPosition = false
				g.Inventory[clickedY][clickedX] = toDrop
			}
		} else {
			g.HoldingInvItem = invItem
			g.HoldingInvItem.ShouldUseDrawPosition = true
			g.Inventory[clickedY][clickedX] = nil
		}
	}
}

func (g *GUI) HandleInventoryRightClick() {
	// check if mouse pressed inventory item
	mousePos := g.Window.MousePosition()

	offsetX := g.Window.Bounds().W()/2 - (8 * 16) - 8*g.Scale
	offsetY := 0.0

	clickedX := int(math.Floor((mousePos.X-offsetX)/(16*g.Scale)) + 1)
	clickedY := int(math.Floor((mousePos.Y - offsetY) / (16 * g.Scale)))

	invItem := g.Inventory[clickedY][clickedX]

	// if invItem is nil it means we're clicking into an inventory spot with nothing in it
	if invItem != nil {
		if g.HoldingInvItem == nil {
			if invItem.Amount >= 2 {
				half := invItem.Amount / 2
				invItem.Amount -= half

				newItem := NewInventoryItem(g.Window, invItem.ItemType, half, invItem.InventoryPosition)
				newItem.ShouldUseDrawPosition = true
				newItem.Count.Clear()
				newItem.Count.WriteString(strconv.Itoa(half))
				g.HoldingInvItem = newItem
			}
		}
	}
}
