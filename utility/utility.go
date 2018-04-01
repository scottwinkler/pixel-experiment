package utility

import (
	"image"
	"image/gif"
	"os"

	"github.com/faiface/pixel"
)

func LoadPicture(path string) (pixel.Picture, error) {
	//fmt.Println("printing cwd")
	//fmt.Println(os.Getwd())
	//fmt.Println("picture path: " + path)
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

func LoadGif(path string) ([]pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	gif, err := gif.DecodeAll(file)
	if err != nil {
		return nil, err
	}
	var pics []pixel.Picture
	for _, img := range gif.Image {
		pics = append(pics, pixel.PictureDataFromImage(img))
	}
	return pics, nil
}
