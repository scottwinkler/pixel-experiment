package sfx

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"github.com/scottwinkler/simple-rpg/utility"
)

//an SFX is a simplified animation that only needs to be run once and is not tied to a gameobject
type SFX struct {
	id          string
	spritesheet *utility.Spritesheet
	v           pixel.Vec
	done        bool
	name        string
	sfxManager  *SFXManager
	index       int
	smooth      bool
	frameRate   int
	frames      []SFXFrame
}

func (s *SFX) Id() string {
	return s.id
}

func (s *SFX) SetV(v pixel.Vec) {
	s.v = v
}

func MappingToSFX(spritesheet *utility.Spritesheet, mapping map[string]interface{}) []*SFX {
	var effects []*SFX
	for name, value := range mapping {
		attributes := value.(map[string]interface{})
		framesArr := attributes["Frames"].([]interface{})
		var sfxFrames []SFXFrame
		for _, value := range framesArr {
			matrix := spritesheet.Matrix()
			frame := value.(float64)
			//does the sprite need to be reflected?
			if frame < 0 {
				frame = math.Abs(frame)
				matrix = matrix.Chained(pixel.IM.Rotated(pixel.ZV, -math.Pi).ScaledXY(pixel.ZV, pixel.V(-1, 1)).Rotated(pixel.ZV, math.Pi))
			}
			//assume that we do not use a custom matrix or color for effects created from spritesheets (maybe a bad guess?)
			sfxFrames = append(sfxFrames, NewSFXFrame(int(frame), matrix, pixel.Alpha(1)))
		}
		smooth := attributes["Smooth"].(bool)
		frameRate := int(attributes["FrameRate"].(float64)) //cast to int because ain't dealing with floating point nonsense
		effects = append(effects, NewSFX(spritesheet, name, sfxFrames, smooth, frameRate))
	}
	return effects
}

//a helper struct for creating effects which have colored masks or which use an adjusted matrix
type SFXFrame struct {
	Frame  int
	Matrix pixel.Matrix
	Mask   pixel.RGBA
}

//constructor utility of SFXFrame objects
func NewSFXFrame(frame int, matrix pixel.Matrix, mask pixel.RGBA) SFXFrame {
	return SFXFrame{
		Frame:  frame,
		Matrix: matrix,
		Mask:   mask,
	}
}

func NewSFX(spritesheet *utility.Spritesheet, name string, frames []SFXFrame, smooth bool, frameRate int) *SFX {
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

//factory method to deep clones the properties from a reference effect
func (s *SFX) Clone() *SFX {
	effect := NewSFX(s.spritesheet, s.name, s.frames, s.smooth, s.frameRate)
	effect.sfxManager = s.sfxManager
	return effect
}

func (s *SFX) Next(tick int) (*pixel.Sprite, SFXFrame) {
	if tick%(60/s.frameRate) == 0 {
		s.index++
		if s.index >= len(s.frames)-1 {
			//kill self
			s.sfxManager.killEffect(s.id)
			s.done = true
		}
	}
	frame := s.frames[s.index].Frame
	return s.spritesheet.Sprites()[frame], s.frames[s.index]
}
