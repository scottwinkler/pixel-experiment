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
		animation.AnimationManager = animationManager
	}
	return animationManager
}
func (am *AnimationManager) AddAnimation(animation *Animation) {
	animation.AnimationManager = am
	am.Animations = append(am.Animations, animation)
}

func (am *AnimationManager) Select(name string) {
	//no need to save state if going to a new animation
	if am.Selected != nil && am.Selected.Name != name {
		am.Selected.Reset()
	}
	for _, animation := range am.Animations {
		if strings.EqualFold(animation.Name, name) {
			am.Selected = animation
			am.Selected.SetPaused(false)
		}
	}
}
