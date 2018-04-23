package entity

import (
	"github.com/scottwinkler/simple-rpg/gui"
	"github.com/scottwinkler/simple-rpg/world"
)

//Player -- the container struct of players
type Player struct {
	entity     *Entity //the in game representation of the player
	camera     *world.Camera
	guiManager *gui.Manager
}

//NewPlayer constructor for player objects
func NewPlayer(config *Configuration) *Player {
	win := config.W.Window
	camera := &world.Camera{
		Zoom:   1.0,
		Window: win,
	}
	camera.SetV(config.V)
	player := &Player{
		camera: camera,
	}
	entity := NewEntity(config, player)
	player.entity = entity
	guiManager := gui.NewManager(player)
	player.guiManager = guiManager

	return player
}

//Camera -- getter method for the players camera. Implementation of world.Player
func (p *Player) Camera() *world.Camera {
	return p.camera
}

//SetCamera -- setter method for the players camera
func (p *Player) SetCamera(camera *world.Camera) {
	p.camera = camera
	p.camera.Window.SetMatrix(p.camera.Matrix)
}

//GameObject -- getter method for underlying entity. Impelemenationf of world.Player
func (p *Player) GameObject() world.GameObject {
	return p.entity
}
