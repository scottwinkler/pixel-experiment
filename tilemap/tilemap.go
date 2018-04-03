package tilemap

import (
	"strings"

	"github.com/faiface/pixel"
)

type Tilemap struct {
	TileWidth   int            `json:"tilewidth"`
	TileHeight  int            `json:"tileheight"`
	Orientation string         `json:"orientation"`
	Version     int            `json:"version"`
	Layers      []TilemapLayer `json:"layers"`
	Tilesets    []Tileset      `json:"tilesets"`
	Width       int            `json:"width"`
	Height      int            `json:"height"`
	bounds      pixel.Rect
}

func (tm *Tilemap) MakeTiles() {
	var new_layers []TilemapLayer
	for _, layer := range tm.Layers {
		new_layers = append(new_layers, *layer.MakeTiles(tm))
	}
	tm.Layers = new_layers
}

//useful for z-indexing
func (tm *Tilemap) DrawLayers(t pixel.Target, layerIndexes []int) {
	for _, index := range layerIndexes {
		tm.Layers[index].Draw(t)
	}
}

func (tm *Tilemap) SetBounds(bounds pixel.Rect){
	tm.bounds = bounds
}
func (tm *Tilemap) Bounds() pixel.Rect {
	return tm.bounds
}

func (tm *Tilemap) GetTileset(name string) *Tileset {
	var tileset *Tileset
	for _, ts := range tm.Tilesets {
		if strings.EqualFold(ts.TilesetData.Name, name) {
			return &ts
		}
	}
	return tileset //should never happen
}

//accepts a coodinate position and a layer index and returns the tile at the position, or null
func (tm *Tilemap) GetTileAtPosition(pos pixel.Vec, layerIndex int) *Tile {
	var tile *Tile
	x := int(pos.X)
	y := int(pos.Y)
	maxY := tm.TileHeight * tm.Height
	maxX := tm.TileWidth * tm.Width
	if x <= 0 || y <= 0 || x >= maxX || y >= maxY {
		return tile
	}
	col := x / tm.TileWidth
	row := (maxY - y) / tm.TileHeight
	tileIndex := row*tm.Width + col
	if tileIndex > tm.Height*tm.Width {
		return tile
	}
	//fmt.Printf("layer: %d, name: %s", layerIndex, tm.Layers[layerIndex].Name)
	return tm.Layers[layerIndex].Tiles[tileIndex]
}

func (tm *Tilemap) Draw(t pixel.Target) {
	for _, layer := range tm.Layers {
		layer.Draw(t)
	}
}
