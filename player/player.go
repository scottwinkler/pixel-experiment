package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/pixel-experiment/animation"
)

type Player struct {
	Sprite *pixel.Sprite
	//Spritesheet      spritesheet.Spritesheet
	Speed            int
	X                int
	Y                int
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
		//Sprite:           spritesheet.Sprites[index],
		//Spritesheet:      spritesheet,
		Speed:  3, //default
		Window: window,
		X:      50,
		Y:      50,
		//AnimationManager: animation.NewAnimationManager(animations),
	}
	animationManager := animation.NewAnimationManager(animations)
	animationManager.Select("WalkDown")
	player.Sprite = animationManager.Selected.Spritesheet.Sprites[animationManager.Selected.Frames[0]]
	player.Matrix = pixel.Matrix(pixel.IM.Moved(pixel.V(float64(player.X), float64(player.Y))))
	player.SetAnimationManager(animationManager)
	return player
}

/*
func (p *Player) SetFrame(index int) {
	p.Sprite = p.Spritesheet.Sprites[index]
}*/

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
	matrix := pixel.Matrix(pixel.IM.Moved(pixel.V(float64(p.X), float64(p.Y))))
	p.Matrix = matrix
}

func (p *Player) Update(tick int) {
	//assume 60 ticks per second
	//so change animation once every 4 ticks

	var (
		camPos = pixel.V(float64(p.X), float64(p.Y)) //position camera centered on player
		//camSpeed = 1.0
		camZoom = 1.0
	)

	win := p.Window
	cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
	win.SetMatrix(cam)
	if win.Pressed(pixelgl.KeyLeft) {
		//fmt.Println("Animating WalkLeft")
		p.AnimationManager.Selected.SetPaused(false)
		p.AnimationManager.Select("WalkLeft")
		p.Move(LEFT)
	} else if win.Pressed(pixelgl.KeyRight) {
		p.AnimationManager.Selected.SetPaused(false)
		p.AnimationManager.Select("WalkRight")
		p.Move(RIGHT)
	} else if win.Pressed(pixelgl.KeyDown) {
		p.AnimationManager.Selected.SetPaused(false)
		p.AnimationManager.Select("WalkDown")
		p.Move(DOWN)
	} else if win.Pressed(pixelgl.KeyUp) {
		p.AnimationManager.Selected.SetPaused(false)
		p.AnimationManager.Select("WalkUp")
		p.Move(UP)

	} else {
		p.AnimationManager.Selected.SetPaused(true)
	}
	if tick == 0 {
		p.Sprite = p.AnimationManager.Selected.Next()
	}
	p.Draw()
}
func (p *Player) Draw() {
	p.Sprite.Draw(p.Window, p.Matrix)
}
