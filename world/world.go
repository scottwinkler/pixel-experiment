package world

import (
	"math"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/simple-rpg/enum"
	"github.com/scottwinkler/simple-rpg/sfx"
	"github.com/scottwinkler/simple-rpg/tilemap"
	"golang.org/x/image/colornames"
)

//Camera -- the struct that holds details about the players camera
type Camera struct {
	Zoom   float64
	V      pixel.Vec
	Matrix pixel.Matrix
	Window *pixelgl.Window
}

//SetV -- sets camera position and updates underlying matrix
func (c *Camera) SetV(v pixel.Vec) {
	c.V = v
	c.Matrix = pixel.IM.Scaled(c.V, c.Zoom).Moved(c.Window.Bounds().Center().Sub(c.V))
}

//GUI -- a simple interface for graphical interfaces
type GUI interface {
	ID() string //probably should have a reference to the player too
	Update(int)
}

//Player -- a simple interface for players
type Player interface {
	GameObject() GameObject
	Camera() *Camera
}

//Callback -- Used for callbacks to gameobjects
type Callback func(interface{})

//GameObject -- interface that all entities (and other in game objects) must implement
type GameObject interface {
	Update(int)
	ID() string
	HandleHit(GameObject, Callback) bool
	Speed() float64
	Direction() int
	V() pixel.Vec
	R() float64
	Material() string
	Damage() float64
	Name() string
	World() *World
}

//World -- a minor diety in game terms
type World struct {
	Groups     map[string][]GameObject //should this be called "names"?
	GUIs       []GUI
	Tilemap    *tilemap.Tilemap
	Window     *pixelgl.Window
	sfxManager *sfx.Manager
}

//SFXManager -- getter function for sfxManager
func (w *World) SFXManager() *sfx.Manager {
	return w.sfxManager
}

//NewWorld -- constructor for world
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
	sfxManager := sfx.NewManager(effects, win)

	world := World{
		Groups:     make(map[string][]GameObject),
		Tilemap:    tilemap,
		Window:     win,
		sfxManager: sfxManager,
	}
	return &world
}

//UpdateGameObjects -- helper function to call update on each object that is part of the world
func (w *World) UpdateGameObjects(tick int) {
	for _, group := range w.Groups {
		for _, gameobject := range group {
			gameobject.Update(tick)
		}
	}
}

//UpdateGUIs -- helper function call update on each gui that is registered with the world. Typically a list
//of gui managers (only one if single player)
func (w *World) UpdateGUIs(tick int) {
	for _, gui := range w.GUIs {
		gui.Update(tick)
	}
}

//Start -- main game loop
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
				w.UpdateGUIs(normalizedTick)
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

//HitEvent -- could be better done with channels
func (w *World) HitEvent(source interface{}, callback Callback) {
	sourceGameObject := source.(GameObject)
	hitCount := 0
	for _, group := range w.Groups {
		for _, gameobject := range group {
			//don't notify the source of the hit
			if !strings.EqualFold(gameobject.ID(), sourceGameObject.ID()) {
				if gameobject.HandleHit(sourceGameObject, callback) { //leave it to the gameobjects to decide what to do
					hitCount++
				}
			}
		}
	}
	//no hits, attack has missed
	if hitCount == 0 {
		callback(nil)
	}
}

//RelativeDirection -- returns the direction of point a relative to point b (b is the center)
func RelativeDirection(posA pixel.Vec, posB pixel.Vec) int {
	posA = posA.Sub(posB)
	top := posA.Y >= posA.X    //above line y=x?
	right := posA.Y >= -posA.X //above line y=-x?
	var direction int
	if top && right {
		direction = enum.Direction.Up
	} else if top {
		direction = enum.Direction.Left
	} else if right {
		direction = enum.Direction.Right
	} else {
		direction = enum.Direction.Down
	}
	return direction
}

//Resize -- resizes window to tilemap dimensions (not really used...)
func (w *World) Resize() {
	maxY := float64(w.Tilemap.TileHeight * w.Tilemap.Height)
	maxX := float64(w.Tilemap.TileWidth * w.Tilemap.Width)
	bounds := pixel.R(0, 0, maxX, maxY)
	w.Window.SetBounds(bounds)
	w.Window.Update()
}

//GameObjectsByName -- gets all game objects of given name (i.e. group)
func (w *World) GameObjectsByName(name string) []GameObject {
	var gameObjects []GameObject
	for _, group := range w.Groups {
		for _, gameObject := range group {
			if strings.EqualFold(gameObject.Name(), name) {
				gameObjects = append(gameObjects, gameObject)
			}
		}
	}
	return gameObjects
}

//GameObjectByID -- gets a registered gameobject by id
func (w *World) GameObjectByID(id string) GameObject {
	var gameobject GameObject
	for _, group := range w.Groups {
		for _, obj := range group {
			if strings.EqualFold(obj.ID(), id) {
				gameobject = obj
				return gameobject
			}
		}
	}
	return gameobject
}

//AddGameObject -- the method that gameobjects call when they want to register themselves with the world
func (w *World) AddGameObject(group string, gameobject GameObject) {
	w.Groups[group] = append(w.Groups[group], gameobject)
}

//DeleteGameObject -- the method that gameobjects call when they want to deregister themselves with the world
func (w *World) DeleteGameObject(gameobject GameObject) {
	var newName string
	var newGroup []GameObject
	for name, group := range w.Groups {
		for i, obj := range group {
			if strings.EqualFold(obj.ID(), gameobject.ID()) {
				//delete this gameobject
				newName = name
				newGroup = append(group[:i], group[i+1:]...)
				break
			}
		}
	}
	w.Groups[newName] = newGroup
}

//AddGUI -- the method that guis call when they want to register themselves with the world
func (w *World) AddGUI(gui GUI) {
	w.GUIs = append(w.GUIs, gui)
}

//DeleteGUI -- the method that guis call when they want to deregister themselves with the world
func (w *World) DeleteGUI(gui GUI) {
	for i, obj := range w.GUIs {
		if strings.EqualFold(obj.ID(), gui.ID()) {
			w.GUIs = append(w.GUIs[:i], w.GUIs[i+1:]...)
			break
		}
	}
}

//CircleCollision -- returns true if the circles collide with each other.
func CircleCollision(v1 pixel.Vec, r1 float64, v2 pixel.Vec, r2 float64) bool {
	distanceSqr := math.Pow(v2.X-v1.X, 2) + math.Pow(v2.Y-v1.Y, 2)
	totalRadius := r1 + r2
	totalRadiusSqr := totalRadius * totalRadius
	return distanceSqr < totalRadiusSqr
}

//Collides -- returns true if a circle with the given point and radius collides with any other collidable entities circle
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
			if !strings.EqualFold(id, gameobject.ID()) && CircleCollision(v1, r1, v2, r2) {
				return true
			}
		}
	}
	return false
}
