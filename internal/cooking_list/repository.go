package cooking_list

import "context"

type CookingListRepository interface {
	CreateCookingList(ctx context.Context, cookingList CookingList) (int, error)
	UpdateCookingList(ctx context.Context, cookingList CookingList) error
	GetCookingList(ctx context.Context, cookingListID int) (CookingList, error)
}
