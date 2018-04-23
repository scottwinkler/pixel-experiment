package gui

import (
	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"github.com/scottwinkler/simple-rpg/world"
)

//Manager -- the struct that manages guis
type Manager struct {
	player world.Player //the game object which owns this gui
	panels []*Panel
	id     string
}

//NewManager -- constructor for gui manager
func NewManager(player world.Player) *Manager {
	var panels []*Panel
	panel := NewPanel(100, 100, pixel.ZV, player)
	panels = append(panels, panel)
	id := xid.New().String()
	manager := &Manager{
		player: player,
		panels: panels,
		id:     id,
	}
	player.GameObject().World().AddGUI(manager) //register self with world
	return manager
}

//ID -- getter function for Id. Statisfies world.GUI interface
func (m *Manager) ID() string {
	return m.id
}

//Update -- getter function for update. Satisfies world.GUI interface
func (m *Manager) Update(tick int) {
	for _, panel := range m.panels {
		panel.Update(tick)
	}
}

//SetVisible -- setter function for visiblity
func (m *Manager) SetVisible(visible bool) {
	//note that this is ghetto and doesn't support multiple panels
	//but i am tired and just want to see it working...
	for _, panel := range m.panels {
		panel.visible = visible
	}
}
