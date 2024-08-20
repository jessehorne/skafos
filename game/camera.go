package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

type Camera struct {
	Position  pixel.Vec
	Speed     float64
	Zoom      float64
	ZoomSpeed float64
}

func NewCamera() *Camera {
	return &Camera{
		Position:  pixel.ZV,
		Speed:     100.0,
		Zoom:      1.0,
		ZoomSpeed: 1.2,
	}
}

func (c *Camera) Update(win *opengl.Window, dt float64) {
	if win.Pressed(pixel.KeyA) {
		c.Position.X -= c.Speed * dt
	}
	if win.Pressed(pixel.KeyD) {
		c.Position.X += c.Speed * dt
	}
	if win.Pressed(pixel.KeyS) {
		c.Position.Y -= c.Speed * dt
	}
	if win.Pressed(pixel.KeyW) {
		c.Position.Y += c.Speed * dt
	}
}
