package entity

import (
	"image/color"
	"math"

	"github.com/scottwinkler/simple-rpg/enum"
	"github.com/scottwinkler/simple-rpg/utility"

	"github.com/scottwinkler/simple-rpg/sound"

	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"github.com/scottwinkler/simple-rpg/animation"
	"github.com/scottwinkler/simple-rpg/sfx"
	"github.com/scottwinkler/simple-rpg/world"
)

//Data -- Subset of entity configuration that excludes ephemeral data
type Data struct {
	Animations []*animation.Animation
	Sounds     []*sound.Sound
	Health     float64
	Speed      float64
	R          float64
	Material   string
	Name       string
	Color      color.Color
}

//Configuration -- needed so we don't have an insane amount of arguments to the constructor
type Configuration struct {
	V    pixel.Vec
	W    *world.World
	Data *Data
	//color color.Color possibly a tint color?  if not specified then no tint?
}

//Entity -- implements world.GameObject, the underlying struct for all living entities
type Entity struct {
	id               string
	parent           interface{} //the container of entity, nil if there isn't one
	speed            float64
	v                pixel.Vec
	r                float64 //used for collider calculations
	animationManager *animation.Manager
	soundManager     *sound.Manager
	world            *world.World
	direction        int
	health           float64
	name             string //the kind of entity e.g. slime, skeleton, goblin, whatever
	material         string
	color            color.Color
	controller       Controller
	damage           float64
}

//NewEntity -- constructor for Entity
func NewEntity(config *Configuration, parent interface{}) *Entity {
	animationManager := animation.NewManager(config.Data.Animations)
	animationManager.Select("Idle") //every entity should have an idle animation
	soundManager := sound.NewManager(config.Data.Sounds)
	id := xid.New().String()
	name := config.Data.Name
	entity := &Entity{
		id:               id,
		parent:           parent,
		speed:            config.Data.Speed,
		v:                config.V,
		r:                config.Data.R,
		animationManager: animationManager,
		soundManager:     soundManager,
		world:            config.W,
		direction:        enum.Direction.Down, //default
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

//ID -- getter method for id
func (e *Entity) ID() string {
	return e.id
}

//V -- getter method for v
func (e *Entity) V() pixel.Vec {
	return e.v
}

//R -- getter method for r
func (e *Entity) R() float64 {
	return e.r
}

//Speed -- getter method for speed
func (e *Entity) Speed() float64 {
	return e.speed
}

//Direction -- getter method for direction
func (e *Entity) Direction() int {
	return e.direction
}

//Material -- getter method for material
func (e *Entity) Material() string {
	return e.material
}

//World -- getter method for world
func (e *Entity) World() *world.World {
	return e.world
}

//AnimationManager -- getter method for animationManager
func (e *Entity) AnimationManager() *animation.Manager {
	return e.animationManager
}

//SoundManager -- getter method for soundManager
func (e *Entity) SoundManager() *sound.Manager {
	return e.soundManager
}

//Damage -- getter method for damage
func (e *Entity) Damage() float64 {
	return e.damage
}

//Name -- getter method for name
func (e *Entity) Name() string {
	return e.name
}

//IsDead -- convenience function to know if an entity is dead
func (e *Entity) IsDead() bool {
	return e.health <= 0
}

//a calculated effect to play once when the entity dies
func (e *Entity) makeDeathEffect() *sfx.SFX {
	var sfxFrames []sfx.Frame
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
		sfxFrame := sfx.NewFrame(frame, matrix, mask)
		sfxFrames = append(sfxFrames, sfxFrame)
	}
	frameRate := animation.FrameRate()
	sfx := sfx.NewSFX(animation.Spritesheet(), "EntityDeath", sfxFrames, false, frameRate)
	return sfx

}

//Kill -- the method which safely kills this object and does so in a cool way
func (e *Entity) Kill() {
	//play a cool death scene.
	//e.soundManager.Play("death0") -> should be moved into callback function
	effect := e.makeDeathEffect()
	e.world.SFXManager().PlayCustomEffect(effect, e.v)
	e.world.DeleteGameObject(e)
}

//Move to given direction. Returns a boolean for success condition
func (e *Entity) Move(direction int) bool {
	nextPos := pixel.V(e.v.X, e.v.Y)
	var moveAnimation string
	switch direction {
	case enum.Direction.Left:
		nextPos.X -= e.speed
		moveAnimation = "MoveLeft"
	case enum.Direction.Right:
		nextPos.X += e.speed
		moveAnimation = "MoveRight"
	case enum.Direction.Down:
		nextPos.Y -= e.speed
		moveAnimation = "MoveDown"
	case enum.Direction.Up:
		nextPos.Y += e.speed
		moveAnimation = "MoveUp"
	}
	e.direction = direction
	e.animationManager.Select(moveAnimation)
	if !e.world.Collides(e.ID(), nextPos, e.r) {
		e.v = nextPos
		return true
	}
	return false
}

//HandleHit -- the method which handles hit events
func (e *Entity) HandleHit(s world.GameObject, cb world.Callback) bool {
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
			relativeDir = enum.Direction.Up
		} else if top {
			relativeDir = enum.Direction.Left
		} else if right {
			relativeDir = enum.Direction.Right
		} else {
			relativeDir = enum.Direction.Down
		}
		//is the source facing the right direction?
		if relativeDir == s.Direction() {
			//let the controller decide if it wants to accept the hit or not
			isHit := e.controller.HitCallback(s)
			am := e.animationManager
			if isHit {
				switch e.direction {
				case enum.Direction.Left:
					am.Select("HitLeft")
				case enum.Direction.Right:
					am.Select("HitRight")
				case enum.Direction.Down:
					am.Select("HitDown")
				case enum.Direction.Up:
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

//Attack -- the method for attacking in a given direction
func (e *Entity) Attack(direction int) {
	e.direction = direction
	switch direction {
	case enum.Direction.Left:
		e.animationManager.Select("AttackLeft")
	case enum.Direction.Right:
		e.animationManager.Select("AttackRight")
	case enum.Direction.Down:
		e.animationManager.Select("AttackDown")
	case enum.Direction.Up:
		e.animationManager.Select("AttackUp")
	}
	cb := e.controller.AttackCallback
	e.world.HitEvent(e, cb)
}

//Update -- the method that gets called by the main game loop
func (e *Entity) Update(tick int) {
	e.controller.Update(tick) //outsource all our work like the plebs we are
	e.Draw(tick)
}

//Draw -- draws the entity onto the window target
func (e *Entity) Draw(tick int) {
	target := e.world.Window
	sprite, animationFrame := e.animationManager.Next(tick)
	matrix := animationFrame.Matrix.Moved(e.v)
	mask := animationFrame.Mask
	sprite.DrawColorMask(target, matrix, mask)
}
