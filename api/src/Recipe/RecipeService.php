<?php

namespace Recipe;

use Recipe\Response\GetRecipeResponse;

final class CookbookService
{
    public function __construct(
        private Repository $repository,
    ) {}

    public function getRecipes(): array
    {
        $recipes = $this->repository->getRecipes();

        return $recipes;
    }

    public function getRecipeById(int $id): GetRecipeResponse
    {
        return $this->repository->getRecipe($id);
    }

    public function addToCookingList(int $recipeId): void
    {

    }
}