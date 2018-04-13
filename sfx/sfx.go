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
	name        string
	frames      []int
	sfxManager  *SFXManager
	index       int
	smooth      bool
	frameRate   int
	matrix      pixel.Matrix
}

func (s *SFX) Id() string {
	return s.id
}

func (s *SFX) SetV(v pixel.Vec) {
	s.v = v
}

func (s *SFX) Matrix() pixel.Matrix {
	return s.matrix
}
func MappingToSFX(spritesheet *utility.Spritesheet, mapping map[string]interface{}) []*SFX {
	var effects []*SFX
	for name, value := range mapping {
		attributes := value.(map[string]interface{})
		framesArr := attributes["Frames"].([]interface{})
		var frames []int
		for _, frame := range framesArr {
			frames = append(frames, int(frame.(float64)))
		}
		smooth := attributes["Smooth"].(bool)
		frameRate := int(attributes["FrameRate"].(float64)) //cast to int because ain't dealing with floating point nonsense
		effects = append(effects, NewSFX(spritesheet, name, frames, smooth, frameRate))
	}
	return effects
}

func NewSFX(spritesheet *utility.Spritesheet, name string, frames []int, smooth bool, frameRate int) *SFX {
	if smooth {
		//smooth sfx by padding it with frames appended in reverse order
		for i := len(frames) - 1; i > 0; i-- {
			frames = append(frames, frames[i])
		}
	}
	id := xid.New().String()
	sfx := SFX{
		spritesheet: spritesheet,
		frames:      frames,
		frameRate:   frameRate,
		name:        name,
		smooth:      smooth,
		id:          id,
		matrix:      spritesheet.Matrix(),
	}
	return &sfx
}

//factory method to deep clones the properties from a reference effect
func (s *SFX) Clone() *SFX {
	effect := NewSFX(s.spritesheet, s.name, s.frames, s.smooth, s.frameRate)
	effect.sfxManager = s.sfxManager
	return effect
}

func (s *SFX) Next(tick int) *pixel.Sprite {
	var frame int
	if tick%(60/s.frameRate) == 0 {
		s.index++
		if s.index >= len(s.frames)-1 {
			//kill self
			s.sfxManager.killEffect(s.id)
		}
		frame = s.frames[s.index]
	} else {
		//always return same frame if paused or not an appropriate time to change animations
		frame = s.frames[s.index]
	}

	//does the sprite need to be reflected?
	if frame < 0 {
		frame = int(math.Abs(float64(frame)))
		s.matrix = s.spritesheet.Matrix().Chained(pixel.IM.Rotated(pixel.ZV, -math.Pi).ScaledXY(pixel.ZV, pixel.V(-1, 1)).Rotated(pixel.ZV, math.Pi)).Chained(pixel.IM.Moved(s.v))
	} else {
		s.matrix = s.spritesheet.Matrix().Chained(pixel.IM.Moved(s.v))
	}
	return s.spritesheet.Sprites()[frame]
}

func (s *SFX) Draw(t *pixel.Target) {

}
