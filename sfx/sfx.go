package sfx

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"github.com/scottwinkler/simple-rpg/utility"
)

//SFX -- an SFX is a simplified animation that only needs to be run once and is not tied to a gameobject
type SFX struct {
	id          string
	spritesheet *utility.Spritesheet
	v           pixel.Vec
	done        bool
	name        string
	manager     *Manager
	index       int
	smooth      bool
	frameRate   int
	frames      []Frame
}

//ID -- getter function for ID
func (s *SFX) ID() string {
	return s.id
}

//SetV -- getter function for V
func (s *SFX) SetV(v pixel.Vec) {
	s.v = v
}

//MappingToSFX -- helper function for creating sfx from config files
func MappingToSFX(spritesheet *utility.Spritesheet, mapping map[string]interface{}) []*SFX {
	var effects []*SFX
	for name, value := range mapping {
		attributes := value.(map[string]interface{})
		framesArr := attributes["Frames"].([]interface{})
		var sfxFrames []Frame
		for _, value := range framesArr {
			matrix := spritesheet.Matrix()
			frame := value.(float64)
			//does the sprite need to be reflected?
			if frame < 0 {
				frame = math.Abs(frame)
				matrix = matrix.Chained(pixel.IM.Rotated(pixel.ZV, -math.Pi).ScaledXY(pixel.ZV, pixel.V(-1, 1)).Rotated(pixel.ZV, math.Pi))
			}
			//assume that we do not use a custom matrix or color for effects created from spritesheets (maybe a bad guess?)
			sfxFrames = append(sfxFrames, NewFrame(int(frame), matrix, pixel.Alpha(1)))
		}
		smooth := attributes["Smooth"].(bool)
		frameRate := int(attributes["FrameRate"].(float64)) //cast to int because ain't dealing with floating point nonsense
		effects = append(effects, NewSFX(spritesheet, name, sfxFrames, smooth, frameRate))
	}
	return effects
}

//Frame -- a helper struct for creating effects which have colored masks or which use an adjusted matrix
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

//NewSFX -- constructor for SFX objects
func NewSFX(spritesheet *utility.Spritesheet, name string, frames []Frame, smooth bool, frameRate int) *SFX {
	if smooth {
		//smooth sfx by padding it with frames appended in reverse order
		for i := len(frames) - 1; i > 0; i-- {
			frames = append(frames, frames[i])
		}
	}
	id := xid.New().String()
	sfx := SFX{
		spritesheet: spritesheet,
		frameRate:   frameRate,
		name:        name,
		smooth:      smooth,
		id:          id,
		frames:      frames,
		done:        false,
	}
	return &sfx
}

//Clone -- factory method to deep clones the properties from a reference effect
func (s *SFX) Clone() *SFX {
	effect := NewSFX(s.spritesheet, s.name, s.frames, s.smooth, s.frameRate)
	effect.manager = s.manager
	return effect
}

//Next -- fetches the next animation or else kills this effect
func (s *SFX) Next(tick int) (*pixel.Sprite, Frame) {
	if tick%(60/s.frameRate) == 0 {
		s.index++
		if s.index >= len(s.frames)-1 {
			//kill self
			s.manager.killEffect(s.id)
			s.done = true
		}
	}
	frame := s.frames[s.index].Frame
	return s.spritesheet.Sprites()[frame], s.frames[s.index]
}
