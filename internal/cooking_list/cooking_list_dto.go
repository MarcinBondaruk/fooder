package cooking_list

type CreateCookingListRequest struct {
	RecipeID int `json:"recipeId"`
}

type AddRecipeToCookingListRequest struct {
	RecipeID int `json:"recipeId"`
}

type CookingListItem struct {
	RecipeID int `json:"recipeId"`
	//RecipeName string `json:"recipeName"`
}

type CookingListResponse struct {
	ID      int               `json:"id"`
	Recipes []CookingListItem `json:"recipes"`
}
