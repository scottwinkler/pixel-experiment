package entity

import "github.com/scottwinkler/simple-rpg/world"

type NullController struct {
	entity *Entity
}

//simple constructor
func NewNullController(entity *Entity) controller {
	return &NullController{
		entity: entity,
	}
}

//implementation of controller interface
func (c *NullController) HitCallback(source interface{}) bool {
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

//implementation of controller interface
func (c *NullController) AttackCallback(s interface{}) {
	//do nothing
}

//implementation of controller interface
func (c *NullController) Update(tick int) {
	var (
		e  = c.entity
		am = e.AnimationManager()
	)
	if am.Ready() {
		am.Select("Idle")
	}
	e.Draw(tick)
}
