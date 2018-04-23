package entity

import (
	"math/rand"

	"github.com/scottwinkler/simple-rpg/enum"

	"github.com/scottwinkler/simple-rpg/world"
)

//SlimeController -- controller for slimes (and possibly other shit)
type SlimeController struct {
	entity *Entity
}

//NewSlimeController -- simple constructor
func NewSlimeController(entity *Entity) Controller {
	return &SlimeController{
		entity: entity,
	}
}

//HitCallback -- implementation method
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

//AttackCallback -- implementation method
func (c *SlimeController) AttackCallback(interface{}) {
	//do nothing
}

//Update -- implementation of controller interface
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
					success := e.Move(dir)
					if !success {
						switch dir {
						case enum.Direction.Down, enum.Direction.Up:
							if heads {
								e.Move(enum.Direction.Left)
							} else {
								e.Move(enum.Direction.Right)
							}
						case enum.Direction.Right, enum.Direction.Left:
							if heads {
								e.Move(enum.Direction.Up)
							} else {
								e.Move(enum.Direction.Down)
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
}
