package sound

import (
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"github.com/scottwinkler/pixel-experiment/utility"
)

type Sound struct {
	Name      string
	Path      string
	Extension string
}

func (s *Sound) Decode() (beep.StreamSeekCloser, beep.Format) {
	file := utility.LoadFile(s.Path)
	var streamer beep.StreamSeekCloser
	var format beep.Format
	switch s.Extension {
	case "wav":
		streamer, format, _ = wav.Decode(file)
	case "flac":
		streamer, format, _ = flac.Decode(file)
	case "mp3":
		streamer, format, _ = mp3.Decode(file)
	}
	return streamer, format
}

//create a new sound from a file
func NewSound(name string, path string) *Sound {
	parts := strings.Split(path, ".")
	extension := parts[len(parts)-1]

	sound := Sound{
		Name:      name,
		Path:      path,
		Extension: extension,
	}
	return &sound
}
