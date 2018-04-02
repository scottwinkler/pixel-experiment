package tilemap

import (
	"github.com/faiface/pixel"
)

type TilemapLayer struct {
	Data    []int   `json:"data"`
	Name    string  `json:"name"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
	Width   int     `json:"width"`  //number of tiles wide
	Height  int     `json:"height"` //number of tiles high
	Opacity float64 `json:"opacity"`
	Visible bool    `json:"visible"`
	Tiles   []*Tile
	Tilemap *Tilemap
}

//insantiate tiles from tileset data
func (layer *TilemapLayer) MakeTiles(tm *Tilemap) *TilemapLayer {
	layer.Tilemap = tm
	tileHeight := tm.TileHeight
	tileWidth := tm.TileWidth
	var tiles []*Tile
	index := 0
	for col := 0; col < layer.Width; col++ {
		for row := 0; row < layer.Height; row++ {
			x := row*tileHeight + tileHeight/2
			y := -col*tileWidth - tileWidth/2 + tileHeight*tm.Height //put it top right quadrant

			var sprite *pixel.Sprite
			var matrix pixel.Matrix = pixel.IM.Moved(pixel.V(float64(x), float64(y)))
			var properties map[string]interface{}
			var tile *Tile
			for _, tileset := range tm.Tilesets {
				gid := layer.Data[index]
				if gid < 1 { //quit early
					break
				}

				if tileset.Contains(gid) {
					//fmt.Printf("gid: %d in tileset", gid)
					sprite = tileset.GetSpriteForGid(gid)
					properties = tileset.GetPropertiesForGid(gid)
					break
				}
			}
			tile = NewTile(layer, x, y, tileWidth, tileHeight, sprite, matrix, properties)
			//fmt.Printf("new tile: %v", tile)
			tiles = append(tiles, tile)
			index++
		}
	}
	layer.Tiles = tiles
	return layer
}

func (layer *TilemapLayer) Draw(t pixel.Target) {
	for _, tile := range layer.Tiles {
		tile.Draw(t)
	}

}
