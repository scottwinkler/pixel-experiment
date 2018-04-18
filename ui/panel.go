package ui

type Panel struct {
 bounds pixel.Rect
 widgets []Widget
 window *pixel.Window
imdraw *imdraw.IMDraw //all children widgets draw to the batch for better performance
}

//constructor for panel
func NewPanel() *Panel {
	panel := &Panel{
		batch
	}
	return panel
}
//called on every update by the main game loop
func (p *Panel) Update(tick int){
//update each of the children widgets
 for _,widget := range p.widgets {
	 widget.Update(tick)
 }
}

//add a widget to the panel
func (p *Panel) AddWidget(w Widget){
	p.widgets = append(p.widgets,w)
}

//remove a widget from the panel
func (p *Panel) RemoveWidget(w Widget){
	for i,widget := range p.widgets {
		if strings.EqualFold(w.Id(),widget.Id()){
			p.widgets = append(p.widgets[:i],p.widgets[i+1:]...)
			return
		}
	}
}