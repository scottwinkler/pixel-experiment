package gui

//Widget the interface for all gui elements
type Widget interface {
	ID() string
	Update(int)
}
