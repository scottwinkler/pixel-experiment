package entity

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"github.com/scottwinkler/simple-rpg/animation"
	"github.com/scottwinkler/simple-rpg/world"
)

type Entity struct {
	id               string
	sprite           *pixel.Sprite
	speed            float64
	v                pixel.Vec
	r                float64 //used for collider calculations
	matrix           pixel.Matrix
	animationManager *animation.AnimationManager
	world            *world.World
	direction        int
	health           float64
}

func NewEntity(v pixel.Vec, r float64, animations []*animation.Animation, w *world.World) *Entity {
	animationManager := animation.NewAnimationManager(animations)
	animationManager.Select("Idle") //every entity should have an idle frame
	sprite := animationManager.Selected.Next(0)
	matrix := pixel.Matrix(pixel.IM.Moved(v))
	id := xid.New().String()
	entity := &Entity{
		id:               id,
		sprite:           sprite,
		speed:            3, //default
		v:                v,
		r:                r,
		matrix:           matrix,
		animationManager: animationManager,
		world:            w,
		direction:        world.DOWN,
		health:           15.0,
	}
	return entity
}

func (e *Entity) Id() string {
	return e.id
}

func (e *Entity) V() pixel.Vec {
	return e.v
}

func (e *Entity) R() float64 {
	return e.r
}

func (e *Entity) Speed() float64 {
	return e.speed
}

func (e *Entity) Direction() int {
	return e.direction
}

func (e *Entity) Kill() {
	fmt.Println("killing...")
	e.world.SFXManager().MakeEffect("BloodExplosion", e.v)
	e.world.DeleteGameObject(e)
}

func (e *Entity) Material() string {
	return world.MATERIAL_FLESH
}

func (e *Entity) Move(direction int) {
	nextPos := pixel.V(e.v.X, e.v.Y)
	switch direction {
	case world.LEFT:
		nextPos.X -= e.speed
	case world.RIGHT:
		nextPos.X += e.speed
	case world.DOWN:
		nextPos.Y -= e.speed
	case world.UP:
		nextPos.Y += e.speed
	}
	e.direction = direction
	if !e.world.Collides(e.Id(), nextPos, e.r) {
		e.v = nextPos
	}
	//update matrix and collision circle
	matrix := pixel.IM.Moved(e.v)
	e.matrix = matrix
}

func (e *Entity) HandleHit(s world.GameObject, cb world.Fn_Callback) bool {
	//am i near enoguh to be affected?
	//draw a slightly bigger circle than the collision circle
	//so that the hit box is reasonable
	hitFactor := 1.2
	if world.CircleCollision(e.v, e.r*hitFactor, s.V(), s.R()+s.Speed()) {
		//where am i relative to the source?
		relativePos := e.v.Sub(s.V())
		top := relativePos.Y >= relativePos.X    //above line y=x?
		right := relativePos.Y >= -relativePos.X //above line y=-x?

		var relativeDir int
		if top && right {
			relativeDir = world.UP
		} else if top {
			relativeDir = world.LEFT
		} else if right {
			relativeDir = world.RIGHT
		} else {
			relativeDir = world.DOWN
		}
		//is the source facing the right direction?
		if relativeDir == s.Direction() {
			switch e.direction {
			case world.LEFT:
				e.animationManager.Select("HitLeft")
			case world.RIGHT:
				e.animationManager.Select("HitRight")
			case world.DOWN:
				e.animationManager.Select("HitDown")
			case world.UP:
				e.animationManager.Select("HitUp")
			}
			e.health -= 3
			fmt.Println(e.health)
			if e.health <= 0 {
				e.Kill()
			}
			cb(e)
			return true
		}
	}

	return false
}

func (e *Entity) Update(tick int) {
	//win := e.World.Window
	if e.animationManager.Selected.Skippable() || e.animationManager.Selected.Done() { //only listen to new events if the current animation is skippable or done playing
		e.animationManager.Select("Idle")
	}
	e.sprite = e.animationManager.Selected.Next(tick)
	e.Draw()
}
func (e *Entity) Draw() {
	animation := e.animationManager.Selected
	//chained methods so that we first scale by spritesheet size, then by reflection, then by position
	matrix := animation.Matrix().Chained(e.matrix)
	e.sprite.Draw(e.world.Window, matrix)
}
