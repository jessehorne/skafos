package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"math"
)

const (
	FloaterTypeDirt byte = 0
)

// When an item is dropped, it "floats" and rotates around on the ground until someone picks it up
type Floater struct {
	Position      pixel.Vec
	OldPosition   pixel.Vec
	ItemType      byte
	Size          pixel.Vec
	Scale         float64
	Solid         bool
	Sprite        *pixel.Sprite
	RotationSpeed float64
	Rotation      float64
	ScaleSpeed    float64
	ScaleMax      float64
	ScaleMin      float64
	ScaleDir      float64
	DebugRect     *pixel.Sprite
	Deleted       bool
}

func NewFloater(win *opengl.Window, t byte, position pixel.Vec, sprite *pixel.Sprite) *Floater {
	return &Floater{
		Position:      position,
		OldPosition:   position,
		ItemType:      t,
		Size:          pixel.V(8, 8),
		Scale:         0.5,
		ScaleSpeed:    0.25,
		ScaleMax:      0.5,
		ScaleMin:      0.4,
		ScaleDir:      0, // 0 == up & 1 == down
		Solid:         true,
		Sprite:        sprite,
		RotationSpeed: 3,
		Rotation:      0,
		DebugRect:     MakeDebugRect(win, 8, 8),
	}
}

func (f *Floater) GetPosition() pixel.Vec {
	return f.Position
}

func (f *Floater) GetSize() pixel.Vec {
	return f.Size
}

func (f *Floater) Collide(c Collideable) {
	if c.GetType() == CollideableTypePlayer {
		f.Deleted = true
	}
}

func (f *Floater) IsSolid() bool {
	return f.Solid
}

func (f *Floater) GetType() byte {
	return CollideableTypeFloater
}

func (f *Floater) GetOldPosition() pixel.Vec {
	return f.OldPosition
}

func (f *Floater) Update(dt float64) {
	f.Rotation -= dt * f.RotationSpeed
	if f.Rotation <= -math.Pi*(2) {
		f.Rotation = 0
	}

	if f.ScaleDir == 0 {
		f.Scale += f.ScaleSpeed * dt
	} else if f.ScaleDir == 1 {
		f.Scale -= f.ScaleSpeed * dt
	}

	if f.Scale >= f.ScaleMax {
		f.ScaleDir = 1
	} else if f.Scale <= f.ScaleMin {
		f.ScaleDir = 0
	}
}

func (f *Floater) Draw(win *opengl.Window) {
	f.Sprite.Draw(win, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, f.Scale).Rotated(pixel.ZV, f.Rotation).Moved(f.Position))
}

func (f *Floater) DrawDebug(win *opengl.Window) {
	f.DebugRect.Draw(win, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, 0.5).Moved(f.Position))
}
