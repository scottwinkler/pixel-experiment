package entity

import (
	"image/color"
	"math"
	"time"

	"github.com/scottwinkler/simple-rpg/utility"

	"github.com/scottwinkler/simple-rpg/sound"

	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"github.com/scottwinkler/simple-rpg/animation"
	"github.com/scottwinkler/simple-rpg/sfx"
	"github.com/scottwinkler/simple-rpg/world"
)

//Subset of entity configuration that excludes ephemeral data
type EntityData struct {
	Animations []*animation.Animation
	Sounds     []*sound.Sound
	Health     float64
	Speed      float64
	R          float64
	Material   string
	Name       string
	Color      color.Color
}

//needed so we don't have an insane amount of arguments to the constructor
type EntityConfiguration struct {
	V    pixel.Vec
	W    *world.World
	Data *EntityData
	//color color.Color possibly a tint color?  if not specified then no tint?
}

type Entity struct {
	id               string
	speed            float64
	v                pixel.Vec
	r                float64 //used for collider calculations
	animationManager *animation.AnimationManager
	soundManager     *sound.SoundManager
	world            *world.World
	direction        int
	health           float64
	name             string //the kind of entity e.g. slime, skeleton, goblin, whatever
	material         string
	color            color.Color
}

func NewEntity(config *EntityConfiguration) *Entity {
	animationManager := animation.NewAnimationManager(config.Data.Animations)
	animationManager.Select("Idle") //every entity should have an idle frame
	soundManager := sound.NewSoundManager(config.Data.Sounds)
	id := xid.New().String()
	entity := &Entity{
		id:               id,
		speed:            config.Data.Speed,
		v:                config.V,
		r:                config.Data.R,
		animationManager: animationManager,
		soundManager:     soundManager,
		world:            config.W,
		direction:        world.DOWN, //default
		health:           config.Data.Health,
		name:             config.Data.Name,
		material:         config.Data.Material,
		color:            config.Data.Color,
	}
	config.W.AddGameObject(config.Data.Name, entity) //register self with world
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

//getter method for material
func (e *Entity) Material() string {
	return e.material
}

//convenience function
func (e *Entity) IsDead() bool {
	return e.health <= 0
}

//a calculated effect to play once when the entity dies
func (e *Entity) MakeDeathEffect() *sfx.SFX {
	var sfxFrames []sfx.SFXFrame
	animation := e.animationManager.Selected()
	_, animationFrame := animation.Current()
	frame := animationFrame.Frame
	//8 frames feels right. this is approximately how long the death sound takes at a frameRate of 10
	framesCount := 8
	for i := framesCount; i >= 0; i-- {
		scaleFactor := float64(i) / float64(framesCount)
		//this exponential function makes it more interesting than a simple linear interpolation
		scaleFactor = math.Exp(1 - 1/math.Pow(scaleFactor, 2))
		matrix := animationFrame.Matrix.ScaledXY(pixel.ZV, pixel.V(1, scaleFactor))
		//add a color mask to the sprite, and make it incrementally smaller
		mask := utility.ToRGBA(e.color, 0.8)
		sfxFrame := sfx.NewSFXFrame(frame, matrix, mask)
		sfxFrames = append(sfxFrames, sfxFrame)
	}
	frameRate := animation.FrameRate()
	sfx := sfx.NewSFX(animation.Spritesheet(), "EntityDeath", sfxFrames, false, frameRate)
	return sfx

}

func (e *Entity) Kill() {
	//play a cool death scene.
	e.soundManager.Play("death0")
	effect := e.MakeDeathEffect()
	e.world.SFXManager().PlayCustomEffect(effect, e.v)
	e.world.DeleteGameObject(e)
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

				e.health -= 3
				if e.health <= 0 {
					e.Kill()
				} else {
					e.soundManager.Play("hit0")
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
