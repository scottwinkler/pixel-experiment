package world

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/pixel-experiment/tilemap"
)

type World struct {
	Entities map[string]interface{}
	Tilemap  *tilemap.Tilemap
	Window   *pixelgl.Window
}

func NewWorld(width float64, height float64) *World {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	world := World{
		Entities: make(map[string]interface{}),
		Window:   win,
	}
	return &world
}

//resizes window to tilemap dimensions
func (w *World) Resize() {
	maxY := float64(w.Tilemap.TileHeight * w.Tilemap.Height)
	maxX := float64(w.Tilemap.TileWidth * w.Tilemap.Width)
	bounds := pixel.R(0, 0, maxX, maxY)
	w.Window.SetBounds(bounds)
	w.Window.Update()
}

func (w *World) SetTilemap(tilemap *tilemap.Tilemap) {
	w.Tilemap = tilemap
}

func (w *World) AddEntity(name string, entity interface{}) {
	w.Entities[name] = entity
}
