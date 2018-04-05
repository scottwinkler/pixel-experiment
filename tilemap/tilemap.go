package tilemap

import (
	"strings"

	"github.com/faiface/pixel"
)

type Tilemap struct {
	TileWidth       int            `json:"tilewidth"`
	TileHeight      int            `json:"tileheight"`
	Orientation     string         `json:"orientation"`
	Version         int            `json:"version"`
	Layers          []TilemapLayer `json:"layers"`
	Tilesets        []Tileset      `json:"tilesets"`
	Width           int            `json:"width"`
	Height          int            `json:"height"`
	bounds          pixel.Rect
	LayersNameIndex map[string]int //a convenient datastructure for mapping a name of a layer to its index
}

// initialize the tiles in the tilemap. Only needs to be run once, should never be called directly
func (tm *Tilemap) MakeTiles() {
	var new_layers []TilemapLayer
	for _, layer := range tm.Layers {
		new_layers = append(new_layers, *layer.MakeTiles(tm))
	}
	tm.Layers = new_layers
}

//useful for z-indexing
func (tm *Tilemap) DrawLayers(t pixel.Target, names []string) {
	for _, name := range names {
		index := tm.LayersNameIndex[name]
		tm.Layers[index].Draw(t)
	}
}

//setter method for Bounds
func (tm *Tilemap) SetBounds(bounds pixel.Rect) {
	tm.bounds = bounds
}

//getter method for Bounds
func (tm *Tilemap) Bounds() pixel.Rect {
	return tm.bounds
}

//retreive a tileset by its name. used by layers to make batch groups for tiles of the same tileset
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
func (tm *Tilemap) GetTileAtPosition(pos pixel.Vec, layerName string) *Tile {
	//fmt.Printf("len, %v", tm.LayersNameIndex)
	layerIndex := tm.LayersNameIndex[layerName]
	//fmt.Printf("len, %d", layerIndex)
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

//delegate drawing responsbility to layers
func (tm *Tilemap) Draw(t pixel.Target) {
	for _, layer := range tm.Layers {
		layer.Draw(t)
	}
}
