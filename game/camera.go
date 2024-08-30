package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

type Camera struct {
	Matrix    pixel.Matrix
	Position  pixel.Vec
	Speed     float64
	Zoom      float64
	ZoomSpeed float64
}

func NewCamera() *Camera {
	return &Camera{
		Matrix:    pixel.IM,
		Position:  pixel.ZV,
		Speed:     100.0,
		Zoom:      4.0,
		ZoomSpeed: 1.2,
	}
}

func (c *Camera) Update(pos pixel.Vec) {
	c.Position = pos
}

func (c *Camera) StartCamera(win *opengl.Window) {
	c.Matrix = pixel.IM.Scaled(c.Position, c.Zoom).Moved(win.Bounds().Center().Sub(c.Position))
	win.SetMatrix(c.Matrix)
}

func (c *Camera) EndCamera(win *opengl.Window) {
	win.SetMatrix(pixel.IM)
}
