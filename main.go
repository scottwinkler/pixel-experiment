package main

import (
	_ "image/png"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Tile struct {
	sprite *pixel.Sprite
	matrix pixel.Matrix
}

type Tileset struct {
	tiles []Tile
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 800, 800),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	tileSpritesheet := loadSpritesheet("spritesheets/minecraft_tiles.png", 80, 80)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, tileSpritesheet.picture)
	var (
		camPos   = pixel.ZV
		camSpeed = 10.0
		camZoom  = 1.0
		tiles    []Tile
	)

	last := time.Now()
	for x := 0; x < int(win.Bounds().Max.X); x += 80 {
		for y := 0; y < int(win.Bounds().Max.Y); y += 80 {

			sprite := tileSpritesheet.sprites[rand.Intn(len(tileSpritesheet.sprites))]
			pos := pixel.V(float64(x), float64(y))
			matrix := pixel.IM.Scaled(pixel.ZV, 1).Moved(pos)
			tile := Tile{sprite: sprite, matrix: matrix}
			tiles = append(tiles, tile)
		}
	}

	for !win.Closed() {
		dt := time.Since(last).Seconds()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)
		/*if win.JustPressed(pixelgl.MouseButtonLeft) {
			tile := pixel.NewSprite(spritesheet, tileFrames[rand.Intn(len(tileFrames))])
			tiles = append(tiles, tile)
			mouse := cam.Unproject(win.MousePosition())
			matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 4).Moved(mouse))
		}*/

		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		win.Clear(colornames.White)
		batch.Clear()
		for i, tile := range tiles {
			tile.sprite.Draw(batch, tiles[i].matrix)
		}
		batch.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
