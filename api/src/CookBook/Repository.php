<?php
namespace CookBook;

use CookBook\Response\GetRecipesResponse;

final class Repository
{
    /**
     * @return GetRecipesResponse[]
     */
    public function getRecipes(): array
    {
        return [
            new GetRecipesResponse(
                1,
                'Eggs\'n\'baky'
            ),
            new GetRecipesResponse(
                2,
                'Spaghetti Pomodoro'
            ),
        ];
    }
}