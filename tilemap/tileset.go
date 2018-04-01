package tilemap

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/faiface/pixel"

	"github.com/scottwinkler/pixel-experiment/utility"
)

type Tileset struct {
	FirstGid    int `json:"firstgid"`
	LastGid     int
	Source      string `json:"source"`
	TilesetData TilesetData
}

type TilesetData struct {
	Image       string `json:"image"`
	Source      string `json:"source"`
	Name        string `json:"name"`
	TileWidth   int    `json:"tilewidth"`
	TileHeight  int    `json:"tileheight"`
	ImageWidth  int    `json:"imagewidth"`
	ImageHeight int    `json:"imageheight"`
	Margin      int    `json:"margin"`
	Spacing     int    `json:"spacing"`
	Columns     int    `json:"columns"`
	Type        string `json:"type"`
	Picture     *pixel.Picture
	TileCount   int                    `json:"tilecount"`
	Properties  map[string]interface{} `json:"tileproperties"`
}

const TMX_DIR = "assets/tmx/"

//fetchs the rest of the tileset data
func (ts Tileset) FetchTilesetData() Tileset {
	//fmt.Println("reading data source file")
	raw, err := ioutil.ReadFile(TMX_DIR + ts.Source)
	if err != nil {
		panic(err)
	}
	var tilesetData TilesetData
	json.Unmarshal(raw, &tilesetData)

	pic, err := utility.LoadPicture(TMX_DIR + tilesetData.Image)
	if err != nil {
		panic(err)
	}
	tilesetData.Picture = &pic
	ts.TilesetData = tilesetData
	ts.LastGid = ts.FirstGid + tilesetData.TileCount
	//fmt.Printf("%+v\n", ts)
	return ts
}

//returns true if the specified gid is in the tileset
func (ts Tileset) Contains(gid int) bool {
	return gid >= ts.FirstGid && gid <= ts.LastGid
}

func (ts Tileset) GetPropertiesForGid(gid int) map[string]interface{} {
	var properties map[string]interface{}
	value := ts.TilesetData.Properties[strconv.Itoa(gid)]
	if value != nil {
		properties = value.(map[string]interface{})
	}
	return properties
}

//returns a sprite for a specified gid
func (ts Tileset) GetSpriteForGid(gid int) *pixel.Sprite {
	pic := ts.TilesetData.Picture
	pd := pixel.PictureDataFromPicture(*pic)
	index := gid - ts.FirstGid
	col := index % ts.TilesetData.Columns
	row := (index - col) / ts.TilesetData.Columns
	//fmt.Printf("index: %d, col: %d,row: %d", index, col, row)
	minX := pd.Bounds().Min.X + float64(col*ts.TilesetData.TileWidth)
	maxX := minX + float64(ts.TilesetData.TileWidth)
	maxY := pd.Bounds().Max.Y - float64(row*ts.TilesetData.TileHeight)
	minY := maxY - float64(ts.TilesetData.TileHeight)
	frame := pixel.R(minX, minY, maxX, maxY)
	//fmt.Printf("frame for gid %d is: %v", gid, frame)
	sprite := pixel.NewSprite(*pic, frame)
	return sprite
}

/*
func (tileset Tileset) Draw(t pixel.Target) {
	for i, tile := range tileset.tiles {
		tile.spritePtr.Draw(t, tileset.tiles[i].matrix)
	}
}*/