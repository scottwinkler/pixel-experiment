package entity

import (
	"math/rand"

	"github.com/scottwinkler/simple-rpg/world"
)

type SlimeAi struct {
	entity *Entity
}

//simple constructor
func NewSlimeAi(entity *Entity) ai {
	return &SlimeAi{
		entity: entity,
	}
}

//implementation of ai interface
func (a *SlimeAi) Update(tick int) {
	var (
		e  = a.entity
		am = e.AnimationManager()
	)

	if am.Ready() {
		w := e.World()
		p := w.GameObjectById("player")
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
					if !e.Move(dir) {
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
	e.Draw(tick)
}
