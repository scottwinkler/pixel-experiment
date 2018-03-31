package main

import (
	"image"
	"os"

	"github.com/faiface/pixel"
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

type Spritesheet struct {
	sprites []*pixel.Sprite
	picture pixel.Picture
}

//parse a spritesheet based on a standard width and height
func loadSpritesheet(path string, frameWidth float64, frameHeight float64) Spritesheet {
	spritesheet, err := loadPicture(path)
	if err != nil {
		panic(err)
	}
	var sprites []*pixel.Sprite
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += frameWidth {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += frameHeight {
			frame := pixel.R(x, y, x+frameWidth, y+frameHeight)
			sprites = append(sprites, pixel.NewSprite(spritesheet, frame))
		}
	}
	return Spritesheet{sprites: sprites, picture: spritesheet}
}
