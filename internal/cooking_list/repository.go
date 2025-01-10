package cooking_list

type Repository interface {
	createCookingList(cookingList CookingList) (int, error)
	updateCookingList(cookingList CookingList) error
	getCookingList(cookingListID int) (CookingList, error)
}
