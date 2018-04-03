package tilemap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/faiface/pixel"
)

func ParseTiledJSON(path string) (*Tilemap, error) {
	raw, err := ioutil.ReadFile(path)
	var tilemap Tilemap
	if err != nil {
		fmt.Printf("%s", err)
		return &tilemap, err
	}
	json.Unmarshal(raw, &tilemap)
	if tilemap.Orientation != "orthogonal" {
		fmt.Println("Only orthogonal map types are supported")
		return &tilemap, nil
	}
	var new_tilesets []Tileset
	for _, tileset := range tilemap.Tilesets {
		tileset.FetchTilesetData()
		new_tilesets = append(new_tilesets, tileset)
	}
	tilemap.Tilesets = new_tilesets
	tilemap.MakeTiles()
	tilemap.SetBounds(pixel.R(0, 0, float64(tilemap.TileWidth*tilemap.Width), float64(tilemap.TileHeight*tilemap.Height)))
	return &tilemap, nil
}
