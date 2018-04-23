package enum

type materialRegistry struct {
	Wood  string
	Metal string
	Flesh string
}

//Material -- an enum
var Material = newMaterialRegistry()

func newMaterialRegistry() *materialRegistry {
	return &materialRegistry{
		Wood:  "wood",
		Metal: "metal",
		Flesh: "flesh",
	}
}
