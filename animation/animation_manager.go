package animation

import (
	"strings"
)

type AnimationManager struct {
	Animations []*Animation
	Selected   *Animation
	//Player     *player.Player
}

func NewAnimationManager(animations []*Animation) *AnimationManager {
	//select first animation as default
	var animationManager *AnimationManager
	/*if len(animations) > 0 {
		selected := animations[0]
		animationManager = &AnimationManager{
			Animations: animations,
			Selected:   selected,
			//Player:     player,
		}
	} else {*/
	animationManager = &AnimationManager{
		Animations: animations,
		//	Player:     player,
	}
	//}
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
	//if am.Current.Name
	for _, animation := range am.Animations {
		if strings.EqualFold(animation.Name, name) {
			am.Selected = animation
			//am.Current.Play(player)
		}
	}
}
