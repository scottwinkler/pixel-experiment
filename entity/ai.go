package entity

//lookup table for name to ai construcot
var (
	AiMap = map[string]Constructor{
		"slime":   NewSlimeAi,
		"spawner": NewNullAi,
		"null":    NewNullAi,
	}
)

//the interface which all ai must implement
type ai interface {
	Update(int)
}

//ai constructor takes an entity and returns an instance of an ai interface
type Constructor func(*Entity) ai
