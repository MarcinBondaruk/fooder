<?php

namespace Recipe;

use Recipe\Api\Response\GetRecipeResponse;
use Recipe\Dto\RecipeDTO;
use Recipe\Model\Recipe;
use Recipe\Repository\RecipeQueryRepository;
use Recipe\Repository\RecipeRepository;

final class CookbookService
{
    public function __construct(
        private RecipeRepository $recipeRepository,
        private RecipeQueryRepository $recipeQueryRepository,
    ) {}

    public function getRecipes(): array
    {
        $recipes = $this->recipeQueryRepository->getRecipes();

        return $recipes;
    }

    public function getRecipeById(int $id): GetRecipeResponse
    {
        return $this->recipeQueryRepository->getRecipe($id);
    }

    public function createRecipe(RecipeDTO $recipeDTO): void
    {
        $recipe = new Recipe(1, 'title', 'prep', []);
        $this->recipeRepository->add($recipe);
    }
}