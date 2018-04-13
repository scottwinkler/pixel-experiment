package animation

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/scottwinkler/simple-rpg/utility"
)

type Animation struct {
	name             string
	spritesheet      *utility.Spritesheet
	matrix           pixel.Matrix
	frames           []int
	animationManager *AnimationManager
	index            int
	loop             bool
	paused           bool
	skippable        bool
	frameRate        int
	done             bool
}

//getter function for matrix
func (a *Animation) Matrix() pixel.Matrix {
	return a.matrix
}

//getter function for spritesheet
func (a *Animation) Spritesheet() *utility.Spritesheet {
	return a.spritesheet
}

//utility function for converting a spritesheet based on a mapping of name:frames to an array of animations

//idea: instead of passing the spritesheet explicitly, if the mapping contains a reference to the path of the spritesheet,
//then we can do this automatically
func MappingToAnimations(spritesheet *utility.Spritesheet, mapping map[string]interface{}) []*Animation {
	var animations []*Animation
	for name, value := range mapping {
		attributes := value.(map[string]interface{})
		framesArr := attributes["Frames"].([]interface{})
		var frames []int
		for _, frame := range framesArr {
			frames = append(frames, int(frame.(float64)))
		}
		loop := attributes["Loop"].(bool)
		skippable := attributes["Skippable"].(bool)
		smooth := attributes["Smooth"].(bool)
		frameRate := int(attributes["FrameRate"].(float64)) //cast to int because ain't dealing with floating point nonsense
		animations = append(animations, NewAnimation(spritesheet, name, frames, loop, skippable, smooth, frameRate))
	}
	return animations
}

func NewAnimation(spritesheet *utility.Spritesheet, name string, frames []int, loop bool, skippable bool, smooth bool, frameRate int) *Animation {
	if smooth {
		//smooth animation by padding it with frames appended in reverse order
		for i := len(frames) - 1; i > 0; i-- {
			frames = append(frames, frames[i])
		}
	}

	animation := Animation{
		name:        name,
		spritesheet: spritesheet,
		frames:      frames,
		index:       0, //when Next() is first called it will go to index=0 which is what we want
		loop:        loop,
		paused:      true,
		matrix:      spritesheet.Matrix(),
		skippable:   skippable,
		frameRate:   frameRate,
		done:        false,
	}
	return &animation
}

func (a *Animation) SetAnimationManager(animationManager *AnimationManager) {
	a.animationManager = animationManager
}

func (a *Animation) Reset() {
	a.done = false
	a.index = 0
}

func (a *Animation) SetPaused(paused bool) {
	a.paused = paused
}

func (a *Animation) Skippable() bool {
	return a.skippable
}

func (a *Animation) Done() bool {
	return a.done
}

func (a *Animation) Next(tick int) *pixel.Sprite {
	var frame int
	if !a.paused && tick%(60/a.frameRate) == 0 {
		a.index++
		if a.index > len(a.frames)-1 {
			a.index = 0
			if !a.loop {
				a.done = true
			}
		}
		frame = a.frames[a.index]
	} else {
		//always return same frame if paused or not an appropriate time to change animations
		frame = a.frames[a.index]
	}

	//does the sprite need to be reflected?
	if frame < 0 {
		frame = int(math.Abs(float64(frame)))
		a.matrix = a.spritesheet.Matrix().Chained(pixel.IM.Rotated(pixel.ZV, -math.Pi).ScaledXY(pixel.ZV, pixel.V(-1, 1)).Rotated(pixel.ZV, math.Pi))
	} else {
		a.matrix = a.spritesheet.Matrix()
	}
	return a.spritesheet.Sprites()[frame]
}
