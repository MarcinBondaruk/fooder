package recipe

type Recipe struct {
	id          int
	name        string
	description string
	ingredients []string
}

func NewRecipe(name string, description string, ingredients []string) Recipe {
	return Recipe{
		name:        name,
		description: description,
		ingredients: ingredients,
	}
}

func (r *Recipe) ID() int {
	return r.id
}

func (r *Recipe) Name() string {
	return r.name
}

func (r *Recipe) Description() string {
	return r.description
}

func (r *Recipe) Ingredients() []string {
	return r.ingredients
}
