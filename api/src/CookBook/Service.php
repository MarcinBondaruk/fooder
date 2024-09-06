<?php

namespace CookBook;

use CookBook\Response\GetRecipeResponse;

final class Service
{
    public function __construct(
        private Repository $repository,
    ) {}

    /**
     * @return GetRecipesResponse[]
     */
    public function getRecipes(): array
    {
        $recipes = $this->repository->getRecipes();

        return $recipes;
    }

    public function getRecipeById(int $id): GetRecipeResponse
    {
        return $this->repository->getRecipe($id);
    }
}