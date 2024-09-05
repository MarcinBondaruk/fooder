<?php

namespace CookBook;

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
}