package cooking_list

type CookingList struct {
	ID      int
	Recipes []int
}

func NewCookingList(recipes []int) CookingList {
	return CookingList{
		Recipes: recipes,
	}
}

func (cl *CookingList) AddRecipe(recipeID int) {
	cl.Recipes = append(cl.Recipes, recipeID)
}
