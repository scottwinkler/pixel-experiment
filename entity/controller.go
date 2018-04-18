package entity

//lookup table for name to controller constructor
var (
	ControllerMap = map[string]Constructor{
		"slime":   NewSlimeController,
		"spawner": NewNullController,
		"null":    NewNullController,
		"player":  NewPlayerController,
	}
)

//the interface which all controllers must implement
type controller interface {
	Update(int)
	HitCallback(interface{}) bool
	AttackCallback(interface{})
}

//ai constructor takes an entity and returns an instance of a controller interface
type Constructor func(*Entity) controller
