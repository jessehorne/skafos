package game

import "github.com/gopxl/pixel/v2"

type IntVec struct {
	X int
	Y int
}

func NewIntVec(x, y int) IntVec {
	return IntVec{
		X: x,
		Y: y,
	}
}

func (i IntVec) ToVec() pixel.Vec {
	return pixel.V(float64(i.X), float64(i.Y))
}
