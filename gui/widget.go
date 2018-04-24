package gui

//Widget the interface for all gui elements. should register interest in events. maybe make a struct that hanldes
//the registration process?
type Widget interface {
	ID() string
	Update(int)
	HandleClick()
	//HandleDoubleClick() //via repeated events?
	HandleMouseOver()
	//HandleDragStart()
	//HandleDragEnd()
}
