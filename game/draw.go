package game

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"image"
)

var (
	DrawColor        = pixel.RGBA{R: 1, G: 1, B: 1}
	RectangleSprites = map[int]map[int]*pixel.Sprite{}
)

func DrawSetColor(color pixel.RGBA) {
	DrawColor = color
}

// HLine draws a horizontal line
func DrawHLine(img *image.RGBA, x1, y, x2 int) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, DrawColor)
	}
}

// VLine draws a veritcal line
func DrawVLine(img *image.RGBA, x, y1, y2 int) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, DrawColor)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func MakeDebugRect(win *opengl.Window, w, h int) *pixel.Sprite {
	_, yExists := RectangleSprites[h]
	if yExists {
		_, xExists := RectangleSprites[h][w]
		if xExists {
			return RectangleSprites[h][w]
		}
	}

	x1 := 0
	y1 := 0
	x2 := w
	y2 := h

	i := image.NewRGBA(image.Rect(x1, y1, x2+2, y2+2))
	DrawHLine(i, x1, y1, x2)
	DrawHLine(i, x1, y2, x2)
	DrawVLine(i, x1, y1, y2)
	DrawVLine(i, x2, y1, y2)

	_, yExists = RectangleSprites[h]
	if !yExists {
		RectangleSprites[h] = map[int]*pixel.Sprite{}
	}

	_, xExists := RectangleSprites[h][w]

	if !xExists {
		RectangleSprites[h][w] = pixel.NewSprite(pixel.PictureDataFromImage(i), pixel.R(float64(x1), float64(y1), float64(x2)+2, float64(y2)+2))
	}

	return RectangleSprites[h][w]
}
