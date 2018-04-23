package enum

type directionRegistry struct {
	Left  int
	Right int
	Down  int
	Up    int
}

//Direction -- an enum
var Direction = newDirectionRegistry()

func newDirectionRegistry() *directionRegistry {
	return &directionRegistry{
		Left:  0,
		Right: 1,
		Down:  2,
		Up:    3,
	}
}
