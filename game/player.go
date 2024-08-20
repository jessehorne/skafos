package game

type Player struct {
	X float64
	Y float64
}

func NewPlayer() *Player {
	return &Player{
		X: 0,
		Y: 0,
	}
}
