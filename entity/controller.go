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

//Controller -- the interface which all controllers must implement
type Controller interface {
	Update(int)
	HitCallback(interface{}) bool
	AttackCallback(interface{})
}

//Constructor takes an entity and returns an instance of a controller interface
type Constructor func(*Entity) Controller
