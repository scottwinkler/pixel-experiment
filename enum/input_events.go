package enum

type mouseEventRegistry struct {
	Click     string
	MouseOver string
}

//MouseEvent -- an enum
var MouseEvent = newMouseEventRegistry()

func newMouseEventRegistry() *mouseEventRegistry {
	return &mouseEventRegistry{
		Click:     "Click",
		MouseOver: "MouseEvent",
	}
}
