package main

import (
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/simple-rpg/animation"
	"github.com/scottwinkler/simple-rpg/entity"
	"github.com/scottwinkler/simple-rpg/player"
	"github.com/scottwinkler/simple-rpg/sfx"
	"github.com/scottwinkler/simple-rpg/sound"
	"github.com/scottwinkler/simple-rpg/tilemap"
	"github.com/scottwinkler/simple-rpg/utility"
	"github.com/scottwinkler/simple-rpg/world"
)

func run() {
	player_sound_mapping := utility.LoadJSON("./_configuration/player/sounds.json")
	player_sounds := sound.MappingToSounds(player_sound_mapping)

	tm, _ := tilemap.ParseTiledJSON("_assets/tmx/world1.json")
	particles_animations_mapping := utility.LoadJSON("./_configuration/entity/effects.json")
	effects_spritesheet := utility.LoadSpritesheet("_assets/particles/blood_bath.png", pixel.R(0, 0, 256, 256), 0.5)
	effects := sfx.MappingToSFX(effects_spritesheet, particles_animations_mapping)
	player_animations_mapping := utility.LoadJSON("./_configuration/player/animations.json")
	slime_animations_mapping := utility.LoadJSON("./_configuration/entity/animations.json")
	player_spritesheet := utility.LoadSpritesheet("_assets/spritesheets/chara_hero.png", pixel.R(0, 0, 48, 48), 2.0)
	player_animations := animation.MappingToAnimations(player_spritesheet, player_animations_mapping)
	slime_spritesheet := utility.LoadSpritesheet("_assets/spritesheets/chara_slime.png", pixel.R(0, 0, 48, 48), 2.0)
	slime_animations := animation.MappingToAnimations(slime_spritesheet, slime_animations_mapping)

	world := world.NewWorld(pixel.R(0, 0, 400, 400), tm, effects)
	slime := entity.NewEntity(pixel.V(150, 300), 16, slime_animations, world)
	player := player.NewPlayer(pixel.V(150, 200), 16, player_animations, player_sounds, world)

	world.AddGameObject("player", player)
	world.AddGameObject("monster", slime)
	world.Start(60.0, 15.0)

}

func main() {
	pixelgl.Run(run)
}
