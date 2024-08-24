package game

import (
	"github.com/gopxl/pixel/v2"
	"math"
)

const (
	CollisionDirectionUp    byte = 1
	CollisionDirectionDown  byte = 2
	CollisionDirectionLeft  byte = 3
	CollisionDirectionRight byte = 4
)

func CollisionBBox(pos1, size1, pos2, size2 pixel.Vec) bool {
	return pos1.X < pos2.X+size2.X &&
		pos1.X+size2.X > pos2.X &&
		pos1.Y < pos2.Y+size2.Y &&
		pos1.Y+size1.Y > pos2.Y
}

func GetCollisionDirection(c1, c2 Collideable) byte {
	pos1 := c1.GetPosition()
	pos2 := c2.GetPosition()
	
	disX := math.Abs(pos1.X - pos2.X)
	disY := math.Abs(pos1.Y - pos2.Y)

	// right
	if pos1.X < pos2.X && (disX > disY) {
		return CollisionDirectionRight
	}

	// left
	if pos1.X > pos2.X && (disX > disY) {
		return CollisionDirectionLeft
	}

	// up
	if pos1.Y < pos2.Y && (disX < disY) {
		return CollisionDirectionUp
	}

	// down
	if pos1.Y > pos2.Y && (disX < disY) {
		return CollisionDirectionDown
	}

	return 0
}
