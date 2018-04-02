package spritesheet

import (
	"fmt"
	"strings"

	"github.com/faiface/pixel"
	"github.com/scottwinkler/pixel-experiment/utility"
)

type Spritesheet struct {
	Sprites  []*pixel.Sprite
	Pictures []*pixel.Picture //typically do not need this
}

//parse a spritesheet based on a standard width and height
func LoadSpritesheet(path string, frameWidth int, frameHeight int) *Spritesheet {
	parts := strings.Split(path, ".")
	ext := parts[len(parts)-1]
	var sprites []*pixel.Sprite
	var pics []*pixel.Picture
	//process a png spritesheet with multiple frames
	if strings.EqualFold(ext, "png") {
		pic, err := utility.LoadPicture(path)
		if err != nil {
			panic(err)
		}

		for x := pic.Bounds().Min.X; x < pic.Bounds().Max.X; x += float64(frameWidth) {
			for y := pic.Bounds().Min.Y; y < pic.Bounds().Max.Y; y += float64(frameHeight) {
				frame := pixel.R(x, y, x+float64(frameWidth), y+float64(frameHeight))
				sprites = append(sprites, pixel.NewSprite(pic, frame))
				pics = append(pics, &pic)
			}
		}
	} else if strings.EqualFold(ext, "gif") { //process a gif spritesheet
		gif, _ := utility.LoadGif(path)

		for _, pic := range gif {
			pics = append(pics, &pic)
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprites = append(sprites, sprite)
		}
	}
	fmt.Printf("len: %d", len(sprites))
	return &Spritesheet{Sprites: sprites, Pictures: pics}
}

//gifs will have more than one picture
/*func (s Spritesheet) IsGif() bool {
	return len(s.Pictures) > 1
}*/
