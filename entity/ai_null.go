package entity

type NullAi struct {
	entity *Entity
}

//simple constructor
func NewNullAi(entity *Entity) ai {
	return &NullAi{
		entity: entity,
	}
}

//implementation of ai interface
func (a *NullAi) Update(tick int) {
	var (
		e  = a.entity
		am = e.AnimationManager()
	)
	if am.Ready() {
		am.Select("Idle")
	}
	e.Draw(tick)
}
