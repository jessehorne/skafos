package game

import (
	"fmt"
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

	img, de, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	fmt.Println(de)

	return &Spritesheet{
		Path:    path,
		Picture: pixel.PictureDataFromImage(img),
	}, nil
}
