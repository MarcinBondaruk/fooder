package cooking_list

import "sync"

type InMemoryRepository struct {
	lock         sync.Mutex
	cookingLists map[int][]int
	nextID       int
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		cookingLists: make(map[int][]int),
	}
}

func (db *InMemoryRepository) CreateCookingList() int {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.nextID++
	db.cookingLists[db.nextID] = make([]int, 0)

	return db.nextID
}

func (db *InMemoryRepository) AddRecipeToCookingList(cookingListID, recipeID int) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.cookingLists[cookingListID] = append(db.cookingLists[cookingListID], recipeID)
}

func (db *InMemoryRepository) ViewCookingList(cookingListID int) []int {
	cookingList, ok := db.cookingLists[cookingListID]
	if !ok {
		return []int{}
	}

	return cookingList
}
