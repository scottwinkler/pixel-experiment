package main

import (
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
	tm, err := tilemap.ParseTiledJSON("assets/tmx/world3.json")
	//tm.MakeWorld()
	//world := world.NewWorld()
	//world.SetTilemap(tm)
	/*spritesheet := spritesheet.LoadSpritesheet("assets/spritesheets/baldricSlashSheet.png", 64, 64)
	var animations []*animation.Animation
	animations = append(animations, animation.NewAnimation(spritesheet, "WalkRight", []int{0, 4, 8, 12, 16, 20, 16, 12, 8, 4, 0}, true))
	animations = append(animations, animation.NewAnimation(spritesheet, "WalkDown", []int{1, 5, 9, 13, 17, 21, 17, 13, 9, 5, 1}, true))
	animations = append(animations, animation.NewAnimation(spritesheet, "WalkLeft", []int{2, 6, 10, 14, 18, 22, 18, 14, 10, 6, 2}, true))
	animations = append(animations, animation.NewAnimation(spritesheet, "WalkUp", []int{3, 7, 11, 15, 19, 23, 19, 15, 11, 7}, true))*/
	//spritesheet := spritesheet.NewSpritesheet()
	spritesheet := spritesheet.LoadSpritesheet("assets/spritesheets/knight_iso_char.png", 84, 84)
	var animations []*animation.Animation
	animations = append(animations, animation.NewAnimation(spritesheet, "Idle", []int{0, 1, 2, 3}))
	animations = append(animations, animation.NewAnimation(spritesheet, "RunDown", []int{4, 5, 6, 7, 8}))
	animations = append(animations, animation.NewAnimation(spritesheet, "RunUp", []int{9, 10, 11, 12, 13}))
	animations = append(animations, animation.NewAnimation(spritesheet, "RunRight", []int{14, 15, 16, 17, 18, 19}))
	animations = append(animations, animation.NewAnimation(spritesheet, "RunLeft", []int{20, 21, 22, 23, 24, 25}))
	animations = append(animations, animation.NewAnimation(spritesheet, "AttackDown", []int{26, 27, 28}))
	animations = append(animations, animation.NewAnimation(spritesheet, "AttackUp", []int{29, 30, 31}))
	animations = append(animations, animation.NewAnimation(spritesheet, "AttackRight", []int{32, 33, 34}))
	animations = append(animations, animation.NewAnimation(spritesheet, "AttackLeft", []int{35, 36, 37}))
	player := player.NewPlayer(animations, win)

	fps := 60
	ticks := 0
	interval := time.Duration(float64(1000) / float64(fps))
	ticker := time.NewTicker(time.Millisecond * interval)
	go func() {
		for {
			select {
			case <-ticker.C: //main game loop @normalized fps is here
				win.Clear(colornames.Black)
				tm.Draw(win)
				player.Update(ticks)
				win.Update()
				ticks++
				//assume 60 ticks per second
				//so change animation once every 6 ticks to achieve a frameRate of 15, which feels reasonable
				if ticks > 6 {
					ticks = 0
				}
			}
		}
	}()
	//need this otherwise the game exits immediantly
	for !win.Closed() {
		time.Sleep(time.Millisecond * interval)
	}
	ticker.Stop()
}

func main() {
	pixelgl.Run(run)
}
