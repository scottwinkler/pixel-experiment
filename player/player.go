package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/scottwinkler/pixel-experiment/animation"
	"github.com/scottwinkler/pixel-experiment/world"
)

type Player struct {
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

func (p *Player) SetAnimationManager(animationManager *animation.AnimationManager) {
	p.AnimationManager = animationManager
}
func NewPlayer(animations []*animation.Animation, world *world.World) *Player {
	player := &Player{
		Speed: 3, //default
		World: world,
		V:     pixel.V(200, 150),
	}
	animationManager := animation.NewAnimationManager(animations)
	animationManager.Select("Idle")
	player.Sprite = animationManager.Selected.Spritesheet.Sprites[animationManager.Selected.Frames[0]]
	player.Matrix = pixel.Matrix(pixel.IM.Moved(player.V).Scaled(player.V, 0.5))
	player.SetAnimationManager(animationManager)
	return player
}

//returns true if the player would have a collision at the given point
func (p *Player) Collides(v pixel.Vec) bool {
	if !p.World.Tilemap.Bounds().Contains(v) {
		return true //out of bounds!
	}
	tile := p.World.Tilemap.GetTileAtPosition(v, 3) //check collision layer
	if tile == nil {
		return false
	}
	//fmt.Printf("got tile: v:%v, isCollidable %t, gid: %d", tile.V, tile.IsCollidable, tile.Gid)
	return tile.IsCollidable
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
	if !p.Collides(nextPos) {
		p.V = nextPos
	}
	matrix := pixel.Matrix(pixel.IM.Moved(p.V).Scaled(p.V, 0.5))
	p.Matrix = matrix
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
	p.Sprite.Draw(p.World.Window, p.Matrix)
}
