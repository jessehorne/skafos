package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

const (
	CollideableTypeBlock   byte = 0
	CollideableTypePlayer  byte = 1
	CollideableTypeFloater byte = 2
)

type Collideable interface {
	GetPosition() pixel.Vec
	GetSize() pixel.Vec
	Collide(Collideable)
	IsSolid() bool
	GetType() byte
	DrawDebug(*opengl.Window)
	GetOldPosition() pixel.Vec
}
