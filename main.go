package main

import (
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/scottwinkler/pixel-experiment/animation"
	"github.com/scottwinkler/pixel-experiment/utility"
	"github.com/scottwinkler/pixel-experiment/world"

	"github.com/scottwinkler/pixel-experiment/player"

	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/pixel-experiment/tilemap"
)

func run() {

	tm, _ := tilemap.ParseTiledJSON("_assets/tmx/world1.json")

	player_spritesheet := utility.LoadSpritesheet("_assets/spritesheets/knight_iso_char.png", pixel.R(0, 0, 84, 84), 0.7)
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
	/*
		slime_spritesheet := spritesheet.LoadSpritesheet("_assets/spritesheets/chara_slime.png", pixel.R(0, 0, 16, 16), 2.0)
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
		}
		slime_animations := animation.AnimationsFromSpritesheet(slime_spritesheet, slime_mapping)*/

	world := world.NewWorld(pixel.R(0, 0, 400, 400), tm)
	//world.SetTilemap(tm)
	//v := pixel.V(float64(tm.TileWidth*tm.Width), float64(tm.TileHeight*tm.Height))
	//world.Resize()
	//win := world.Window

	//the world is not ready for you yet, my friend
	//entity := entity.NewEntity(slime_animations, world)
	player := player.NewPlayer(pixel.V(150, 200), player_animations, world)

	world.AddGameObject("player", player)
	world.Start(60)

}

func main() {
	pixelgl.Run(run)
}
