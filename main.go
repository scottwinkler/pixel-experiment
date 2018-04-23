package main

import (
	_ "image/png"

	"github.com/scottwinkler/simple-rpg/enum"

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
	particlesAnimationsMapping := utility.LoadJSON("./_configuration/entity/slime/effects.json")
	effectsSpritesheet := utility.LoadSpritesheet("_assets/particles/blood_bath.png", pixel.R(0, 0, 256, 256), 0.5)
	effects := sfx.MappingToSFX(effectsSpritesheet, particlesAnimationsMapping)

	//player
	playerSoundMapping := utility.LoadJSON("./_configuration/player/sounds.json")
	playerSounds := sound.MappingToSounds(playerSoundMapping)
	playerAnimationsMapping := utility.LoadJSON("./_configuration/player/animations.json")
	playerSpritesheet := utility.LoadSpritesheet("_assets/spritesheets/chara_hero.png", pixel.R(0, 0, 48, 48), 2.0)
	playerAnimations := animation.MappingToAnimations(playerSpritesheet, playerAnimationsMapping)

	//slime
	slimeSpritesheet := utility.LoadSpritesheet("_assets/spritesheets/chara_slime.png", pixel.R(0, 0, 48, 48), 2.0)
	slimeAnimationsMapping := utility.LoadJSON("./_configuration/entity/slime/animations.json")
	slimeAnimations := animation.MappingToAnimations(slimeSpritesheet, slimeAnimationsMapping)
	slimeSoundMapping := utility.LoadJSON("./_configuration/entity/slime/sounds.json")
	slimeSounds := sound.MappingToSounds(slimeSoundMapping)

	//slime spawner
	//todo: make sounds configuration include volumes
	spawnerSpritesheet := utility.LoadSpritesheet("_assets/spritesheets/portalRings2.png", pixel.R(0, 0, 32, 32), 2.0)
	spawnerAnimationsMapping := utility.LoadJSON("./_configuration/entity/spawner/animations.json")
	spawnerAnimations := animation.MappingToAnimations(spawnerSpritesheet, spawnerAnimationsMapping)
	spawnerSoundMapping := utility.LoadJSON("./_configuration/entity/spawner/sounds.json")
	spawnerSounds := sound.MappingToSounds(spawnerSoundMapping)

	w := world.NewWorld(pixel.R(0, 0, 400, 400), tm, effects)

	//player := player.NewPlayer(pixel.V(150, 200), 16, player_animations, player_sounds, w)
	playerConfig := &entity.Configuration{
		V: pixel.V(150, 200),
		W: w,
		Data: &entity.Data{
			R:          16,
			Animations: playerAnimations,
			Sounds:     playerSounds,
			Health:     12,
			Speed:      4,
			Name:       enum.Entity.Player,
			Material:   enum.Material.Flesh,
			Color:      colornames.Red,
		},
	}
	entity.NewPlayer(playerConfig)

	spawnerConfig := &entity.Configuration{
		V: pixel.V(150, 300),
		W: w,
		Data: &entity.Data{
			R:          30,
			Animations: spawnerAnimations,
			Sounds:     spawnerSounds,
			Health:     12,
			Speed:      0,
			Name:       enum.Entity.Spawner,
			Material:   enum.Material.Wood,
			Color:      colornames.White,
		},
	}
	spawnData := &entity.Data{
		Animations: slimeAnimations,
		Sounds:     slimeSounds,
		Health:     9,
		Speed:      1,
		R:          16,
		Material:   enum.Material.Wood,
		Name:       enum.Entity.Slime,
		Color:      colornames.Green,
	}
	entity.NewSpawner(spawnerConfig, spawnData)
	w.Start(60.0, 15.0)
}

func main() {
	pixelgl.Run(run)
}
