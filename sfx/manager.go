package sfx

import (
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//Manager - world should be the only one that has an SFXManager
type Manager struct {
	Effects []*SFX //a library of effects to use
	Running []*SFX //the currently running effects managed by this object
	window  *pixelgl.Window
}

//NewManager -- constructor for manager
func NewManager(effects []*SFX, window *pixelgl.Window) *Manager {
	var running []*SFX
	manager := &Manager{
		Effects: effects,
		Running: running,
		window:  window,
	}
	//add reference to sfxManager for each reference effect
	for _, effect := range effects {
		effect.manager = manager
	}
	return manager
}

//helper method for getting an effect from the library
func (m *Manager) getEffectByName(name string) *SFX {
	var output *SFX
	for _, effect := range m.Effects {
		if effect.name == name {
			output = effect
			break
		}
	}
	return output
}

func (m *Manager) killEffect(id string) {
	for index, runningEffect := range m.Running {
		if strings.EqualFold(runningEffect.ID(), id) {
			//delete the offending element. We didnt like him anyways
			m.Running = append(m.Running[:index], m.Running[index+1:]...)
			break
		}
	}
}

//PlayEffect -- plays an effect from the library at the given point
func (m *Manager) PlayEffect(name string, v pixel.Vec) {
	referenceEffect := m.getEffectByName(name)
	newEffect := referenceEffect.Clone()
	newEffect.SetV(v)
	m.Running = append(m.Running, newEffect)
}

//PlayCustomEffect -- plays an effect that isn't cached in the library. useful for playing computed effects that only ever get run once
func (m *Manager) PlayCustomEffect(sfx *SFX, v pixel.Vec) {
	sfx.manager = m
	sfx.SetV(v)
	m.Running = append(m.Running, sfx)
}

//Update -- called by main game loop. Loop through all effects and update each one
func (m *Manager) Update(tick int) {
	for _, effect := range m.Running {
		sprite, sfxFrame := effect.Next(tick)
		if !effect.done {
			target := m.window
			matrix := sfxFrame.Matrix.Moved(effect.v)
			mask := sfxFrame.Mask
			//fmt.Printf(" sfxFrame.Frame: %v", sfxFrame.Frame)
			sprite.DrawColorMask(target, matrix, mask)
		}
	}
}
