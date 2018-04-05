package entity

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/pixel-experiment/animation"
	"github.com/scottwinkler/pixel-experiment/world"
)

type Entity struct {
	Sprite           *pixel.Sprite
	Speed            float64
	V                pixel.Vec
	Matrix           pixel.Matrix
	AnimationManager *animation.AnimationManager
	World            *world.World
}

const (
	LEFT  = 0
	RIGHT = 1
	DOWN  = 2
	UP    = 3
)

func (e *Entity) SetAnimationManager(animationManager *animation.AnimationManager) {
	e.AnimationManager = animationManager
}
func NewEntity(animations []*animation.Animation, world *world.World) *Entity {
	entity := &Entity{
		Speed: 3, //default
		World: world,
		V:     pixel.V(200, 250), //probably shouldnt hardcode this
	}
	animationManager := animation.NewAnimationManager(animations)
	animationManager.Select("Idle")
	entity.Sprite = animationManager.Selected.Spritesheet.Sprites[animationManager.Selected.Frames[0]]
	entity.Matrix = pixel.Matrix(pixel.IM.Moved(entity.V).Scaled(entity.V, 0.5))
	entity.SetAnimationManager(animationManager)
	return entity
}

//returns true if the entity would have a collision at the given point
func (e *Entity) Collides(v pixel.Vec) bool {
	if !e.World.Tilemap.Bounds().Contains(v) {
		return true //out of bounds!
	}
	tile := e.World.Tilemap.GetTileAtPosition(v, "Collision") //check tile in collision layer
	if tile == nil {
		return false
	}
	return tile.IsCollidable
}

func (e *Entity) Move(direction int) {
	nextPos := pixel.V(e.V.X, e.V.Y)
	switch direction {
	case LEFT:
		nextPos.X -= e.Speed
	case RIGHT:
		nextPos.X += e.Speed
	case DOWN:
		nextPos.Y -= e.Speed
	case UP:
		nextPos.Y += e.Speed
	}
	if !e.Collides(nextPos) {
		e.V = nextPos
	}
	matrix := pixel.Matrix(pixel.IM.Moved(e.V).Scaled(e.V, -0.5)) //.Rotated(e.V, -math.Pi)
	//matrix = Reflection(matrix, e.V, math.Pi)
	e.Matrix = matrix
}

func (e *Entity) Update(tick int) {
	win := e.World.Window
	if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
		e.AnimationManager.Select("MoveLeft")
		e.Move(LEFT)
	} else if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
		e.AnimationManager.Select("MoveLeft")
		e.Move(RIGHT)
	} else if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
		e.AnimationManager.Select("MoveDown")
		e.Move(DOWN)
	} else if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
		e.AnimationManager.Select("MoveUp")
		e.Move(UP)
	} else {
		e.AnimationManager.Select("Idle")
	}
	if tick == 0 {
		e.Sprite = e.AnimationManager.Selected.Next()
	}
	e.Draw()
}
func (e *Entity) Draw() {
	e.Sprite.Draw(e.World.Window, e.Matrix)
}
