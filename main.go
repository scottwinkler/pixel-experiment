package main

import (
	_ "image/png"
	"time"

	"github.com/faiface/pixel"
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

	player_spritesheet := spritesheet.LoadSpritesheet("assets/spritesheets/knight_iso_char.png", pixel.R(0, 0, 84, 84), 1.0)
	//var animations []*animation.Animation
	//should make a method that makes this proceess easier. for example
	//animations.AddAnimation(name string,frames[]int) []*animation.Animation0
	//all of this garbage should be put in either a config file or a preloader function
	player_mapping := map[string][]int{
		"Idle":        []int{0, 1, 2, 3},
		"RunDown":     []int{4, 5, 6, 7, 8},
		"RunUp":       []int{9, 10, 11, 12, 13},
		"RunRight":    []int{14, 15, 16, 17, 18, 19},
		"RunLeft":     []int{-14, -15, -16, -17, -18, -19},
		"AttackDown":  []int{26, 27, 28},
		"AttackUp":    []int{29, 30, 31},
		"AttackRight": []int{32, 33, 34},
		"AttackLeft":  []int{35, 36, 37},
	}
	player_animations := animation.AnimationsFromSpritesheet(player_spritesheet, player_mapping)

	//load spritesheet should accept a rectangle, maybe? and a resize factor?

	/*slime_spritesheet := spritesheet.LoadSpritesheet("assets/spritesheets/chara_slime.png", pixel.R(0, 0, 16, 16), 2.0)
	//if animations receives a negative number then it should know to flip it
	slime_mapping := map[string][]int{
		"Idle":        []int{0, 1, 2},
		"Wake":        []int{4, 5, 6},
		"MoveDown":    []int{8, 9, 10, 11},
		"MoveRight":   []int{12, 13, 14, 15},
		"MoveLeft":    []int{-12, -13, -14, -15},
		"MoveUp":      []int{16, 17, 18, 19},
		"AttackDown":  []int{20, 21, 22, 23},
		"AttackRight": []int{24, 25, 26, 27},
		"AttackLeft":  []int{-24, -25, -26, -27},
		"AttackUp":    []int{28, 29, 30, 31},
		"HitDown":     []int{32, 33, 34},
		"HitRight":    []int{36, 37, 38},
		"HitLeft":     []int{-36, -37, -38},
		"HitUp":       []int{40, 41, 42},
	}*/
	//slime_animations := animation.AnimationsFromSpritesheet(slime_spritesheet, slime_mapping)

	world := world.NewWorld(400, 400)
	world.SetTilemap(tm)
	//v := pixel.V(float64(tm.TileWidth*tm.Width), float64(tm.TileHeight*tm.Height))
	//world.Resize()
	win := world.Window

	//the world is not ready for you yet, my friend
	//entity := entity.NewEntity(slime_animations, world)
	player := player.NewPlayer(player_animations, world)
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
				tm.DrawLayers(win, []string{"Ground", "Rocks"}) //draw base layer
				player.Update(ticks)
				//	entity.Update(ticks)
				tm.DrawLayers(win, []string{"Treetops", "Collision"}) //draw everything else
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
