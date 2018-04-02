package world

import (
	"time"

	"github.com/scottwinkler/pixel-experiment/player"
	"github.com/scottwinkler/pixel-experiment/tilemap"
)

type World struct {
	Player  *player.Player
	Tilemap *tilemap.Tilemap
	Second  <-chan time.Time
}

func NewWorld() *World {
	world := World{
		Second: time.Tick(time.Second), //global timer
	}
	return &world
}

func (w *World) SetTilemap(tilemap *tilemap.Tilemap) {
	w.Tilemap = tilemap
}

func (w *World) SetPlayer(player *player.Player) {
	w.Player = player
}
