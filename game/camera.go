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
		Zoom:      4.0,
		ZoomSpeed: 1.2,
	}
}

func (c *Camera) Update(pos pixel.Vec) {
	c.Position = pos
}

func (c *Camera) StartCamera(win *opengl.Window) {
	camera := pixel.IM.Scaled(c.Position, c.Zoom).Moved(win.Bounds().Center().Sub(c.Position))
	win.SetMatrix(camera)
}

func (c *Camera) EndCamera(win *opengl.Window) {
	win.SetMatrix(pixel.IM)
}
