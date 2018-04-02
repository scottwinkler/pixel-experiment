package tilemap

import (
	"github.com/faiface/pixel"
)

type Tile struct {
	IsCollidable bool
	Height       int           //the height of the tile in pixels
	Width        int           //the width of the tile in pixels
	Layer        *TilemapLayer //The layer in the Tilemap data that this tile belongs to
	X            int           //the x coordinate of this tile
	Y            int           //the y coordinate of this tile
	Sprite       *pixel.Sprite
	Matrix       pixel.Matrix
	Properties   map[string]interface{}
	Alpha        pixel.RGBA
}

// NewTile creates a Tile from the
func NewTile(layer *TilemapLayer, x int, y int, width int, height int, sprite *pixel.Sprite, matrix pixel.Matrix, properties map[string]interface{}) *Tile {
	isCollidable := false
	if properties != nil {
		value := properties["collidable"]
		if value != nil {
			isCollidable = value.(bool)
		}
	}
	tile := Tile{
		Layer:        layer,
		X:            x,
		Y:            y,
		Width:        width,
		Height:       height,
		Sprite:       sprite,
		Matrix:       matrix,
		Properties:   properties,
		IsCollidable: isCollidable,
		Alpha:        pixel.Alpha(layer.Opacity),
	}
	return &tile
}

/*
func (Tile) ContainsPoint(x int, y int) bool {

}
*/
func (tile *Tile) Draw(t pixel.Target) {
	if tile.Sprite != nil {
		tile.Sprite.DrawColorMask(t, tile.Matrix, tile.Alpha)
	}
}
