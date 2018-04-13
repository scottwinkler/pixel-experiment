package tilemap

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"strconv"

	"github.com/faiface/pixel"

	"github.com/scottwinkler/simple-rpg/utility"
)

type Tileset struct {
	//these fields need to be public because golang does not support reflection
	FirstGid    int `json:"firstgid"`
	lastGid     int
	Source      string `json:"source"`
	tilesetData TilesetData
	batch       *pixel.Batch
}

type TilesetData struct {
	//these fields need to be public because golang does not support reflection
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
	picture     *pixel.Picture
	TileCount   int                    `json:"tilecount"`
	Properties  map[string]interface{} `json:"tileproperties"`
}

const TMX_DIR = "_assets/tmx/"

//fetchs the rest of the tileset data
//todo: allow it to read from sets of images, not just a single image. Will need to change the datastructure of TilesetData to allow this
func (ts *Tileset) FetchTilesetData() {
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
	tilesetData.picture = &pic
	ts.tilesetData = tilesetData
	ts.lastGid = ts.FirstGid + tilesetData.TileCount - 1
	ts.batch = pixel.NewBatch(&pixel.TrianglesData{}, *ts.tilesetData.picture)
}

//returns true if the specified gid is in the tileset
func (ts *Tileset) Contains(gid int) bool {
	return gid >= ts.FirstGid && gid <= ts.lastGid
}

//get the properties for the gid
func (ts *Tileset) GidToProperties(gid int) map[string]interface{} {
	var properties map[string]interface{}
	value := ts.tilesetData.Properties[strconv.Itoa(gid-1)] //something is weird about the csv data. the value is one more than it should be
	if value != nil {
		properties = value.(map[string]interface{})
	}
	return properties
}

//returns a sprite for a specified gid
func (ts *Tileset) GidToSprite(gid int) *pixel.Sprite {
	pic := ts.tilesetData.picture
	pd := pixel.PictureDataFromPicture(*pic)
	index := gid - ts.FirstGid
	col := (index % ts.tilesetData.Columns)                                       //if index=0 then i am in col=0
	row := int(math.Ceil((float64(index)+1)/float64(ts.tilesetData.Columns))) - 1 //if index=0 then i am in row=0
	minX := pd.Bounds().Min.X + float64(col*ts.tilesetData.TileWidth)
	maxX := minX + float64(ts.tilesetData.TileWidth)
	maxY := pd.Bounds().Max.Y - float64(row*ts.tilesetData.TileHeight)
	minY := maxY - float64(ts.tilesetData.TileHeight)
	frame := pixel.R(minX, minY, maxX, maxY)
	sprite := pixel.NewSprite(*pic, frame)
	return sprite
}
