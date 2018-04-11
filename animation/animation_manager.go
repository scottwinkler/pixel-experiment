package animation

import (
	"strings"
)

type AnimationManager struct {
	Animations []*Animation
	Selected   *Animation
}

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
func (am *AnimationManager) AddAnimation(animation *Animation) {
	animation.animationManager = am
	am.Animations = append(am.Animations, animation)
}

func (am *AnimationManager) Select(name string) {
	//only reset if its not a looping animation.
	if am.Selected != nil && !am.Selected.loop {
		am.Selected.Reset()
	}
	for _, animation := range am.Animations {
		if strings.EqualFold(animation.name, name) {
			am.Selected = animation
			am.Selected.SetPaused(false)
		}
	}
}
