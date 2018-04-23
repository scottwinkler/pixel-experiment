package enum

type entityRegistry struct {
	Slime   string
	Spawner string
	Player  string
}

//Entity -- an enum
var Entity = newEntityRegistry()

func newEntityRegistry() *entityRegistry {
	return &entityRegistry{
		Slime:   "slime",
		Spawner: "spawner",
		Player:  "player",
	}
}
