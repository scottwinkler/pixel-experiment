package tilemap

import (
	"github.com/faiface/pixel"
)

type Tile struct {
	collidable bool
	height     int           //the height of the tile in pixels
	width      int           //the width of the tile in pixels
	layer      *TilemapLayer //The layer in the Tilemap data that this tile belongs to
	v          pixel.Vec     //the draw position
	gid        int
	sprite     *pixel.Sprite
	matrix     pixel.Matrix
	properties map[string]interface{}
	alpha      pixel.RGBA
}

func (t *Tile) Collidable() bool {
	return t.collidable
}

// NewTile creates a Tile from the given data
func NewTile(layer *TilemapLayer, v pixel.Vec, width int, height int, gid int, sprite *pixel.Sprite, matrix pixel.Matrix, properties map[string]interface{}) *Tile {
	collidable := false
	if properties != nil {
		value := properties["collidable"]
		if value != nil {
			collidable = value.(bool)
		}
	}
	tile := Tile{
		layer:      layer,
		v:          v,
		gid:        gid,
		width:      width,
		height:     height,
		sprite:     sprite,
		matrix:     matrix,
		properties: properties,
		collidable: collidable,
		alpha:      pixel.Alpha(layer.Opacity),
	}
	return &tile
}

func (tile *Tile) Draw(t pixel.Target) {
	if tile.sprite != nil {
		tile.sprite.DrawColorMask(t, tile.matrix, tile.alpha)
	}
}
