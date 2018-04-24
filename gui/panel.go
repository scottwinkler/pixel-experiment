package gui

import (
	"strings"

	"github.com/scottwinkler/simple-rpg/utility"
	"github.com/scottwinkler/simple-rpg/world"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

//Panel the struct that acts as a container for all other widgets
type Panel struct { //panel should be a special kind of widget that holds other widgets
	//widget  Widget
	height  float64
	width   float64
	v       pixel.Vec
	widgets []Widget //child widgets know how to draw themselves
	player  world.Player
	imd     *imdraw.IMDraw
	visible bool
}

//NewPanel constructor for panel
func NewPanel(height float64, width float64, v pixel.Vec, player world.Player) *Panel {
	panel := &Panel{
		imd:     imdraw.New(nil),
		height:  height,
		width:   width,
		v:       v,
		player:  player,
		visible: true, //default!!!
	}
	//register widgets with panel, then add to panel
	pic, _ := utility.LoadPicture("_assets/ui/close_normal.png")
	sprite := pixel.NewSprite(pic, pic.Bounds())
	buttonMatrix := pixel.IM.Scaled(pixel.ZV, 0.05).Moved(pixel.V(165, 165))
	var testWidgets []Widget
	testButton := NewButton(panel, sprite, buttonMatrix)
	testWidgets = append(testWidgets, testButton)
	panel.widgets = testWidgets
	return panel
}

//SetV setter method for v
func (p *Panel) SetV(v pixel.Vec) {
	p.v = v
}

//Update called on every update by the main game loop
func (p *Panel) Update(tick int) {
	playerPos := p.player.GameObject().V()
	p.v = playerPos.Add(pixel.V(75, 75)) //draw relative to player position
	p.Draw()
	//update each of the children widgets
	for _, widget := range p.widgets {
		widget.Update(tick)
	}
}

//AddWidget add a widget to the panel
func (p *Panel) AddWidget(w Widget) {
	p.widgets = append(p.widgets, w)
}

//RemoveWidget removes a widget from the panel
func (p *Panel) RemoveWidget(w Widget) {
	for i, widget := range p.widgets {
		if strings.EqualFold(w.ID(), widget.ID()) {
			p.widgets = append(p.widgets[:i], p.widgets[i+1:]...)
			return
		}
	}
}

//Draw draws the panel onto the window target
func (p *Panel) Draw() {

	if p.visible {
		minX := p.v.X
		minY := p.v.Y
		maxX := p.v.X + p.width
		maxY := p.v.Y + p.height
		imd := p.imd
		imd.Clear()
		//draw background
		imd.Color = colornames.Darkgray
		imd.Push(pixel.V(minX, minY))
		imd.Push(pixel.V(maxX, maxY))
		imd.Rectangle(0)
		//draw border
		imd.Color = colornames.Black
		imd.Push(pixel.V(minX, minY))
		imd.Push(pixel.V(maxX, maxY))
		imd.Rectangle(1)
		imd.Draw(p.player.GameObject().World().Window)
	}
}
