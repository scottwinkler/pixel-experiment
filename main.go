package main

import (
	"fmt"
	_ "image/png"
	"time"

	"github.com/scottwinkler/pixel-experiment/player"
	"github.com/scottwinkler/pixel-experiment/spritesheet"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/pixel-experiment/tilemap"
	"golang.org/x/image/colornames"
)

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
	tm, err := tilemap.ParseTiledJSON("assets/tmx/world2.json")
	tm.MakeWorld()
	spritesheet := spritesheet.LoadSpritesheet("assets/spritesheets/baldricSlashSheet.png", 64, 64)
	player := player.NewPlayer(0, spritesheet)
	fmt.Printf("len: %d", len(spritesheet.Sprites))
	var (
		camPos   = pixel.ZV
		camSpeed = 1.0
		camZoom  = 1.0
		//tiles    []Tile
	)

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
			player.SetFrame(2)
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
			player.SetFrame(0)
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
			player.SetFrame(1)
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
			player.SetFrame(3)
		}
		win.Clear(colornames.Black)
		//batch.Clear()
		tm.Draw(win)
		//spritesheet.Sprites[0].Draw(win, pixel.IM)
		player.Draw(win)
		win.SetTitle(cam.Unproject(win.MousePosition()).String())
		//batch.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
