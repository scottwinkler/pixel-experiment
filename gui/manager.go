package gui

import (
	"strings"

	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"github.com/scottwinkler/simple-rpg/world"
)

//Manager -- the struct that manages guis
type Manager struct {
	player    world.Player //the game object which owns this gui
	panels    []*Panel
	id        string
	listeners map[string][]Widget //used for input callbacks
}

//NewManager -- constructor for gui manager
func NewManager(player world.Player) *Manager {
	var panels []*Panel
	panel := NewPanel(100, 100, pixel.ZV, player)
	panels = append(panels, panel)
	id := xid.New().String()
	listeners := map[string][]Widget{
		"MouseOver": []Widget{},
		"Click":     []Widget{},
	}
	manager := &Manager{
		player:    player,
		panels:    panels,
		id:        id,
		listeners: listeners,
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

//registers a widget with a particular event
func (m *Manager) addEventListener(event string, widget Widget) {
	m.listeners[event] = append(m.listeners[event], widget)
}

//unregisters a widget with a particular event
func (m *Manager) deleteEventListener(event string, widget Widget) {
	widgets := m.listeners[event]
	for i, obj := range widgets {
		if strings.EqualFold(obj.ID(), widget.ID()) {
			widgets = append(widgets[:i], widgets[i+1:]...)
			m.listeners[event] = widgets
			break
		}
	}
}
