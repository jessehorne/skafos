package game

import (
	"github.com/gopxl/pixel/v2"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type Spritesheet struct {
	Path    string
	Picture pixel.Picture
}

func NewSpritesheet(path string) (*Spritesheet, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return &Spritesheet{
		Path:    path,
		Picture: pixel.PictureDataFromImage(img),
	}, nil
}
