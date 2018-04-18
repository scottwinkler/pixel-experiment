package entity

import (
	"math/rand"

	"github.com/scottwinkler/simple-rpg/world"
)

type SlimeController struct {
	entity *Entity
}

//simple constructor
func NewSlimeController(entity *Entity) controller {
	return &SlimeController{
		entity: entity,
	}
}

func (c *SlimeController) HitCallback(source interface{}) bool {
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

func (c *SlimeController) AttackCallback(interface{}) {
	//do nothing
}

//implementation of controller interface
func (c *SlimeController) Update(tick int) {
	var (
		e  = c.entity
		am = e.AnimationManager()
	)

	if am.Ready() {
		w := e.World()
		players := w.GameObjectsByName("player") //always get first player... wont work for multiplayer
		//condition on death of player
		if len(players) == 0 {
			am.Select("Idle")
			return
		}
		p := players[0] //select first player
		distanceToPlayer := p.V().To(e.V()).Len()
		noticeRange := 5 * e.R() //magic numbers, beware
		aggroRange := 10 * e.R()
		lastState := am.Selected().Name()
		hitDistance := p.R() + e.R() + e.Speed()

		switch lastState {
		case "Idle":
			if distanceToPlayer <= noticeRange {
				am.Select("Wake")
			} else {
				am.Select("Idle")
			}
		case "AttackUp", "AttackDown", "AttackLeft", "AttackRight":
			//take a break so the attack is not so relentless
			am.Select("Rest")
		default:
			if distanceToPlayer <= aggroRange {
				dir := world.RelativeDirection(p.V(), e.V())
				if distanceToPlayer <= hitDistance {
					e.Attack(dir)
				} else {
					//add some randomness in the choice of direction to choose with a coin toss
					heads := rand.Intn(2) == 0
					//try moving in the obvious direction first
					moveSuccess := e.Move(dir)
					if !moveSuccess {
						switch dir {
						case world.DOWN, world.UP:
							if heads {
								e.Move(world.LEFT)
							} else {
								e.Move(world.RIGHT)
							}
						case world.RIGHT, world.LEFT:
							if heads {
								e.Move(world.UP)
							} else {
								e.Move(world.DOWN)
							}
						}
					}
				}
			} else {
				//out of range
				am.Select("Idle")
			}
		}
	}
	//e.Draw(tick)
}
