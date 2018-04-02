package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/pixel-experiment/animation"
)

type Player struct {
	Sprite           *pixel.Sprite
	Speed            int
	X                int
	Y                int
	V                pixel.Vec
	Matrix           pixel.Matrix
	AnimationManager *animation.AnimationManager
	Window           *pixelgl.Window
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
func NewPlayer(animations []*animation.Animation, window *pixelgl.Window) *Player {
	player := &Player{
		Speed:  3, //default
		Window: window,
		X:      200,
		Y:      150,
		V:      pixel.V(float64(200), float64(150)),
	}
	animationManager := animation.NewAnimationManager(animations)
	animationManager.Select("Idle")
	player.Sprite = animationManager.Selected.Spritesheet.Sprites[animationManager.Selected.Frames[0]]
	player.Matrix = pixel.Matrix(pixel.IM.Moved(player.V).Scaled(player.V, 0.5))
	player.SetAnimationManager(animationManager)
	return player
}

func (p *Player) Move(direction int) {
	switch direction {
	case LEFT:
		p.X -= p.Speed
	case RIGHT:
		p.X += p.Speed
	case DOWN:
		p.Y -= p.Speed
	case UP:
		p.Y += p.Speed
	}
	p.V = pixel.V(float64(p.X), float64(p.Y))
	matrix := pixel.Matrix(pixel.IM.Moved(p.V).Scaled(p.V, 0.5))
	p.Matrix = matrix
}

func (p *Player) Update(tick int) {
	var (
		camPos = p.V //position camera centered on player
		//camSpeed = 1.0
		camZoom = 1.0
	)

	win := p.Window
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
	p.Sprite.Draw(p.Window, p.Matrix)
}
