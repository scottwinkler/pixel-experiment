package tilemap

import (
	"github.com/faiface/pixel"
)

type Tilemap struct {
	Width       int            `json:"width"`
	Height      int            `json:"height"`
	TileWidth   int            `json:"tilewidth"`
	TileHeight  int            `json:"tileheight"`
	Orientation string         `json:"orientation"`
	Version     string         `json:"version"`
	Layers      []TilemapLayer `json:"layers"`
	Tilesets    []Tileset      `json:"tilesets"`
}

func (tm *Tilemap) MakeTiles() {
	var new_layers []TilemapLayer
	for _, layer := range tm.Layers {
		new_layers = append(new_layers, *layer.MakeTiles(tm))
	}
	tm.Layers = new_layers
}

func (tm *Tilemap) Draw(t pixel.Target) {
	for _, layer := range tm.Layers {
		layer.Draw(t)
	}
}
