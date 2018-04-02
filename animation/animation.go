package animation

import (
	"github.com/faiface/pixel"
	"github.com/scottwinkler/pixel-experiment/spritesheet"
)

//func(*FrameData) AddFrame(pic *pixel.Picture,frame)

type Animation struct {
	Name        string
	Spritesheet *spritesheet.Spritesheet
	Frames      []int
	//FrameRate    int
	AnimationManager *AnimationManager
	Index            int
	Loop             bool
	Paused           bool
	//Ticker       *ticker
}

func NewAnimation(spritesheet *spritesheet.Spritesheet, name string, frames []int, loop bool) *Animation {
	animation := Animation{
		Name:        name,
		Spritesheet: spritesheet,
		Frames:      frames,
		//	FrameRate:    frameRate,
		Index:  0,
		Loop:   loop,
		Paused: true,
	}
	return &animation
}

func (a *Animation) SetAnimationManager(animationManager *AnimationManager) {
	a.AnimationManager = animationManager
}

func (a *Animation) Reset() {
	a.Index = 0
	a.Paused = true

}
func (a *Animation) SetPaused(paused bool) {
	//fmt.Println("pausing animation")
	a.Paused = paused
	//sprite := a.Next()
	//return a.Next()
	//	sprite.Draw(target, matrix)
}

/*func (a *Animation) Play() *pixel.Sprite{
	//fmt.Println("playing animation")
	a.Paused = false
	sprite := a.Next()
	//sprite.Draw(target, matrix)
}*/

func (a *Animation) Next() *pixel.Sprite {
	//interval := int(float64(60) / float64(a.FrameRate))
	nextIndex := a.Index
	if !a.Paused {
		//if elapsedSeconds%interval == 0 { //time to change frames
		if a.Loop {
			nextIndex = (a.Index + 1) % len(a.Frames)
		}
		//}
	}
	a.Index = nextIndex
	frame := a.Frames[a.Index]
	return a.Spritesheet.Sprites[frame]
}
