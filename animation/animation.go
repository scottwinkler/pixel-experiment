package animation

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/scottwinkler/simple-rpg/utility"
)

//a helper struct for creating animations which have color masks unique for each frame
type AnimationFrame struct {
	Frame  int
	Matrix pixel.Matrix
	Mask   pixel.RGBA
}

//constructor utility of AnimationFrame objects
func NewAnimationFrame(frame int, matrix pixel.Matrix, mask pixel.RGBA) AnimationFrame {
	return AnimationFrame{
		Frame:  frame,
		Matrix: matrix,
		Mask:   mask,
	}
}

//properties for animation
type Animation struct {
	name             string
	spritesheet      *utility.Spritesheet
	frames           []AnimationFrame
	animationManager *AnimationManager
	index            int
	loop             bool
	paused           bool
	skippable        bool
	frameRate        int
	done             bool
}

//constructor for animation
func NewAnimation(spritesheet *utility.Spritesheet, name string, frames []AnimationFrame, loop bool, skippable bool, smooth bool, frameRate int) *Animation {
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
		index:       0,
		loop:        loop,
		paused:      true,
		skippable:   skippable,
		frameRate:   frameRate,
		done:        false,
	}
	return &animation
}

//getter method for index
func (a *Animation) Index() int {
	return a.index
}

//getter method for frameRate
func (a *Animation) FrameRate() int {
	return a.frameRate
}

//getter method for spritesheet
func (a *Animation) Spritesheet() *utility.Spritesheet {
	return a.spritesheet
}

//setter method for animationManager
func (a *Animation) SetAnimationManager(animationManager *AnimationManager) {
	a.animationManager = animationManager
}

//setter method for paused
func (a *Animation) SetPaused(paused bool) {
	a.paused = paused
}

//utility function for converting a spritesheet based on a mapping of name:frames to an array of animations

//idea: instead of passing the spritesheet explicitly, if the mapping contains a reference to the path of the spritesheet,
//then we can do this automatically
func MappingToAnimations(spritesheet *utility.Spritesheet, mapping map[string]interface{}) []*Animation {
	var animations []*Animation
	for name, value := range mapping {
		attributes := value.(map[string]interface{})
		framesArr := attributes["Frames"].([]interface{})
		var animationFrames []AnimationFrame
		for _, value := range framesArr {
			matrix := spritesheet.Matrix()
			frame := value.(float64)
			//does the sprite need to be reflected?
			if frame < 0 {
				frame = math.Abs(frame)
				matrix = matrix.Chained(pixel.IM.Rotated(pixel.ZV, -math.Pi).ScaledXY(pixel.ZV, pixel.V(-1, 1)).Rotated(pixel.ZV, math.Pi))
			}
			//assume that we do not use a custom matrix or color for effects created from spritesheets (maybe a bad guess?)
			animationFrames = append(animationFrames, NewAnimationFrame(int(frame), matrix, pixel.Alpha(1)))
		}

		loop := attributes["Loop"].(bool)
		skippable := attributes["Skippable"].(bool)
		smooth := attributes["Smooth"].(bool)
		frameRate := int(attributes["FrameRate"].(float64)) //cast to int because ain't dealing with floating point nonsense
		animations = append(animations, NewAnimation(spritesheet, name, animationFrames, loop, skippable, smooth, frameRate))
	}
	return animations
}

//resets animation to initial state
func (a *Animation) Reset() {
	a.done = false
	a.index = 0
}

//returns current frame data. ready only, does not have side effects
func (a *Animation) Current() (*pixel.Sprite, AnimationFrame) {
	frame := a.frames[a.index].Frame
	return a.spritesheet.Sprites()[frame], a.frames[a.index]
}

//returns the next sprite and frame data. should be invoked on every update.
func (a *Animation) Next(tick int) (*pixel.Sprite, AnimationFrame) {
	if !a.paused && tick%(60/a.frameRate) == 0 {
		a.index++
		if a.index > len(a.frames)-1 {
			a.index = 0
			if !a.loop {
				a.done = true
			}
		}
	}
	frame := a.frames[a.index].Frame
	return a.spritesheet.Sprites()[frame], a.frames[a.index]
}
