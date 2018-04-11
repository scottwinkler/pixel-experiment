package utility

import (
	"strings"

	"github.com/faiface/pixel"
)

type Spritesheet struct {
	sprites []*pixel.Sprite
	matrix  pixel.Matrix
}

func (s *Spritesheet) Matrix() pixel.Matrix {
	return s.matrix
}

func (s *Spritesheet) Sprites() []*pixel.Sprite {
	return s.sprites
}

//returns an empty spritesheet
func NewSpritesheet() *Spritesheet {
	var sprites []*pixel.Sprite
	spritesheet := &Spritesheet{
		sprites: sprites,
		matrix:  pixel.IM,
	}
	return spritesheet
}

//helper method, used with addsprite
func LoadSprite(path string) *pixel.Sprite {
	pic, err := LoadPicture(path)
	if err != nil {
		panic(err)
	}
	return pixel.NewSprite(pic, pic.Bounds())
}

//useful if choosing to create a spritesheet from individual sprite
func (s *Spritesheet) AddSprite(sprite *pixel.Sprite) {
	s.sprites = append(s.sprites, sprite)
}

//parse a spritesheet based on a standard width and height
func LoadSpritesheet(path string, frame pixel.Rect, scale float64) *Spritesheet {
	parts := strings.Split(path, ".")
	ext := parts[len(parts)-1]
	var sprites []*pixel.Sprite
	//process a png spritesheet with multiple frames
	if strings.EqualFold(ext, "png") {
		pic, err := LoadPicture(path)
		if err != nil {
			panic(err)
		}
		for y := pic.Bounds().Max.Y; y > pic.Bounds().Min.Y; y -= frame.H() {
			for x := pic.Bounds().Min.X; x < pic.Bounds().Max.X; x += frame.W() {
				frame := pixel.R(x, y-frame.H(), x+frame.W(), y)
				sprites = append(sprites, pixel.NewSprite(pic, frame))
			}
		}
	} else if strings.EqualFold(ext, "gif") { //process a gif spritesheet
		gif, _ := LoadGif(path)

		for _, pic := range gif {
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprites = append(sprites, sprite)
		}
	}
	matrix := pixel.IM.Scaled(pixel.V(0, 0), scale)
	return &Spritesheet{sprites: sprites, matrix: matrix}
}
