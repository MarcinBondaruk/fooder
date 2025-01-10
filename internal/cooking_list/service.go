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
	return s.repository.CreateCookingList(recipeID)
}

func (s *Service) AddRecipeToCookingList(cookingListID, recipeID int) error {
	// get cooking list
	// add recipe to cooking list
	// update cooking list
	s.repository.AddRecipeToCookingList(cookingListID, recipeID)

	return nil
}

func (s *Service) ViewCookingList(id int) []int {
	return s.repository.ViewCookingList(id)
}
