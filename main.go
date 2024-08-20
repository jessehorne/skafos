package main

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/jessehorne/skafos/game"
	"golang.org/x/image/colornames"
	"image"
	"log"
	"math"
	"os"
	"time"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := opengl.WindowConfig{
		Title:  "Skafos (pre-alpha) by JesseH",
		Bounds: pixel.R(0, 0, 1600, 900),
		VSync:  true,
	}

	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	win.Clear(colornames.Black)

	// create new game
	g, err := game.NewGame("test")
	if err != nil {
		log.Fatalln(err)
	}

	for y := 0; y < 4; y++ {
		for x := 0; x < 7; x++ {
			g.Map.GenerateAllDirtChunk(x, y, true)
		}
	}

	g.Map.RefreshDrawBatch()

	maxFPS := float64(1 / 30)
	currentFrame := float64(0)
	last := time.Now()

	win.SetScrollCallback(func(win *opengl.Window, scroll pixel.Vec) {
		g.Camera.Zoom *= math.Pow(g.Camera.ZoomSpeed, scroll.Y)
	})

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		currentFrame += dt

		// only draw 30 fps
		if currentFrame >= maxFPS {
			win.Clear(colornames.Black)
			g.Draw(win)
			currentFrame = 0
		}

		g.Update(win, dt)
		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
