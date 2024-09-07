package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"strconv"
)

const (
	UnderlyingTypePlaceableBlock byte = 0
)

type InventoryItem struct {
	UnderlyingType        byte
	ItemType              byte
	Frame                 byte
	Amount                int
	InventoryPosition     pixel.Vec
	DrawPosition          pixel.Vec
	ShouldUseDrawPosition bool
	Count                 *text.Text
	Font                  font.Face
	Atlas                 *text.Atlas
}

func NewInventoryItem(underType, itemType, frame byte, amt int, inventoryPos pixel.Vec) *InventoryItem {
	newItem := &InventoryItem{
		UnderlyingType:    underType,
		ItemType:          itemType,
		Frame:             frame,
		Amount:            amt,
		InventoryPosition: inventoryPos,
	}
	count := text.New(pixel.V(0, 0), Atlas)
	count.Color = colornames.White
	count.WriteString("0")
	newItem.Count = count
	return newItem
}

func (i *InventoryItem) GetCraftingPosition(win *opengl.Window, scale float64) pixel.Vec {
	craftingOffsetX := win.Bounds().W()/2 + 16*scale
	craftingOffsetY := win.Bounds().H()/2 + 7*scale

	posX := craftingOffsetX + i.InventoryPosition.X*(16*scale)
	posY := craftingOffsetY + i.InventoryPosition.Y*(16*scale)

	return pixel.V(posX, posY)
}

func (i *InventoryItem) DrawCraftingItem(win *opengl.Window, scale float64) {
	pos := pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, 3.0).Moved(i.GetCraftingPosition(win, scale))
	Tiles[i.ItemType][0].Draw(win, pos)
	i.Count.Clear()
	i.Count.WriteString(strconv.Itoa(i.Amount))
	i.Count.Draw(win, pixel.IM)
}

func (i *InventoryItem) Draw(win *opengl.Window) {
	if i.ShouldUseDrawPosition {
		drawPos := pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, 3.0).Moved(i.DrawPosition)
		Tiles[i.ItemType][0].Draw(win, drawPos)
	} else {
		pos := i.GetDrawPosition(win)

		drawPos := pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, 3.0).Moved(pos)
		Tiles[i.ItemType][0].Draw(win, drawPos)

		i.Count.Clear()
		i.Count.WriteString(strconv.Itoa(i.Amount))
		i.Count.Orig = pos
		i.Count.Draw(win, pixel.IM)
	}
}

func (i *InventoryItem) GetDrawPosition(win *opengl.Window) pixel.Vec {
	scale := 4.0
	posX := i.InventoryPosition.X * 16 * scale
	posY := i.InventoryPosition.Y * 16 * scale

	offsetX := win.Bounds().W()/2 - (8 * 16) - 16*scale
	offsetY := 4.0

	return pixel.V(offsetX+posX, offsetY+posY+36)
}

func GetInventoryItemDrawPosition(win *opengl.Window, x, y int) pixel.Vec {
	scale := 4.0
	posX := float64(x) * 16 * scale
	posY := float64(y) * 16 * scale

	offsetX := win.Bounds().W()/2 - (8 * 16) - 16*scale
	offsetY := 4.0

	return pixel.V(offsetX+posX, offsetY+posY+36)
}
