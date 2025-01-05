package cooking_list

type Repository interface {
	CreateCookingList() int
	AddRecipeToCookingList(cookingListID, recipeID int)
	ViewCookingList(cookingListID int) []int
}
