package entity

import (
	"math/rand"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/simple-rpg/world"
)

type PlayerController struct {
	entity *Entity
}

//simple constructor
func NewPlayerController(entity *Entity) controller {
	return &PlayerController{
		entity: entity,
	}
}

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
		case world.MATERIAL_FLESH:
			sm.Play("humanattacking" + num)
		case world.MATERIAL_METAL:
			sm.Play("humanattacking"+num, "swordhitmetal")
		case world.MATERIAL_WOOD:
			sm.Play("humanattacking"+num, "swordhitwood")
		}
	} else {
		sm.Play("humanattacking"+num, "swordswing")
	}
}

//implementation of controller interface
func (c *PlayerController) Update(tick int) {
	var (
		e       = c.entity
		camPos  = e.v //position camera centered on player
		camZoom = 1.0
		am      = e.animationManager
	)
	win := e.world.Window
	cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
	win.SetMatrix(cam)

	if am.Ready() { //only listen to new events if the animation manager is ready to accept new input
		if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
			e.Move(world.LEFT)
		} else if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
			e.Move(world.RIGHT)
		} else if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			e.Move(world.DOWN)
		} else if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			e.Move(world.UP)
		} else if win.Pressed(pixelgl.MouseButtonLeft) {
			//determine what quadrant relative to the player the mouse click happened
			mouse := cam.Unproject(win.MousePosition())
			dir := world.RelativeDirection(mouse, e.v)
			e.Attack(dir)
		} else {
			am.Select("Idle")
		}
	}
	e.Draw(tick)
}
