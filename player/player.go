package player

import (
	"github.com/faiface/pixel"
	"github.com/scottwinkler/pixel-experiment/spritesheet"
)

type Player struct {
	Sprite      *pixel.Sprite
	Spritesheet spritesheet.Spritesheet
}

func NewPlayer(index int, spritesheet spritesheet.Spritesheet) Player {
	player := Player{
		Sprite:      spritesheet.Sprites[index],
		Spritesheet: spritesheet,
	}
	return player
}

func (p *Player) SetFrame(index int) {
	p.Sprite = p.Spritesheet.Sprites[index]
}

func (p *Player) Draw(t pixel.Target) {
	p.Sprite.Draw(t, pixel.IM)
}
