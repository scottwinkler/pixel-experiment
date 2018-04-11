package utility

import (
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"io/ioutil"
	"os"

	"github.com/faiface/pixel"
)

func LoadFile(path string) *os.File {
	//cwd, _ := os.Getwd()
	//fmt.Printf("cwd %s", cwd)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return file
}

func LoadJSON(path string) map[string]interface{} {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var f interface{}
	json.Unmarshal(raw, &f)
	return f.(map[string]interface{})
}

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
