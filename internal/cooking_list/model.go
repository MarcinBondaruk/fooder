package cooking_list

type CookingList struct {
	ID      int
	recipes []int
}

func (cl *CookingList) addRecipe(recipeID int) {
	cl.recipes = append(cl.recipes, recipeID)
}
