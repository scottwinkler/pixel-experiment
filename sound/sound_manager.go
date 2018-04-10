package sound

import (
	"math/rand"
	"strings"
	"time"

	"github.com/faiface/beep/speaker"
)

type SoundManager struct {
	Sounds []*Sound
}

func (sm *SoundManager) Play(name string) {
	for _, group := range sm.Sounds {
		for _, sound := range group {
			if strings.EqualFold(sound.Name, name) {
				streamer, format := sound.Decode()
				speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
				speaker.Play(streamer)
			}
		}
	}
}

func (sm *SoundManager) PlayRandom(group string) {
	soundsGroup := sm.Sounds[group]
	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(soundsGroup) - 1)
	name := soundsGroup[index].Name
	sm.Play(name)
}

func NewSoundManager(sounds []*Sound) *SoundManager {
	var soundManager *SoundManager
	soundManager = &SoundManager{
		Sounds: sounds,
	}
	return soundManager
}
