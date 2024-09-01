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
	Velocity      pixel.Vec
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

func NewFloater(win *opengl.Window, t byte, position, velocity pixel.Vec) *Floater {
	f := &Floater{
		Position:      position,
		Velocity:      velocity,
		OldPosition:   position,
		ItemType:      t,
		Size:          pixel.V(8, 8),
		Scale:         0.5,
		ScaleSpeed:    0.25,
		ScaleMax:      0.5,
		ScaleMin:      0.4,
		ScaleDir:      0, // 0 == up & 1 == down
		Solid:         false,
		RotationSpeed: 3,
		Rotation:      0,
		DebugRect:     MakeDebugRect(win, 8, 8),
	}

	if t == BlockTypeDirt {
		f.Sprite = Tiles[BlockTypeDirt][BlockTypeDirtFrameDirt]
	} else if t == BlockTypeGrass {
		f.Sprite = Tiles[BlockTypeGrass][BlockTypeGrassFrame1]
	}

	return f
}

func (f *Floater) GetPosition() pixel.Vec {
	return f.Position
}

func (f *Floater) GetSize() pixel.Vec {
	return f.Size
}

func (f *Floater) Collide(c Collideable) {

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

	f.Position = f.Position.Add(pixel.V(f.Velocity.X*dt, f.Velocity.Y*dt))

	change := 90.0
	diff := 0.02
	if f.Velocity.X < 0 {
		f.Velocity.X += change * dt

		if f.Velocity.X > -diff {
			f.Velocity.X = 0
		}
	} else if f.Velocity.X > 0 {
		f.Velocity.X -= change * dt

		if f.Velocity.X < diff {
			f.Velocity.X = 0
		}
	}

	if f.Velocity.Y < 0 {
		f.Velocity.Y += change * dt

		if f.Velocity.Y > -diff {
			f.Velocity.Y = 0
		}
	} else if f.Velocity.Y > 0 {
		f.Velocity.Y -= change * dt

		if f.Velocity.Y < diff {
			f.Velocity.Y = 0
		}
	}

	if math.Abs(f.Velocity.X) < 0.1 && math.Abs(f.Velocity.Y) < 0.1 {
		f.Solid = true
	}
}

func (f *Floater) Draw(win *opengl.Window) {
	pos := pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, f.Scale).Rotated(pixel.ZV, f.Rotation).Moved(f.Position)
	FloaterBorderSprite.Draw(win, pos)
	f.Sprite.Draw(win, pos)
}

func (f *Floater) DrawDebug(win *opengl.Window) {
	f.DebugRect.Draw(win, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, 0.5).Moved(f.Position))
}
