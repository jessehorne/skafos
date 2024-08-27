package game

import "github.com/gopxl/pixel/v2"

type InventoryItem struct {
	ItemType byte
	Amount   int
	Position pixel.Vec
}

func NewInventoryItem(itemType byte, amt int, pos pixel.Vec) *InventoryItem {
	return &InventoryItem{
		ItemType: itemType,
		Amount:   amt,
		Position: pos,
	}
}
