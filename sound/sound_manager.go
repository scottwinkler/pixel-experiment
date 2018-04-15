package sound

import (
	"fmt"
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
			if strings.EqualFold(sound.name, name) {
				buffer := sound.buffer
				streamer := buffer.Streamer(0, buffer.Len())
				streamers = append(streamers, streamer)
			}
		}
	}
	if len(streamers) == 0 {
		fmt.Printf("[WARN] -- SoundManager could not find sound(s) to play: %v", names)
	}
	speaker.Play(streamers...)
}

func (sm *SoundManager) Play(names ...string) {
	go sm.asyncPlay(names...) //run asynchronously
}

//a helper method for playing sounds some time in the future
func (sm *SoundManager) DelayedPlay(d time.Duration, names ...string) {
	go func() {
		time.Sleep(d)
		sm.Play(names...)
	}()
}

func NewSoundManager(sounds []*Sound) *SoundManager {
	//having performance issues when calling speaker.Init() too many times. Instead we will
	//initialize the speaker once and assume all sound uses the same format
	speaker.Init(48000, 4800)
	var soundManager *SoundManager
	soundManager = &SoundManager{
		Sounds: sounds,
	}
	return soundManager
}
