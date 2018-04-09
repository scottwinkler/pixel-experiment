package sound

import (
	"strings"
	"time"

	"github.com/faiface/beep/speaker"
)

type SoundManager struct {
	Sounds []*Sound
}

func (sm *SoundManager) Play(name string) {
	for _, sound := range sm.Sounds {
		if strings.EqualFold(sound.Name, name) {
			streamer, format := sound.Decode()
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			speaker.Play(streamer)
		}
	}
}

func NewSoundManager(sounds []*Sound) *SoundManager {
	var soundManager *SoundManager
	soundManager = &SoundManager{
		Sounds: sounds,
	}
	return soundManager
}
