package tilemap

import (
	"github.com/faiface/pixel"
)

type TilemapLayer struct {
	//these fields need to be public because golang does not support reflection
	Data        []int   `json:"data"`
	Name        string  `json:"name"`
	X           int     `json:"x"`
	Y           int     `json:"y"`
	Width       int     `json:"width"`  //number of tiles wide
	Height      int     `json:"height"` //number of tiles high
	Opacity     float64 `json:"opacity"`
	Visible     bool    `json:"visible"`
	tiles       []*Tile
	tilemap     *Tilemap
	batchGroups map[string][]*Tile //a useful mapping of tileset names to tiles, for more efficient drawing
}

//insantiate tiles from tileset data
func (l *TilemapLayer) MakeTiles(tm *Tilemap) *TilemapLayer {
	tileSetGroups := make(map[string][]*Tile)
	for _, ts := range tm.Tilesets {
		tileSetGroups[ts.tilesetData.Name] = []*Tile{}
	}
	l.tilemap = tm
	tileHeight := tm.TileHeight
	tileWidth := tm.TileWidth
	var tiles []*Tile
	index := 0
	gid := 0
	for col := 0; col < l.Width; col++ {
		for row := 0; row < l.Height; row++ {
			x := row*tileHeight + tileHeight/2
			y := -col*tileWidth - tileWidth/2 + tileHeight*tm.Height //put it top right quadrant
			var sprite *pixel.Sprite
			var matrix pixel.Matrix = pixel.IM.Moved(pixel.V(float64(x), float64(y)))
			var properties map[string]interface{}
			var tile *Tile
			var tileSetName string
			for _, tileset := range tm.Tilesets {
				gid = l.Data[index]
				if gid < 1 { //quit early
					break
				}

				if tileset.Contains(gid) {
					tileSetName = tileset.tilesetData.Name
					if l.Visible { //necessary to hide collision tiles
						sprite = tileset.GidToSprite(gid)
					}
					properties = tileset.GidToProperties(gid)
					break
				}
			}

			tile = NewTile(l, pixel.V(float64(x), float64(y)), gid, tileWidth, tileHeight, sprite, matrix, properties)
			if tileSetName != "" {
				tileSetGroup := tileSetGroups[tileSetName]
				tileSetGroup = append(tileSetGroup, tile) //add tile to tileset group for future reference
				tileSetGroups[tileSetName] = tileSetGroup
			}
			tiles = append(tiles, tile)
			index++
		}
	}
	l.batchGroups = tileSetGroups
	l.tiles = tiles
	return l
}

func (l *TilemapLayer) Draw(t pixel.Target) {

	for name, tiles := range l.batchGroups {
		ts := l.tilemap.GetTileset(name)

		batch := ts.batch
		batch.Clear()
		for _, tile := range tiles {
			tile.Draw(batch)
		}
		batch.Draw(t)
	}

}
