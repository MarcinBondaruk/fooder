<?php
namespace CookBook;

use CookBook\Views\PreviewRecipe;

final class Repository
{
    /**
     * @return PreviewRecipe[]
     */
    public function getRecipesPreview(): array
    {
        return [
            new PreviewRecipe(
                1,
                'Eggs\'n\'baky'
            ),
            new PreviewRecipe(
                2,
                'Spaghetti Pomodoro'
            ),
        ];
    }
}