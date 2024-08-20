package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"math"
)

const (
	PlayerDirectionUp    byte = 0
	PlayerDirectionDown  byte = 1
	PlayerDirectionLeft  byte = 2
	PlayerDirectionRight byte = 3

	PlayerWalking byte = 0
	PlayerRunning byte = 1
)

type Player struct {
	Position          pixel.Vec
	Speed             map[byte]float64 // pixels per second
	WalkingOrRunning  byte
	Spritesheet       *Spritesheet
	Frames            map[byte][]*pixel.Sprite
	FrameSpeed        map[byte]float64
	CurrentFrame      float64
	MaxMovementFrame  float64
	MovementDirection byte
}

func NewPlayer() (*Player, error) {
	s, err := NewSpritesheet("./assets/player/character.png")
	if err != nil {
		return nil, err
	}

	//w := s.Picture.Bounds().W()
	h := s.Picture.Bounds().H()

	return &Player{
		Position: pixel.V(0, 0),
		Speed: map[byte]float64{
			PlayerWalking: 32,
			PlayerRunning: 64,
		},
		WalkingOrRunning: PlayerWalking,
		Spritesheet:      s,
		Frames: map[byte][]*pixel.Sprite{
			PlayerDirectionLeft: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h-(3*32), 16, h-(4*32))),
				pixel.NewSprite(s.Picture, pixel.R(16, h-(3*32), 2*16, h-(4*32))),
				pixel.NewSprite(s.Picture, pixel.R(2*16, h-(3*32), 3*16, h-(4*32))),
				pixel.NewSprite(s.Picture, pixel.R(3*16, h-(3*32), 4*16, h-(4*32))),
			},
			PlayerDirectionRight: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h-(32), 16, h-(2*32))),
				pixel.NewSprite(s.Picture, pixel.R(16, h-(32), 2*16, h-(2*32))),
				pixel.NewSprite(s.Picture, pixel.R(2*16, h-(32), 3*16, h-(2*32))),
				pixel.NewSprite(s.Picture, pixel.R(3*16, h-(32), 4*16, h-(2*32))),
			},
			PlayerDirectionUp: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h-(2*32), 16, h-(3*32))),
				pixel.NewSprite(s.Picture, pixel.R(16, h-(2*32), 2*16, h-(3*32))),
				pixel.NewSprite(s.Picture, pixel.R(2*16, h-(2*32), 3*16, h-(3*32))),
				pixel.NewSprite(s.Picture, pixel.R(3*16, h-(2*32), 4*16, h-(3*32))),
			},
			PlayerDirectionDown: []*pixel.Sprite{
				pixel.NewSprite(s.Picture, pixel.R(0, h, 16, h-(32))),
				pixel.NewSprite(s.Picture, pixel.R(16, h, 2*16, h-(32))),
				pixel.NewSprite(s.Picture, pixel.R(2*16, h, 3*16, h-(32))),
				pixel.NewSprite(s.Picture, pixel.R(3*16, h, 4*16, h-(32))),
			},
		},
		FrameSpeed: map[byte]float64{
			PlayerWalking: 4,
			PlayerRunning: 8,
		}, // change frame this many times per second
		CurrentFrame:     0,
		MaxMovementFrame: 4,
	}, nil
}

func (p *Player) Update(win *opengl.Window, dt float64) {
	if win.Pressed(pixel.KeyA) {
		p.Position.X -= p.Speed[p.WalkingOrRunning] * dt
		p.CurrentFrame += p.FrameSpeed[p.WalkingOrRunning] * dt
		p.MovementDirection = PlayerDirectionLeft
	}
	if win.Pressed(pixel.KeyD) {
		p.Position.X += p.Speed[p.WalkingOrRunning] * dt
		p.CurrentFrame += p.FrameSpeed[p.WalkingOrRunning] * dt
		p.MovementDirection = PlayerDirectionRight
	}
	if win.Pressed(pixel.KeyS) {
		p.Position.Y -= p.Speed[p.WalkingOrRunning] * dt
		p.CurrentFrame += p.FrameSpeed[p.WalkingOrRunning] * dt
		p.MovementDirection = PlayerDirectionDown
	}
	if win.Pressed(pixel.KeyW) {
		p.Position.Y += p.Speed[p.WalkingOrRunning] * dt
		p.CurrentFrame += p.FrameSpeed[p.WalkingOrRunning] * dt
		p.MovementDirection = PlayerDirectionUp
	}
	if win.Pressed(pixel.KeyLeftControl) {
		p.WalkingOrRunning = PlayerRunning
	} else {
		p.WalkingOrRunning = PlayerWalking
	}

	if p.CurrentFrame > p.MaxMovementFrame {
		p.CurrentFrame = 0
	}
}

func (p *Player) Draw(win *opengl.Window) {
	currentFrame := int(math.Floor(p.CurrentFrame))

	p.Frames[p.MovementDirection][currentFrame].Draw(win, pixel.IM.Moved(p.Position))
}
