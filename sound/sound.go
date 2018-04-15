package sound

import (
	"fmt"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"github.com/scottwinkler/simple-rpg/utility"
)

type Sound struct {
	name   string
	buffer *beep.Buffer
}

//utility function for converting a spritesheet based on a mapping of name:frames to an array of animations
func MappingToSounds(mapping map[string]interface{}) []*Sound {
	var sounds []*Sound
	for _, value := range mapping {
		arr := value.([]interface{})
		for _, obj := range arr {
			attributes := obj.(map[string]interface{})
			name := attributes["Name"].(string)
			path := attributes["Path"].(string)
			streamer, format := decode(path)
			buffer := beep.NewBuffer(format)
			buffer.Append(streamer)
			sounds = append(sounds, NewSound(name, buffer))
		}
	}
	return sounds
}

//private method for reading sound from files
func decode(path string) (beep.StreamSeekCloser, beep.Format) {
	file := utility.LoadFile(path)
	var streamer beep.StreamSeekCloser
	var format beep.Format
	parts := strings.Split(path, ".")
	extension := parts[len(parts)-1]
	var err error
	switch extension {
	case "wav":
		streamer, format, err = wav.Decode(file)
	case "flac":
		streamer, format, err = flac.Decode(file)
	case "mp3":
		streamer, format, err = mp3.Decode(file)
	}
	if err != nil {
		fmt.Printf("[ERROR]: Could not decode sound file: %v", err)
	}
	return streamer, format
}

//sound constructor
func NewSound(name string, buffer *beep.Buffer) *Sound {
	sound := Sound{
		name:   name,
		buffer: buffer,
	}
	return &sound
}
