package world

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/pixel-experiment/tilemap"
	"golang.org/x/image/colornames"
)

type GameObject interface {
	Update(int)
	Collider() (pixel.Vec, float64)
	Id() string
}

type World struct {
	Groups  map[string][]GameObject
	Tilemap *tilemap.Tilemap
	Window  *pixelgl.Window
}

func NewWorld(bounds pixel.Rect, tilemap *tilemap.Tilemap) *World {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: bounds,
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	world := World{
		Groups:  make(map[string][]GameObject),
		Tilemap: tilemap,
		Window:  win,
	}
	return &world
}
func (w *World) UpdateGameObjects(tick int) {
	for _, group := range w.Groups {
		for _, gameobject := range group {
			gameobject.Update(tick)
		}
	}
}

func (w *World) Start(fps float64) {
	tick := 0
	interval := time.Duration(float64(1000) / float64(fps))
	ticker := time.NewTicker(time.Millisecond * interval)
	win := w.Window
	tm := w.Tilemap
	maxTick := 1
	if fps >= 10 {
		maxTick = int(fps) / 10
	}
	go func() {
		for {
			select {
			case <-ticker.C: //main game loop @normalized fps is here
				win.Clear(colornames.Black)
				tm.DrawLayers(win, []string{"Ground", "Rocks"}) //draw base layers
				w.UpdateGameObjects(tick)
				tm.DrawLayers(win, []string{"Treetops", "Collision"}) //draw top layers
				win.Update()
				tick++
				//toggle for changing animation frame
				if tick > maxTick {
					tick = 0
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

//resizes window to tilemap dimensions
func (w *World) Resize() {
	maxY := float64(w.Tilemap.TileHeight * w.Tilemap.Height)
	maxX := float64(w.Tilemap.TileWidth * w.Tilemap.Width)
	bounds := pixel.R(0, 0, maxX, maxY)
	w.Window.SetBounds(bounds)
	w.Window.Update()
}

func (w *World) AddGameObject(group string, gameobject GameObject) {
	w.Groups[group] = append(w.Groups[group], gameobject)
}

//CircleCollision returns true if the circles collide with each other.
func CircleCollision(v1 pixel.Vec, r1 float64, v2 pixel.Vec, r2 float64) bool {
	differenceV := pixel.V(v1.X, v2.Y).Sub(v1)
	totalRadius := r1 + r2
	totalRadiusSqr := totalRadius * totalRadius
	distanceSqr := differenceV.Dot(differenceV)
	return distanceSqr > totalRadiusSqr
}

//returns true if a circle with the given point and radius collides with any other collidable entities circle
//or with any of the predefined collision tiles
func (w *World) Collides(id string, v1 pixel.Vec, r1 float64) bool {
	//check collision tile
	if !w.Tilemap.Bounds().Contains(v1) {
		return true //out of bounds!
	}
	tile := w.Tilemap.GetTileAtPosition(v1, "Collision")
	if tile != nil && tile.IsCollidable {
		return true
	}

	//check if it collides with an existing game object which is not this one
	for _, gameobjects := range w.Groups {
		for _, gameobject := range gameobjects {
			v2, r2 := gameobject.Collider()
			if CircleCollision(v1, r1, v2, r2) && id != gameobject.Id() {
				return true
			}
		}
	}
	return false
}
