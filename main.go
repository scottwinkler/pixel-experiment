package main

import (
	_ "image/png"
	"time"

	"github.com/scottwinkler/pixel-experiment/animation"
	"github.com/scottwinkler/pixel-experiment/world"

	"github.com/scottwinkler/pixel-experiment/player"
	"github.com/scottwinkler/pixel-experiment/spritesheet"

	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/pixel-experiment/tilemap"
	"golang.org/x/image/colornames"
)

func run() {

	tm, _ := tilemap.ParseTiledJSON("assets/tmx/world1.json")

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

	world := world.NewWorld(400, 400)
	world.SetTilemap(tm)
	//v := pixel.V(float64(tm.TileWidth*tm.Width), float64(tm.TileHeight*tm.Height))
	//world.Resize()
	win := world.Window
	player := player.NewPlayer(animations, world)
	world.AddEntity("player", player)

	fps := 60
	ticks := 0
	interval := time.Duration(float64(1000) / float64(fps))
	ticker := time.NewTicker(time.Millisecond * interval)

	go func() {
		for {
			select {
			case <-ticker.C: //main game loop @normalized fps is here
				win.Clear(colornames.Black)
				tm.DrawLayers(win, []int{0, 1}) //draw base layer
				player.Update(ticks)
				tm.DrawLayers(win, []int{2, 3}) //draw everything else
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
