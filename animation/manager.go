package animation

import (
	"strings"

	"github.com/faiface/pixel"
)

//Manager -- manager for animations
type Manager struct {
	Animations []*Animation
	selected   *Animation //the currently playing animation
}

//NewManager constructor for animation manager
func NewManager(animations []*Animation) *Manager {
	var manager *Manager
	manager = &Manager{
		Animations: animations,
	}

	//add reference to animation manager for each animation
	for _, animation := range animations {
		animation.manager = manager
	}
	return manager
}

//Next -- wrapper around internal state.
func (m *Manager) Next(tick int) (*pixel.Sprite, Frame) {
	return m.selected.Next(tick)
}

//Current -- wrapper around internal state.
func (m *Manager) Current() (*pixel.Sprite, Frame) {
	return m.selected.Current()
}

//Selected -- getter for selected. Ask yourself if you really need this before using it. Check yourself before you wreck yourself
func (m *Manager) Selected() *Animation {
	return m.selected
}

//Ready -- a helper method to know if the current animation is ready to accept new input
//i.e current animation is done or can be skipped
func (m *Manager) Ready() bool {
	return m.selected.skippable || m.selected.done
}

//Select -- selects a new animation by name
func (m *Manager) Select(name string) {
	//only reset if its not a looping animation.
	if m.selected != nil && !m.selected.loop {
		m.selected.Reset()
	}
	for _, animation := range m.Animations {
		if strings.EqualFold(animation.name, name) {
			m.selected = animation
			m.selected.SetPaused(false)
		}
	}
}
