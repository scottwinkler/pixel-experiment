package ui

type Button struct {
	panel *Panel
	bounds pixel.Rect
	sprite *pixel.Sprite //picture of button. use IMDraw for additional details
	window *pixel.Window
}

func NewButton() *Button{
	button = &Button{

	}
	return button
}

//called on every tick of the main game loop
func (b *Button) Update(tick int){

	Draw()
}

//do something when the button is clicked
func (b *Button) HandleClick(){

}

//do something when the mouse is over this
func (b *Button) HandleMouseOver(){

}

func (b *Button) Draw(){

}