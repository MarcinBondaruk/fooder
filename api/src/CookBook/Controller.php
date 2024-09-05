<?php
namespace CookBook;

use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

final class Controller {
    public function __construct(
        private Service $recipeService,
    ) {}

    public function getRecipes(Request $request): Response {
        $recipes = $this->recipeService->getRecipes();

        $response = new Response(json_encode($recipes));
        $response->setStatusCode(200);
        $response->headers->set('Content-Type', 'application/json');
        return $response;
    }
}