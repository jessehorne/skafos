package game

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
