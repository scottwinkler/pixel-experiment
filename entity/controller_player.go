package entity

import (
	"math/rand"
	"strconv"

	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/simple-rpg/enum"
	"github.com/scottwinkler/simple-rpg/world"
)

//PlayerController -- controller for player
type PlayerController struct {
	entity *Entity
}

//NewPlayerController -- constructor
func NewPlayerController(entity *Entity) Controller {
	return &PlayerController{
		entity: entity,
	}
}

//HitCallback -- impelmentation method
func (c *PlayerController) HitCallback(source interface{}) bool {
	var (
		s  = source.(world.GameObject)
		e  = c.entity
		sm = e.soundManager
	)
	//is this a killing blow?
	if e.health-s.Damage() <= 0 {
		sm.Play("death0")
	} else {
		sm.Play("hit0")
	}
	return true
}

//AttackCallback -- implementation method
func (c *PlayerController) AttackCallback(source interface{}) {
	var (
		e  = c.entity
		sm = e.soundManager
	)
	num := strconv.Itoa(rand.Intn(2) + 1) //for playing a random attack sound
	if source != nil {
		var (
			s        = source.(world.GameObject)
			material = s.Material()
		)
		//fmt.Printf("[DEBUG] -- material: %s", material)
		switch material {
		case enum.Material.Flesh:
			sm.Play("humanattacking" + num)
		case enum.Material.Metal:
			sm.Play("humanattacking"+num, "swordhitmetal")
		case enum.Material.Wood:
			sm.Play("humanattacking"+num, "swordhitwood")
		}
	} else {
		sm.Play("humanattacking"+num, "swordswing")
	}
}

//Update -- implementation of controller interface
func (c *PlayerController) Update(tick int) {
	var (
		e   = c.entity
		p   = e.parent.(*Player)
		am  = e.animationManager
		win = e.world.Window
	)

	if am.Ready() { //only listen to new events if the animation manager is ready to accept new input
		if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
			e.Move(enum.Direction.Left)
		} else if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
			e.Move(enum.Direction.Right)
		} else if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			e.Move(enum.Direction.Down)
		} else if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			e.Move(enum.Direction.Up)
		} else if win.Pressed(pixelgl.MouseButtonLeft) {
			//determine what quadrant relative to the player the mouse click happened
			mouse := p.Camera().Matrix.Unproject(win.MousePosition())
			dir := world.RelativeDirection(mouse, e.v)
			e.Attack(dir)
		} else {
			am.Select("Idle")
		}
	}
	//other mouse events
	if win.Pressed(pixelgl.KeyI) {
		p.guiManager.SetVisible(true)
	}

	//update camera
	cam := p.Camera()
	cam.SetV(e.v)
	p.SetCamera(cam)
}
