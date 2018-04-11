package sound

import (
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type SoundManager struct {
	Sounds []*Sound
}

func (sm *SoundManager) asyncPlay(names ...string) {
	var streamers []beep.Streamer
	for _, sound := range sm.Sounds {
		for _, name := range names {
			if strings.EqualFold(sound.Name, name) {
				buffer := sound.Buffer
				streamer := buffer.Streamer(0, buffer.Len())
				format := buffer.Format()
				streamers = append(streamers, streamer)
				//this could cause problems if sampling rates are same for both
				speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			}
		}
	}
	speaker.Play(streamers...)
	//fmt.Println("finished playing")
}

func (sm *SoundManager) Play(names ...string) {
	go sm.asyncPlay(names...) //run asynchronously
}

func NewSoundManager(sounds []*Sound) *SoundManager {
	var soundManager *SoundManager
	soundManager = &SoundManager{
		Sounds: sounds,
	}
	return soundManager
}
