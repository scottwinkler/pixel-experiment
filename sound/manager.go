package sound

import (
	"fmt"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

//Manager -- manager for sounds
type Manager struct {
	Sounds []*Sound
}

//NewManager -- constructor for manager
func NewManager(sounds []*Sound) *Manager {
	//having performance issues when calling speaker.Init() too many times. Instead we will
	//initialize the speaker once and assume all sound uses the same format
	speaker.Init(48000, 4800)
	var manager *Manager
	manager = &Manager{
		Sounds: sounds,
	}
	return manager
}

func (m *Manager) asyncPlay(names ...string) {
	var streamers []beep.Streamer
	for _, sound := range m.Sounds {
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

//Play -- plays a sound (or two or three...)
func (m *Manager) Play(names ...string) {
	go m.asyncPlay(names...) //run asynchronously
}

//PlayWithDelay -- a helper method for playing sounds some time in the future
func (m *Manager) PlayWithDelay(d time.Duration, names ...string) {
	go func() {
		time.Sleep(d)
		m.Play(names...)
	}()
}
