package entity

import (
	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"github.com/scottwinkler/pixel-experiment/animation"
	"github.com/scottwinkler/pixel-experiment/world"
)

type Entity struct {
	id               string
	Sprite           *pixel.Sprite
	speed            float64
	v                pixel.Vec
	r                float64 //used for collider calculations
	Matrix           pixel.Matrix
	AnimationManager *animation.AnimationManager
	World            *world.World
	direction        int
}

func NewEntity(v pixel.Vec, r float64, animations []*animation.Animation, w *world.World) *Entity {
	animationManager := animation.NewAnimationManager(animations)
	animationManager.Select("Idle") //every entity should have an idle frame
	sprite := animationManager.Selected.Spritesheet.Sprites[animationManager.Selected.Frames[0]]
	matrix := pixel.Matrix(pixel.IM.Moved(v))
	//radius := float64(world.Tilemap.TileHeight / 3) //not a great solution for rectangles
	id := xid.New().String()
	entity := &Entity{
		id:               id,
		Sprite:           sprite,
		speed:            3, //default
		v:                v,
		r:                r,
		Matrix:           matrix,
		AnimationManager: animationManager,
		World:            w,
		direction:        world.DOWN,
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
	if !e.World.Collides(e.Id(), nextPos, e.r) {
		e.v = nextPos
	}
	//update matrix and collision circle
	matrix := pixel.IM.Moved(e.v)
	e.Matrix = matrix
}

func (e *Entity) HandleHit(s world.GameObject) {
	//am i near enoguh to be affected?
	//draw a slightly bigger circle than the collision circle
	//so that the hit box is reasonable

	//fmt.Printf("handling hit from: %s", s.Id())
	hitFactor := 1.2
	if world.CircleCollision(e.v, e.r*hitFactor, s.V(), s.R()+s.Speed()) {
		//fmt.Println("collided")
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
			//fmt.Println("direction correct")
			switch e.direction {
			case world.LEFT:
				e.AnimationManager.Select("HitLeft")
			case world.RIGHT:
				e.AnimationManager.Select("HitRight")
			case world.DOWN:
				e.AnimationManager.Select("HitDown")
			case world.UP:
				e.AnimationManager.Select("HitUp")
			}
		}
	}
}

func (e *Entity) Update(tick int) {
	//win := e.World.Window
	if e.AnimationManager.Selected.Skippable() || e.AnimationManager.Selected.Done() { //only listen to new events if the current animation is skippable or done playing
		e.AnimationManager.Select("Idle")
	}
	if tick == 0 {
		e.Sprite = e.AnimationManager.Selected.Next()
	}
	e.Draw()
}
func (e *Entity) Draw() {
	animation := e.AnimationManager.Selected
	//chained methods so that we first scale by spritesheet size, then by reflection, then by position
	matrix := animation.Matrix.Chained(e.Matrix)
	e.Sprite.Draw(e.World.Window, matrix)
}
