package gui

import (
	"image/color"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"golang.org/x/image/colornames"
)

//Button the structure for button ui elements
type Button struct {
	id      string
	panel   *Panel        //the containing class
	sprite  *pixel.Sprite //picture of button. use IMDraw for additional details
	matrix  pixel.Matrix  //on second thought we need a matrix to allow scaling of pictures
	visible bool
}

//NewButton constructor for button
func NewButton(panel *Panel, sprite *pixel.Sprite, matrix pixel.Matrix) *Button {
	id := xid.New().String()
	button := &Button{
		id:      id,
		panel:   panel,
		visible: true,
		sprite:  sprite,
		matrix:  matrix,
	}
	return button
}

//ID returns unique id
func (b *Button) ID() string {
	return b.id
}

//Update called on every tick of the main game loop
func (b *Button) Update(tick int) {
	if b.visible && b.panel.visible {
		p := b.panel.player
		win := b.panel.player.GameObject().World().Window
		mouse := p.Camera().Matrix.Unproject(win.MousePosition())
		if b.containsPoint(mouse) {
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				b.HandleClick()
			} else {
				b.HandleMouseOver()
			}
		} else {
			b.Draw(nil)
		}
	}
}

//HandleClick do something when the button is clicked
func (b *Button) HandleClick() {
	b.panel.visible = false //simply toggle visiblity of the panel
}

//HandleMouseOver do something when the mouse is over this
func (b *Button) HandleMouseOver() {
	b.Draw(colornames.Green)
}

//helper function for mouse handling events
func (b *Button) containsPoint(v pixel.Vec) bool {
	//need to project based on matrix transformation,
	//then subtract the center, because that is the origin around which sprites get drawn
	//then add the logical offset of the camera
	offset := b.panel.player.Camera().V
	frame := b.sprite.Frame()
	center := b.sprite.Frame().Center()
	pos1 := b.matrix.Project(frame.Min.Sub(center)).Add(offset)
	pos2 := b.matrix.Project(frame.Max.Sub(center)).Add(offset)
	bounds := pixel.R(pos1.X, pos1.Y, pos2.X, pos2.Y)
	return bounds.Contains(v)
}

//Draw draws the button onto the window target
func (b *Button) Draw(mask color.Color) {
	t := b.panel.player.GameObject().World().Window
	offset := b.panel.player.Camera().V
	matrix := b.matrix.Moved(offset) //translate by panel position relative to player
	b.sprite.DrawColorMask(t, matrix, mask)
}
