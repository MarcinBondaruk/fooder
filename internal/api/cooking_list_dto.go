package api

type CreateCookingListRequest struct {
	RecipeID int `json:"recipeId"`
}

type AddRecipeToCookingListRequest struct {
	RecipeID int `json:"recipeId"`
}
