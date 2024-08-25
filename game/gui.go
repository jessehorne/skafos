package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"
	"image"
)

type GUI struct {
	Window      *opengl.Window
	Spritesheet *Spritesheet
	BarSprite   *pixel.Sprite

	OffsetX float64
	OffsetY float64
	Scale   float64

	Health            float64
	HealthBarPosition pixel.Vec
	HealthBarImage    *image.RGBA
	HealthBarSprite   *pixel.Sprite

	Hunger            float64
	HungerBarPosition pixel.Vec
	HungerBarImage    *image.RGBA
	HungerBarSprite   *pixel.Sprite

	Thirst            float64
	ThirstBarPosition pixel.Vec
	ThirstBarImage    *image.RGBA
	ThirstBarSprite   *pixel.Sprite

	NeedsRedraw bool
}

func NewGUI(win *opengl.Window) (*GUI, error) {
	s, err := NewSpritesheet("./assets/gui.png")
	if err != nil {
		return nil, err
	}

	barSprite := pixel.NewSprite(s.Picture, pixel.R(0, s.Picture.Bounds().H(), 4*16, s.Picture.Bounds().H()-16))

	healthBarImage, healthBarSprite := MakeRect(46, 4, colornames.Red)
	hungerBarImage, hungerBarSprite := MakeRect(46, 4, colornames.Red)
	thirstBarImage, thirstBarSprite := MakeRect(46, 4, colornames.Red)

	g := &GUI{
		Window:      win,
		Spritesheet: s,
		BarSprite:   barSprite,
		OffsetX:     8 * 16,
		OffsetY:     2 * 16,
		Scale:       4.0,

		Health:          100,
		HealthBarImage:  healthBarImage,
		HealthBarSprite: healthBarSprite,

		Hunger:          100,
		HungerBarImage:  hungerBarImage,
		HungerBarSprite: hungerBarSprite,

		Thirst:          100,
		ThirstBarImage:  thirstBarImage,
		ThirstBarSprite: thirstBarSprite,
	}

	//g.HealthBarPosition = pixel.V(g.OffsetX, g.Window.Bounds().H()-g.OffsetY)
	g.HealthBarPosition = pixel.V(16, g.Window.Bounds().H())
	g.HungerBarPosition = pixel.V(16, g.Window.Bounds().H()-2*16)
	g.ThirstBarPosition = pixel.V(16, g.Window.Bounds().H()-4*16)

	return g, nil
}

func (g *GUI) Draw() {
	g.RedrawBars()
}

func (g *GUI) RedrawBars() {
	// health
	g.UpdateHealth(50)
	g.BarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.HealthBarPosition.Add(pixel.V(g.OffsetX, -g.OffsetY))))

	// hunger
	g.UpdateHunger(50)
	g.BarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.HungerBarPosition.Add(pixel.V(g.OffsetX, -g.OffsetY))))

	// thirst
	g.UpdateThirst(50)
	g.BarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.ThirstBarPosition.Add(pixel.V(g.OffsetX, -g.OffsetY))))
}

func (g *GUI) UpdateHealth(v float64) {
	amt := 0.48 * v

	for y := 0; y < 4; y++ {
		for x := 0; x < 46; x++ {
			if x <= int(amt) {
				g.HealthBarImage.Set(x, y, colornames.Red)
			} else {
				g.HealthBarImage.Set(x, y, pixel.RGBA{R: 0, G: 0, B: 0, A: 1})
			}
		}
	}

	g.HealthBarSprite = pixel.NewSprite(pixel.PictureDataFromImage(g.HealthBarImage), pixel.R(0, 0, 46, 4))

	g.HealthBarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.HealthBarPosition.Add(pixel.V(g.OffsetX-32, -g.OffsetY))))
}

func (g *GUI) UpdateHunger(v float64) {
	amt := 0.48 * v

	for y := 0; y < 4; y++ {
		for x := 0; x < 46; x++ {
			if x <= int(amt) {
				g.HungerBarImage.Set(x, y, colornames.Green)
			} else {
				g.HungerBarImage.Set(x, y, pixel.RGBA{R: 0, G: 0, B: 0, A: 1})
			}
		}
	}

	g.HungerBarSprite = pixel.NewSprite(pixel.PictureDataFromImage(g.HungerBarImage), pixel.R(0, 0, 46, 4))

	g.HungerBarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.HungerBarPosition.Add(pixel.V(g.OffsetX-32, -g.OffsetY))))
}

func (g *GUI) UpdateThirst(v float64) {
	amt := 0.48 * v

	for y := 0; y < 4; y++ {
		for x := 0; x < 46; x++ {
			if x <= int(amt) {
				g.ThirstBarImage.Set(x, y, colornames.Blue)
			} else {
				g.ThirstBarImage.Set(x, y, pixel.RGBA{R: 0, G: 0, B: 0, A: 1})
			}
		}
	}

	g.ThirstBarSprite = pixel.NewSprite(pixel.PictureDataFromImage(g.ThirstBarImage), pixel.R(0, 0, 46, 4))

	g.ThirstBarSprite.Draw(g.Window, pixel.IM.Moved(pixel.ZV).Scaled(pixel.ZV, g.Scale).Moved(g.ThirstBarPosition.Add(pixel.V(g.OffsetX-32, -g.OffsetY))))
}
