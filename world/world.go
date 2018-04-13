package world

import (
	"math"
	"strings"
	"time"

	"github.com/scottwinkler/simple-rpg/sfx"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/simple-rpg/tilemap"
	"golang.org/x/image/colornames"
)

//Directions
const (
	LEFT           = 0
	RIGHT          = 1
	DOWN           = 2
	UP             = 3
	MATERIAL_WOOD  = "wood"
	MATERIAL_FLESH = "flesh"
	MATERIAL_METAL = "metal"
)

//Used for callbacks to gameobjects
type Fn_Callback func(interface{})

type GameObject interface {
	Update(int)
	Id() string
	HandleHit(GameObject, Fn_Callback) bool
	Speed() float64
	Direction() int
	V() pixel.Vec
	R() float64
	Material() string
}

type World struct {
	Groups     map[string][]GameObject
	Tilemap    *tilemap.Tilemap
	Window     *pixelgl.Window
	sfxManager *sfx.SFXManager
}

func (w *World) SFXManager() *sfx.SFXManager {
	return w.sfxManager
}

func NewWorld(bounds pixel.Rect, tilemap *tilemap.Tilemap, effects []*sfx.SFX) *World {
	cfg := pixelgl.WindowConfig{
		Title:  "Simple RPG!",
		Bounds: bounds,
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	sfxManager := sfx.NewSFXManager(effects, win)

	world := World{
		Groups:     make(map[string][]GameObject),
		Tilemap:    tilemap,
		Window:     win,
		sfxManager: sfxManager,
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

//could be better done with channels
func (w *World) HitEvent(source interface{}, callback Fn_Callback) {
	sourceGameObject := source.(GameObject)
	hitCount := 0
	for _, group := range w.Groups {
		for _, gameobject := range group {
			//don't notify the source of the hit
			if !strings.EqualFold(gameobject.Id(), sourceGameObject.Id()) {
				if gameobject.HandleHit(sourceGameObject, callback) { //leave it to the gameobjects to decide what to do
					hitCount++
				}
			}
		}
	}
	//fmt.Printf("hitcount: %d", hitCount)
	if hitCount == 0 {
		//fmt.Println("a miss")
		callback(nil)
	}
}

//returns the direction of point a relative to point b (b is the center)
func RelativeDirection(posA pixel.Vec, posB pixel.Vec) int {
	posA = posA.Sub(posB)
	top := posA.Y >= posA.X    //above line y=x?
	right := posA.Y >= -posA.X //above line y=-x?
	var direction int
	if top && right {
		direction = UP
	} else if top {
		direction = LEFT
	} else if right {
		direction = RIGHT
	} else {
		direction = DOWN
	}
	return direction
}

func (w *World) Start(fps float64, animationSpeed float64) {
	tick := 1
	interval := time.Duration(float64(1000) / float64(fps))
	ticker := time.NewTicker(time.Millisecond * interval)
	win := w.Window
	tm := w.Tilemap
	go func() {
		for {
			select {
			case <-ticker.C: //main game loop @normalized fps is here
				win.Clear(colornames.Black)
				tm.DrawLayers(win, []string{"Ground", "Rocks"}) //draw base layers
				//the calculations for animations are easier if they can assume a normalized 60 ticks per second
				normalizedTick := int((60 / int(fps)) * tick)
				w.UpdateGameObjects(normalizedTick)
				w.sfxManager.Update(normalizedTick)
				tm.DrawLayers(win, []string{"Treetops", "Collision"}) //draw top layers
				win.Update()
				tick++
				if tick > int(fps) {
					tick = 1
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

func (w *World) DeleteGameObject(gameobject GameObject) {
	var new_name string
	var new_group []GameObject
	for name, group := range w.Groups {
		for i, obj := range group {
			if strings.EqualFold(obj.Id(), gameobject.Id()) {
				//delete this gameobject
				new_name = name
				new_group = append(group[:i], group[i+1:]...)
				break
			}
		}
	}
	w.Groups[new_name] = new_group
}

//CircleCollision returns true if the circles collide with each other.
func CircleCollision(v1 pixel.Vec, r1 float64, v2 pixel.Vec, r2 float64) bool {
	distanceSqr := math.Pow(v2.X-v1.X, 2) + math.Pow(v2.Y-v1.Y, 2)
	totalRadius := r1 + r2
	totalRadiusSqr := totalRadius * totalRadius
	return distanceSqr < totalRadiusSqr
}

//returns true if a circle with the given point and radius collides with any other collidable entities circle
//or with any of the predefined collision tiles
func (w *World) Collides(id string, v1 pixel.Vec, r1 float64) bool {
	//check collision tile
	if !w.Tilemap.Bounds().Contains(v1) {
		return true //out of bounds!
	}
	tile := w.Tilemap.GetTileAtPosition(v1, "Collision")
	if tile != nil && tile.Collidable() {
		return true
	}

	//check if it collides with an existing game object which is not this one
	for _, gameobjects := range w.Groups {
		for _, gameobject := range gameobjects {
			v2 := gameobject.V()
			r2 := gameobject.R()
			if !strings.EqualFold(id, gameobject.Id()) && CircleCollision(v1, r1, v2, r2) {
				return true
			}
		}
	}
	return false
}
