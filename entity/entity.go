package entity

import (
	"image/color"
	"math"

	"github.com/scottwinkler/simple-rpg/utility"

	"github.com/scottwinkler/simple-rpg/sound"

	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"github.com/scottwinkler/simple-rpg/animation"
	"github.com/scottwinkler/simple-rpg/sfx"
	"github.com/scottwinkler/simple-rpg/world"
)

//list of all valid entity names
const (
	ENTITY_SLIME   = "slime"
	ENTITY_SPAWNER = "spawner"
	ENTITY_PLAYER  = "player"
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
	controller       controller
	damage           float64
}

func NewEntity(config *EntityConfiguration) *Entity {
	animationManager := animation.NewAnimationManager(config.Data.Animations)
	animationManager.Select("Idle") //every entity should have an idle animation
	soundManager := sound.NewSoundManager(config.Data.Sounds)
	id := xid.New().String()
	name := config.Data.Name
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
		name:             name,
		material:         config.Data.Material,
		color:            config.Data.Color,
		damage:           3.0, //fix this
	}
	entity.controller = ControllerMap[name](entity) //fetch ai to use based on entity name
	config.W.AddGameObject(name, entity)            //register self with world
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

//getter method for world
func (e *Entity) World() *world.World {
	return e.world
}

//getter method for animationManager
func (e *Entity) AnimationManager() *animation.AnimationManager {
	return e.animationManager
}

//getter method for soundManager
func (e *Entity) SoundManager() *sound.SoundManager {
	return e.soundManager
}

//getter method for damage
func (e *Entity) Damage() float64 {
	return e.damage
}

//getter method for name
func (e *Entity) Name() string {
	return e.name
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
		mask := utility.ToRGBA(e.color, 0.6)
		sfxFrame := sfx.NewSFXFrame(frame, matrix, mask)
		sfxFrames = append(sfxFrames, sfxFrame)
	}
	frameRate := animation.FrameRate()
	sfx := sfx.NewSFX(animation.Spritesheet(), "EntityDeath", sfxFrames, false, frameRate)
	return sfx

}

func (e *Entity) Kill() {
	//play a cool death scene.
	//e.soundManager.Play("death0") -> should be moved into callback function
	effect := e.MakeDeathEffect()
	e.world.SFXManager().PlayCustomEffect(effect, e.v)
	e.world.DeleteGameObject(e)
}

//Move to given direction. Returns a boolean for success condition
func (e *Entity) Move(direction int) bool {
	nextPos := pixel.V(e.v.X, e.v.Y)
	var moveAnimation string
	switch direction {
	case world.LEFT:
		nextPos.X -= e.speed
		moveAnimation = "MoveLeft"
	case world.RIGHT:
		nextPos.X += e.speed
		moveAnimation = "MoveRight"
	case world.DOWN:
		nextPos.Y -= e.speed
		moveAnimation = "MoveDown"
	case world.UP:
		nextPos.Y += e.speed
		moveAnimation = "MoveUp"
	}
	e.direction = direction
	if !e.world.Collides(e.Id(), nextPos, e.r) {
		e.v = nextPos
		return false
	}
	e.animationManager.Select(moveAnimation)
	return true
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
			//let the controller decide if it wants to accept the hit or not
			isHit := e.controller.HitCallback(s)
			am := e.animationManager
			if isHit {
				switch e.direction {
				case world.LEFT:
					am.Select("HitLeft")
				case world.RIGHT:
					am.Select("HitRight")
				case world.DOWN:
					am.Select("HitDown")
				case world.UP:
					am.Select("HitUp")
				}
				e.health -= s.Damage()
				if e.health <= 0 {
					e.Kill()
				}
			}
		}
		cb(e)
		return true
	}
	return false
}

//shitty copy paste
func (e *Entity) Attack(direction int) {
	e.direction = direction
	switch direction {
	case world.LEFT:
		e.animationManager.Select("AttackLeft")
	case world.RIGHT:
		e.animationManager.Select("AttackRight")
	case world.DOWN:
		e.animationManager.Select("AttackDown")
	case world.UP:
		e.animationManager.Select("AttackUp")
	}
	cb := e.controller.AttackCallback
	e.world.HitEvent(e, cb)
}

func (e *Entity) Update(tick int) {
	e.controller.Update(tick) //outsource all our work like the plebs we are
	e.Draw(tick)
}

func (e *Entity) Draw(tick int) {
	target := e.world.Window
	sprite, animationFrame := e.animationManager.Next(tick)
	matrix := animationFrame.Matrix.Moved(e.v)
	mask := animationFrame.Mask
	sprite.DrawColorMask(target, matrix, mask)
}
