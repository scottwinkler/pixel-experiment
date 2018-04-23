package animation

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/scottwinkler/simple-rpg/utility"
	"golang.org/x/image/colornames"
)

//Frame -- helper struct for creating animations which have color masks unique for each frame
type Frame struct {
	Frame  int
	Matrix pixel.Matrix
	Mask   pixel.RGBA
}

//NewFrame -- constructor utility of Frame objects
func NewFrame(frame int, matrix pixel.Matrix, mask pixel.RGBA) Frame {
	return Frame{
		Frame:  frame,
		Matrix: matrix,
		Mask:   mask,
	}
}

//Animation -- properties for animation
type Animation struct {
	name        string
	spritesheet *utility.Spritesheet
	frames      []Frame
	manager     *Manager
	index       int
	loop        bool
	paused      bool
	skippable   bool
	frameRate   int
	done        bool
}

//NewAnimation -- constructor for animation
func NewAnimation(spritesheet *utility.Spritesheet, name string, frames []Frame, loop bool, skippable bool, smooth bool, frameRate int) *Animation {
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

//Name -- getter method for name
func (a *Animation) Name() string {
	return a.name
}

//Index -- etter method for index
func (a *Animation) Index() int {
	return a.index
}

//FrameRate -- getter method for frameRate
func (a *Animation) FrameRate() int {
	return a.frameRate
}

//Spritesheet -- getter method for spritesheet
func (a *Animation) Spritesheet() *utility.Spritesheet {
	return a.spritesheet
}

//SetManager -- setter method for manager
func (a *Animation) SetManager(manager *Manager) {
	a.manager = manager
}

//SetPaused -- setter method for paused
func (a *Animation) SetPaused(paused bool) {
	a.paused = paused
}

//MappingToAnimations -- utility function for converting a spritesheet based on a mapping of name:frames to an array of animations
//idea: instead of passing the spritesheet explicitly, if the mapping contains a reference to the path of the spritesheet,
//then we can do this automatically
func MappingToAnimations(spritesheet *utility.Spritesheet, mapping map[string]interface{}) []*Animation {
	var animations []*Animation
	for name, value := range mapping {
		attributes := value.(map[string]interface{})
		mask := pixel.Alpha(1) //default

		//if a mask was supplied, then use that
		if value, ok := attributes["Mask"]; ok {
			maskData := value.(map[string]interface{})
			colorName := maskData["Color"].(string)
			color := colornames.Map[colorName]
			alpha := maskData["Alpha"].(float64)
			mask = utility.ToRGBA(color, alpha)
		}
		framesArr := attributes["Frames"].([]interface{})
		var animationFrames []Frame
		for _, value := range framesArr {
			matrix := spritesheet.Matrix()
			frame := value.(float64)
			//does the sprite need to be reflected?
			if frame < 0 {
				frame = math.Abs(frame)
				matrix = matrix.Chained(pixel.IM.Rotated(pixel.ZV, -math.Pi).ScaledXY(pixel.ZV, pixel.V(-1, 1)).Rotated(pixel.ZV, math.Pi))
			}
			//assume that we do not use a custom matrix or color for effects created from spritesheets (maybe a bad guess?)
			animationFrames = append(animationFrames, NewFrame(int(frame), matrix, mask))
		}

		loop := attributes["Loop"].(bool)
		skippable := attributes["Skippable"].(bool)
		smooth := attributes["Smooth"].(bool)

		frameRate := int(attributes["FrameRate"].(float64)) //cast to int because ain't dealing with floating point nonsense
		animations = append(animations, NewAnimation(spritesheet, name, animationFrames, loop, skippable, smooth, frameRate))
	}
	return animations
}

//Reset -- resets animation to initial state
func (a *Animation) Reset() {
	a.done = false
	a.index = 0
}

//Current -- returns current frame data. ready only, does not have side effects
func (a *Animation) Current() (*pixel.Sprite, Frame) {
	frame := a.frames[a.index].Frame
	return a.spritesheet.Sprites()[frame], a.frames[a.index]
}

//Next -- returns the next sprite and frame data. should be invoked on every update.
func (a *Animation) Next(tick int) (*pixel.Sprite, Frame) {
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
