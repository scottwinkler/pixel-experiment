package animation

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/scottwinkler/pixel-experiment/utility"
)

type Animation struct {
	Name             string
	Spritesheet      *utility.Spritesheet
	Matrix           pixel.Matrix
	Frames           []int
	AnimationManager *AnimationManager
	Index            int
	Loop             bool
	Paused           bool
	skippable        bool
	done             bool
}

//utility function for converting a spritesheet based on a mapping of name:frames to an array of animations
func MappingToAnimations(spritesheet *utility.Spritesheet, mapping map[string]interface{}) []*Animation {
	var animations []*Animation
	for key, value := range mapping {
		attributes := value.(map[string]interface{})
		framesArr := attributes["Frames"].([]interface{})
		var frames []int
		for _, frame := range framesArr {
			frames = append(frames, int(frame.(float64)))
		}
		loop := attributes["Loop"].(bool)
		skippable := attributes["Skippable"].(bool)
		animations = append(animations, NewAnimation(spritesheet, key, frames, loop, skippable))
	}
	return animations
}

func NewAnimation(spritesheet *utility.Spritesheet, name string, frames []int, loop bool, skippable bool) *Animation {

	//smooth animation by reversing it after it has completed
	for i := len(frames) - 1; i > 0; i-- {
		frames = append(frames, frames[i])
	}

	animation := Animation{
		Name:        name,
		Spritesheet: spritesheet,
		Frames:      frames,
		Index:       0, //when Next() is first called it will go to index=0 which is what we want
		Loop:        loop,
		Paused:      true,
		Matrix:      spritesheet.Matrix,
		skippable:   skippable,
		done:        false,
	}
	return &animation
}

func (a *Animation) SetAnimationManager(animationManager *AnimationManager) {
	a.AnimationManager = animationManager
}

func (a *Animation) Reset() {
	a.done = false
	a.Index = 0
}

func (a *Animation) SetPaused(paused bool) {
	a.Paused = paused
}

func (a *Animation) Skippable() bool {
	return a.skippable
}

func (a *Animation) Done() bool {
	return a.done
}

func (a *Animation) Next() *pixel.Sprite {
	var frame int
	if !a.Paused {

		//fmt.Printf("index: %d/%d", a.Index, len(a.Frames)-1)
		a.Index++
		if a.Index > len(a.Frames)-1 {
			a.Index = 0
			if !a.Loop {
				//fmt.Printf("done!")
				a.done = true
			}
		}

		frame = a.Frames[a.Index]
	} else { //always return same frame if paused or not an appropriate time to change animations
		frame = a.Frames[a.Index]
	}

	//does the sprite need to be reflected?
	if frame < 0 {
		frame = int(math.Abs(float64(frame)))
		a.Matrix = a.Spritesheet.Matrix.Chained(pixel.IM.Rotated(pixel.ZV, -math.Pi).ScaledXY(pixel.ZV, pixel.V(-1, 1)).Rotated(pixel.ZV, math.Pi))
	} else {
		a.Matrix = a.Spritesheet.Matrix
	}
	return a.Spritesheet.Sprites[frame]
}
