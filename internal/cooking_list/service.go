package cooking_list

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) CreateCookingList(recipeID int) int {
	cookingList := CookingList{recipes: []int{recipeID}}

	return s.repository.createCookingList(cookingList)
}

func (s *Service) AddRecipeToCookingList(cookingListID, recipeID int) error {
	cookingList, err := s.repository.getCookingList(cookingListID)
	if err != nil {
		return err
	}

	cookingList.addRecipe(recipeID)

	err = s.repository.updateCookingList(cookingList)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ViewCookingList(id int) (CookingList, error) {
	cookingList, err := s.repository.getCookingList(id)
	if err != nil {
		return CookingList{}, err
	}
	return cookingList, nil
}
