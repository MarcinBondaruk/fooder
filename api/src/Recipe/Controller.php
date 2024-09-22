<?php
namespace Recipe;

use Symfony\Component\HttpFoundation\Response;

final class Controller
{
    public function __construct(
        private CookbookService $recipeService,
    ) {}

    public function getRecipes(): Response {
        $recipes = $this->recipeService->getRecipes();

        $response = new Response(json_encode($recipes));
        $response->setStatusCode(200);
        $response->headers->set('Content-Type', 'application/json');
        return $response;
    }

    public function getRecipeById(int $id): Response {
        $recipe = $this->recipeService->getRecipeById($id);

        $response = new Response(json_encode($recipe));
        $response->setStatusCode(200);
        $response->headers->set('Content-Type', 'application/json');
        return $response;
    }
}