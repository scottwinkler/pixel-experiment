package entity

import (
	"time"

	"github.com/faiface/pixel"
)

type Spawner struct {
	entity        *Entity     //the spawner itself
	spawnLocation pixel.Rect  //where the creatures are spawned
	spawnData     *EntityData //the kind of creature it spawns
}

//constructor for spawner
func NewSpawner(config *EntityConfiguration, spawnData *EntityData) *Spawner {
	entity := NewEntity(config)
	spawner := &Spawner{
		entity:        entity,
		spawnData:     spawnData,
		spawnLocation: pixel.R(entity.v.X-30, entity.v.Y-80, entity.v.X+30, entity.v.Y-50), //where the creatures are spawned
	}

	//create monsters for as long as the spawner is alive
	go func() {
		for {
			if entity.IsDead() {
				break
			}
			spawner.Spawn()
			time.Sleep(5000 * time.Millisecond)
		}
	}()
	return spawner
}

//helper method for determing a point where to spawn. returns a point, and a boolean for success condition
func (s *Spawner) getSpawnPosition() (pixel.Vec, bool) {
	minX := s.spawnLocation.Min.X
	minY := s.spawnLocation.Min.Y
	maxX := s.spawnLocation.Max.X
	maxY := s.spawnLocation.Max.Y
	spawnDiameter := s.spawnData.R * 2

	//not going to try to do optimal packing because who the fuck cares
	var spawnPos pixel.Vec
	for x := minX; x <= maxX; x += spawnDiameter {
		for y := maxY; y >= minY; y -= spawnDiameter {
			spawnPos = pixel.V(float64(x), float64(y))
			//ensure there won't be a collision by placing an entity at the given point
			if !s.entity.world.Collides("DEADBEEF", spawnPos, s.spawnData.R) {
				return spawnPos, true
			}
		}
	}
	return spawnPos, false
}

func (s *Spawner) Spawn() {
	spawnPos, success := s.getSpawnPosition()
	if !success {
		return
	}
	//create a new entity at the spawner location
	config := &EntityConfiguration{
		V:    spawnPos, //spawn nearby
		W:    s.entity.world,
		Data: s.spawnData,
	}
	NewEntity(config)
}
