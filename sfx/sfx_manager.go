package sfx

import (
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//world should be the only one that has an SFXManager
type SFXManager struct {
	Effects []*SFX //a library of effects to use
	Running []*SFX //the currently running effects managed by this object
	window  *pixelgl.Window
}

func NewSFXManager(effects []*SFX, window *pixelgl.Window) *SFXManager {
	var running []*SFX
	sfxManager := &SFXManager{
		Effects: effects,
		Running: running,
		window:  window,
	}
	//add reference to sfxManager for each reference effect
	for _, effect := range effects {
		effect.sfxManager = sfxManager
	}
	return sfxManager
}

//helper method for getting an effect from the library
func (sm *SFXManager) getEffectByName(name string) *SFX {
	var output *SFX
	for _, effect := range sm.Effects {
		if effect.name == name {
			output = effect
			break
		}
	}
	return output
}

func (sm *SFXManager) killEffect(id string) {
	//fmt.Printf("len %d", len(sm.Running))
	for index, runningEffect := range sm.Running {
		if strings.EqualFold(runningEffect.Id(), id) {
			//delete the offending element. We didnt like him anyways
			sm.Running = append(sm.Running[:index], sm.Running[index+1:]...)
			break
		}
	}
}

//plays an effect at the given point
func (sm *SFXManager) MakeEffect(name string, v pixel.Vec) {
	referenceEffect := sm.getEffectByName(name)
	newEffect := referenceEffect.Clone()
	newEffect.SetV(v)
	//fmt.Printf("appending effect %s", newEffect.name)
	sm.Running = append(sm.Running, newEffect)
}

//loop through all effects and update each one
func (sm *SFXManager) Update(tick int) {
	for _, effect := range sm.Running {
		//fmt.Println("calling next on effect")
		sprite := effect.Next(tick)
		sprite.Draw(sm.window, effect.Matrix())
	}
}
