package cooking_list

type Repository interface {
	CreateCookingList(recipeID int) int
	AddRecipeToCookingList(cookingListID, recipeID int)
	ViewCookingList(cookingListID int) []int
}
