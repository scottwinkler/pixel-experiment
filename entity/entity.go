package entity

import (
	"fmt"
	"time"

	"github.com/scottwinkler/simple-rpg/sound"

	"github.com/scottwinkler/simple-rpg/sfx"

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
	soundManager     *sound.SoundManager
	world            *world.World
	direction        int
	health           float64
}

func NewEntity(v pixel.Vec, r float64, animations []*animation.Animation, sounds []*sound.Sound, w *world.World) *Entity {
	animationManager := animation.NewAnimationManager(animations)
	animationManager.Select("Idle") //every entity should have an idle frame
	soundManager := sound.NewSoundManager(sounds)
	id := xid.New().String()
	entity := &Entity{
		id:               id,
		speed:            3, //default
		v:                v,
		r:                r,
		animationManager: animationManager,
		soundManager:     soundManager,
		world:            w,
		direction:        world.DOWN, //default
		health:           15.0,       //default -- todo: should be changed by parameter
	}
	return entity
}

//getter method for id
func (e *Entity) Id() string {
	return e.id
}

//getter method for v
func (e *Entity) V() pixel.Vec {
	return e.v
}

//getter method for r
func (e *Entity) R() float64 {
	return e.r
}

//getter method for speed
func (e *Entity) Speed() float64 {
	return e.speed
}

//getter method for direction
func (e *Entity) Direction() int {
	return e.direction
}

//a calculated effect to play once when the entity dies
func (e *Entity) PlayDeathEffect() {
	var sfxFrames []sfx.SFXFrame
	animation := e.animationManager.Selected()
	_, animationFrame := animation.Current()
	//matrix := animationFrame.Matrix
	mask := animationFrame.Mask
	frame := animationFrame.Frame
	//8 frames feels right. this is approximately how long the death sound takes at a frameRate of 10
	framesCount := 8
	for i := framesCount; i >= 0; i-- {
		scaleFactor := float64(i) / float64(framesCount)
		fmt.Printf("scalefactor: %f", scaleFactor)
		//this exponential function makes it more interesting than a simple linear interpolation
		//y := math.Exp(1 - 1/math.Pow(scaleFactor, 2))
		//scaleMatrix := pixel.IM.
		matrix := animationFrame.Matrix //.ScaledXY(e.v, pixel.V(1, scaleFactor))
		matrix = matrix.ScaledXY(pixel.ZV, pixel.V(1, scaleFactor))
		//add a green color mask to the sprite, and make it incrementally smaller
		mask = mask.Mul(pixel.RGB(0.2, 1, 0.5))
		//todo: come up with a method that converts hex colors to rgb masks because this is really hard to get right
		sfxFrame := sfx.NewSFXFrame(frame, matrix, mask)
		sfxFrames = append(sfxFrames, sfxFrame)
	}
	frameRate := animation.FrameRate()
	sfx := sfx.NewSFX(animation.Spritesheet(), "EntityDeath", sfxFrames, false, frameRate)
	//the slime looks better if its translated slightly down... not critical though
	//and should be tested with other sprites
	//pos := pixel.V(e.v.X, e.v.Y-8)
	e.world.SFXManager().PlayCustomEffect(sfx, e.v)
}

func (e *Entity) Kill() {
	fmt.Println("killing...")
	//play a cool death scene.
	e.PlayDeathEffect()

	//figure out a way to clean this up so i dont need to put this everywhere. maybe make a function in
	//sound manager called PlayWithDelay()?
	e.soundManager.DelayedPlay(150*time.Millisecond, "entdeath1_1")

	//unregister this entity from the world
	e.world.DeleteGameObject(e)
}

//todo: should be a parameter
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
	//am i near enough to be affected?
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
			go func() {
				time.Sleep(150 * time.Millisecond)
				e.soundManager.Play("enthit1_1")
				e.health -= 3
				if e.health <= 0 {
					e.Kill()
				}
				cb(e)
			}()

			return true
		}
	}

	return false
}

func (e *Entity) Update(tick int) {
	if e.animationManager.Ready() {
		e.animationManager.Select("Idle")
	}
	e.Draw(tick)
}

func (e *Entity) Draw(tick int) {
	target := e.world.Window
	sprite, animationFrame := e.animationManager.Next(tick)
	matrix := animationFrame.Matrix.Moved(e.v)
	mask := animationFrame.Mask
	sprite.DrawColorMask(target, matrix, mask)
}
