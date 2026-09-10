package ingredient

type Unit string

const (
	UnitGram       Unit = "g"
	UnitKilogram   Unit = "kg"
	UnitMilliliter Unit = "ml"
	UnitLiter      Unit = "l"
	UnitPiece      Unit = "piece"
	UnitTablespoon Unit = "tbsp"
	UnitTeaspoon   Unit = "tsp"
)

var validUnits = map[Unit]bool{
	UnitGram:       true,
	UnitKilogram:   true,
	UnitMilliliter: true,
	UnitLiter:      true,
	UnitPiece:      true,
	UnitTablespoon: true,
	UnitTeaspoon:   true,
}

func (u Unit) Valid() bool {
	return validUnits[u]
}

func (u Unit) String() string {
	return string(u)
}

type Ingredient struct {
	ID   int
	Name string
}

func NewIngredient(name string) Ingredient {
	return Ingredient{Name: name}
}
