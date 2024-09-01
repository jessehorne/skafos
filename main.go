package main

import (
	"fmt"
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
		Title:  "Skafos v0.0.1",
		Bounds: pixel.R(0, 0, 900, 600),
		VSync:  true,
	}

	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	win.Clear(colornames.Black)

	// create new game
	g, err := game.NewGame("test", win)
	if err != nil {
		log.Fatalln(err)
	}

	g.Init(win)
	g.Map.RefreshDrawBatch()
	g.CollideablesDrawDebug = false

	maxFPS := float64(1 / 30)
	currentFrame := float64(0)
	last := time.Now()

	win.SetScrollCallback(func(win *opengl.Window, scroll pixel.Vec) {
		if scroll.Y == 1 {
			if game.Cam.Zoom < 42 {
				game.Cam.Zoom *= math.Pow(game.Cam.ZoomSpeed, scroll.Y)
			}
		} else {
			if game.Cam.Zoom > 1.6 {
				game.Cam.Zoom *= math.Pow(game.Cam.ZoomSpeed, scroll.Y)
			}
		}
	})

	win.SetButtonCallback(func(win *opengl.Window, button pixel.Button, action pixel.Action) {
		g.ButtonCallback(button, action)
	})

	win.SetCharCallback(func(win *opengl.Window, r rune) {
		g.CharCallback(r)
	})

	frames := 0
	second := time.Tick(time.Second)

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

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	opengl.Run(run)
}
