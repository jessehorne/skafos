package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"strconv"
)

type InventoryItem struct {
	ItemType          byte
	Amount            int
	InventoryPosition pixel.Vec
	Count             *text.Text
	Font              font.Face
	Atlas             *text.Atlas
}

func NewInventoryItem(win *opengl.Window, itemType byte, amt int, inventoryPos pixel.Vec) *InventoryItem {
	newItem := &InventoryItem{
		ItemType:          itemType,
		Amount:            amt,
		InventoryPosition: inventoryPos,
	}
	count := text.New(newItem.GetDrawPosition(win), Atlas)
	count.Color = colornames.White
	count.WriteString("0")
	newItem.Count = count
	return newItem
}

func (i *InventoryItem) Draw(win *opengl.Window) {
	pos := i.GetDrawPosition(win)

	drawPos := pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, 3.0).Moved(pos)
	Tiles[i.ItemType][0].Draw(win, drawPos)

	i.Count.Clear()
	i.Count.WriteString(strconv.Itoa(i.Amount))
	i.Count.Draw(win, pixel.IM)
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
