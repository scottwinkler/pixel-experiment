package main

import (
	_ "image/png"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/simple-rpg/animation"
	"github.com/scottwinkler/simple-rpg/entity"
	"github.com/scottwinkler/simple-rpg/sfx"
	"github.com/scottwinkler/simple-rpg/sound"
	"github.com/scottwinkler/simple-rpg/tilemap"
	"github.com/scottwinkler/simple-rpg/utility"
	"github.com/scottwinkler/simple-rpg/world"
)

func run() {
	tm, _ := tilemap.ParseTiledJSON("_assets/tmx/world1.json")
	//particles
	particles_animations_mapping := utility.LoadJSON("./_configuration/entity/slime/effects.json")
	effects_spritesheet := utility.LoadSpritesheet("_assets/particles/blood_bath.png", pixel.R(0, 0, 256, 256), 0.5)
	effects := sfx.MappingToSFX(effects_spritesheet, particles_animations_mapping)

	//player
	player_sound_mapping := utility.LoadJSON("./_configuration/player/sounds.json")
	player_sounds := sound.MappingToSounds(player_sound_mapping)
	player_animations_mapping := utility.LoadJSON("./_configuration/player/animations.json")
	player_spritesheet := utility.LoadSpritesheet("_assets/spritesheets/chara_hero.png", pixel.R(0, 0, 48, 48), 2.0)
	player_animations := animation.MappingToAnimations(player_spritesheet, player_animations_mapping)

	//slime
	slime_spritesheet := utility.LoadSpritesheet("_assets/spritesheets/chara_slime.png", pixel.R(0, 0, 48, 48), 2.0)
	slime_animations_mapping := utility.LoadJSON("./_configuration/entity/slime/animations.json")
	slime_animations := animation.MappingToAnimations(slime_spritesheet, slime_animations_mapping)
	slime_sound_mapping := utility.LoadJSON("./_configuration/entity/slime/sounds.json")
	slime_sounds := sound.MappingToSounds(slime_sound_mapping)

	//slime spawner
	//todo: make sounds configuration include volumes
	spawner_spritesheet := utility.LoadSpritesheet("_assets/spritesheets/portalRings2.png", pixel.R(0, 0, 32, 32), 2.0)
	spawner_animations_mapping := utility.LoadJSON("./_configuration/entity/spawner/animations.json")
	spawner_animations := animation.MappingToAnimations(spawner_spritesheet, spawner_animations_mapping)
	spawner_sound_mapping := utility.LoadJSON("./_configuration/entity/spawner/sounds.json")
	spawner_sounds := sound.MappingToSounds(spawner_sound_mapping)

	w := world.NewWorld(pixel.R(0, 0, 400, 400), tm, effects)

	//player := player.NewPlayer(pixel.V(150, 200), 16, player_animations, player_sounds, w)
	playerConfig := &entity.EntityConfiguration{
		V: pixel.V(150, 200),
		W: w,
		Data: &entity.EntityData{
			R:          16,
			Animations: player_animations,
			Sounds:     player_sounds,
			Health:     12,
			Speed:      3,
			Name:       entity.ENTITY_PLAYER,
			Material:   world.MATERIAL_FLESH,
			Color:      colornames.Red,
		},
	}
	entity.NewEntity(playerConfig)

	spawnerConfig := &entity.EntityConfiguration{
		V: pixel.V(150, 300),
		W: w,
		Data: &entity.EntityData{
			R:          30,
			Animations: spawner_animations,
			Sounds:     spawner_sounds,
			Health:     12,
			Speed:      0,
			Name:       entity.ENTITY_SPAWNER,
			Material:   world.MATERIAL_WOOD,
			Color:      colornames.White,
		},
	}
	spawnData := &entity.EntityData{
		Animations: slime_animations,
		Sounds:     slime_sounds,
		Health:     9,
		Speed:      1,
		R:          16,
		Material:   world.MATERIAL_FLESH,
		Name:       entity.ENTITY_SLIME,
		Color:      colornames.Green,
	}
	entity.NewSpawner(spawnerConfig, spawnData)
	//	w.AddGameObject("player", player)
	w.Start(60.0, 15.0)

}

func main() {
	pixelgl.Run(run)
}
