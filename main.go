package main

import (
	"fmt"
	_ "image/png"
	"time"

	"github.com/scottwinkler/pixel-experiment/animation"
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
	win.SetTitle("Aww yisss!!!")
	tm, err := tilemap.ParseTiledJSON("assets/tmx/world2.json")
	//tm.MakeWorld()
	//world := world.NewWorld()
	//world.SetTilemap(tm)
	spritesheet := spritesheet.LoadSpritesheet("assets/spritesheets/baldricSlashSheet.png", 64, 64)
	var animations []*animation.Animation
	animations = append(animations, animation.NewAnimation(spritesheet, "WalkRight", []int{0, 4, 8, 12, 16, 20, 16, 12, 8, 4, 0}, true))
	animations = append(animations, animation.NewAnimation(spritesheet, "WalkDown", []int{1, 5, 9, 13, 17, 21, 17, 13, 9, 5, 1}, true))
	animations = append(animations, animation.NewAnimation(spritesheet, "WalkLeft", []int{2, 6, 10, 14, 18, 22, 18, 14, 10, 6, 2}, true))
	animations = append(animations, animation.NewAnimation(spritesheet, "WalkUp", []int{3, 7, 11, 15, 19, 23, 19, 15, 11, 7}, true))
	player := player.NewPlayer(animations, win)
	//world.SetPlayer(pla)
	fmt.Printf("len: %d", len(spritesheet.Sprites))
	//var (
	//	camPos   = pixel.ZV
	//	camSpeed = 1.0
	//	camZoom  = 1.0
	//tiles    []Tile
	//second = time.Tick(time.Second)
	//)

	//last := time.Now()
	fps := 60
	times := 0
	interval := time.Duration(float64(1000) / float64(fps))
	//fmt.Printf("interval: %d", interval)
	ticker := time.NewTicker(time.Millisecond * interval)
	quit := make(chan struct{})
	go func() {
		//for t := range ticker.C {
		for {
			select {
			case <-ticker.C:

				win.Clear(colornames.Black)
				tm.Draw(win)
				//if(times)
				player.Update(times)
				win.Update()
				/*if win.Closed() {
					ticker.Stop()
				}*/
				times++
				if times > 3 {
					times = 0
				}
				if win.Closed() {
					close(quit)
				}
			case <-quit:
				ticker.Stop()
				return

				//}
			}
		}
	}()
	for !win.Closed() {
		time.Sleep(time.Millisecond * interval)
	}
	//ticker.Stop()
}

func main() {
	pixelgl.Run(run)
}
