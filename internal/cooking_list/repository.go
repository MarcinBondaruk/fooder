package cooking_list

import "context"

type Repository interface {
	createCookingList(ctx context.Context, cookingList CookingList) (int, error)
	updateCookingList(ctx context.Context, cookingList CookingList) error
	getCookingList(ctx context.Context, cookingListID int) (CookingList, error)
}
