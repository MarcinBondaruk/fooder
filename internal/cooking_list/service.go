package cooking_list

import "context"

type Service struct {
	repository CookingListRepository
}

func NewService(repository CookingListRepository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) CreateCookingList(ctx context.Context, recipeID int) (int, error) {
	cookingList := NewCookingList([]int{recipeID})

	id, err := s.repository.CreateCookingList(ctx, cookingList)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) AddRecipeToCookingList(ctx context.Context, cookingListID, recipeID int) error {
	cookingList, err := s.repository.GetCookingList(ctx, cookingListID)
	if err != nil {
		return err
	}

	cookingList.AddRecipe(recipeID)

	err = s.repository.UpdateCookingList(ctx, cookingList)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ViewCookingList(ctx context.Context, id int) (CookingList, error) {
	cookingList, err := s.repository.GetCookingList(ctx, id)
	if err != nil {
		return CookingList{}, err
	}

	return cookingList, nil
}
