package animation

import (
	"strings"

	"github.com/faiface/pixel"
)

type AnimationManager struct {
	Animations []*Animation
	selected   *Animation //the currently playing animation
}

//constructor for animation manager
func NewAnimationManager(animations []*Animation) *AnimationManager {
	var animationManager *AnimationManager
	animationManager = &AnimationManager{
		Animations: animations,
	}

	//add reference to animation manager for each animation
	for _, animation := range animations {
		animation.animationManager = animationManager
	}
	return animationManager
}

//wrapper around internal state.
func (am *AnimationManager) Next(tick int) (*pixel.Sprite, AnimationFrame) {
	return am.selected.Next(tick)
}

//wrapper around internal state.
func (am *AnimationManager) Current() (*pixel.Sprite, AnimationFrame) {
	return am.selected.Current()
}

//getter for selected. Ask yourself if you really need this before using it. Check yourself before you wreck yourself
func (am *AnimationManager) Selected() *Animation {
	return am.selected
}

//a helper method to know if the current animation is ready to accept new input
//i.e current animation is done or can be skipped
func (am *AnimationManager) Ready() bool {
	return am.selected.skippable || am.selected.done
}

func (am *AnimationManager) Select(name string) {
	//only reset if its not a looping animation.
	if am.selected != nil && !am.selected.loop {
		am.selected.Reset()
	}
	for _, animation := range am.Animations {
		if strings.EqualFold(animation.name, name) {
			am.selected = animation
			am.selected.SetPaused(false)
		}
	}
}
