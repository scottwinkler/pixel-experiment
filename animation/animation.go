package animation

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/scottwinkler/pixel-experiment/spritesheet"
)

type Animation struct {
	Name             string
	Spritesheet      *spritesheet.Spritesheet
	Matrix           pixel.Matrix
	Frames           []int
	AnimationManager *AnimationManager
	Index            int
	Loop             bool
	Paused           bool
}

//utility function for converting a spritesheet based on a mapping of name:frames to an array of animations
func AnimationsFromSpritesheet(spritesheet *spritesheet.Spritesheet, mapping map[string][]int) []*Animation {
	var animations []*Animation
	for name, frames := range mapping {
		animations = append(animations, NewAnimation(spritesheet, name, frames))
	}
	return animations
}

func NewAnimation(spritesheet *spritesheet.Spritesheet, name string, frames []int) *Animation {

	//smooth animation by reversing it after it has completed
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
		Matrix:      pixel.IM,
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

	//does the sprite need to be reflected?
	if frame < 0 {
		frame = int(math.Abs(float64(frame)))
		a.Matrix = a.Spritesheet.Matrix.Chained(pixel.IM.Rotated(pixel.ZV, -math.Pi).ScaledXY(pixel.ZV, pixel.V(-1, 1)).Rotated(pixel.ZV, math.Pi))
		//a.Matrix = Reflected(pixel.IM, origin, math.Pi).Scaled(origin, -1)
	} else {
		a.Matrix = pixel.IM
	}
	return a.Spritesheet.Sprites[frame]
}

func Reflected(m pixel.Matrix, around pixel.Vec, angle float64) pixel.Matrix {
	sin2t, cos2t := math.Sincos(2 * angle)
	m[4], m[5] = m[4]-around.X, m[5]-around.Y
	m = m.Chained(pixel.Matrix{cos2t, sin2t, sin2t, -cos2t, 0, 0})
	m[4], m[5] = m[4]+around.X, m[5]+around.Y
	return m
}
