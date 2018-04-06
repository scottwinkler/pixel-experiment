package player

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/rs/xid"
	"github.com/scottwinkler/pixel-experiment/animation"
	"github.com/scottwinkler/pixel-experiment/world"
)

type Player struct {
	id               string
	Sprite           *pixel.Sprite
	Speed            float64
	V                pixel.Vec
	R                float64 //used for collider calculations
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

func (p *Player) SetAnimationManager(animationManager *animation.AnimationManager) {
	p.AnimationManager = animationManager
}
func NewPlayer(v pixel.Vec, animations []*animation.Animation, world *world.World) *Player {
	animationManager := animation.NewAnimationManager(animations)
	animationManager.Select("Idle")
	sprite := animationManager.Selected.Spritesheet.Sprites[animationManager.Selected.Frames[0]]
	matrix := pixel.Matrix(pixel.IM.Moved(v))
	radius := float64(world.Tilemap.TileHeight / 3) //not a great solution for rectangles
	fmt.Printf("r: %f", radius)
	id := xid.New().String()
	player := &Player{
		id:               id,
		Sprite:           sprite,
		Speed:            3, //default
		V:                v,
		R:                radius,
		Matrix:           matrix,
		AnimationManager: animationManager,
		World:            world,
	}
	return player
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) Move(direction int) {
	nextPos := pixel.V(p.V.X, p.V.Y)
	switch direction {
	case LEFT:
		nextPos.X -= p.Speed
	case RIGHT:
		nextPos.X += p.Speed
	case DOWN:
		nextPos.Y -= p.Speed
	case UP:
		nextPos.Y += p.Speed
	}
	if !p.World.Collides(p.Id(), nextPos, p.R) {
		p.V = nextPos
	}
	//update matrix and collision circle
	matrix := pixel.IM.Moved(p.V)
	p.Matrix = matrix
}

func (p *Player) Collider() (pixel.Vec, float64) {
	return p.V, p.R
}

func (p *Player) Update(tick int) {
	var (
		camPos = p.V //position camera centered on player
		//camSpeed = 1.0
		camZoom = 1.0
	)

	win := p.World.Window
	cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
	win.SetMatrix(cam)
	if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
		p.AnimationManager.Select("RunLeft")
		p.Move(LEFT)
	} else if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
		p.AnimationManager.Select("RunRight")
		p.Move(RIGHT)
	} else if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
		p.AnimationManager.Select("RunDown")
		p.Move(DOWN)
	} else if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
		p.AnimationManager.Select("RunUp")
		p.Move(UP)
	} else if win.Pressed(pixelgl.MouseButtonLeft) {
		//determine what quadrant relative to the player the mouse click happened
		mouse := cam.Unproject(win.MousePosition()).Sub(p.V)
		top := mouse.Y >= mouse.X    //above line y=x?
		right := mouse.Y >= -mouse.X //above line y=-x?
		if top && right {
			p.AnimationManager.Select("AttackUp")
		} else if top {
			p.AnimationManager.Select("AttackLeft")
		} else if right {
			p.AnimationManager.Select("AttackRight")
		} else {
			p.AnimationManager.Select("AttackDown")
		}
	} else {
		p.AnimationManager.Select("Idle")
	}
	if tick == 0 {
		p.Sprite = p.AnimationManager.Selected.Next()
	}
	p.Draw()
}
func (p *Player) Draw() {
	animation := p.AnimationManager.Selected
	//chained methods so that we first scale by spritesheet size, then by reflection, then by position
	matrix := animation.Spritesheet.Matrix.Chained(p.Matrix)
	p.Sprite.Draw(p.World.Window, matrix)
}
