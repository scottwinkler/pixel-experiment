package animation

import (
	"github.com/faiface/pixel"
	"github.com/scottwinkler/pixel-experiment/spritesheet"
)

type Animation struct {
	Name             string
	Spritesheet      *spritesheet.Spritesheet
	Frames           []int
	AnimationManager *AnimationManager
	Index            int
	Loop             bool
	Paused           bool
}

func NewAnimation(spritesheet *spritesheet.Spritesheet, name string, frames []int) *Animation {
	//smooth animation
	for i := len(frames) - 1; i > 0; i-- {
		frames = append(frames, frames[i])
	}

	animation := Animation{
		Name:        name,
		Spritesheet: spritesheet,
		Frames:      frames,
		Index:       0,
		Loop:        true,
		Paused:      true,
	}
	return &animation
}

func (a *Animation) SetAnimationManager(animationManager *AnimationManager) {
	a.AnimationManager = animationManager
}

func (a *Animation) Reset() {
	a.Index = 0
}

func (a *Animation) SetPaused(paused bool) {
	a.Paused = paused
}

func (a *Animation) Next() *pixel.Sprite {
	nextIndex := a.Index
	if !a.Paused {
		if a.Loop {
			nextIndex = (a.Index + 1) % len(a.Frames)
		}
	}
	a.Index = nextIndex
	frame := a.Frames[a.Index]
	return a.Spritesheet.Sprites[frame]
}
