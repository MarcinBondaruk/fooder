<?php
namespace Controller;

use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

final class Recipe {
    public function getRecipes(Request $request): Response {
        $recipes = [
            [
                'id' => 1,
                'title' => 'Eggs and Bacon',
                'preparation' => 'Lorem Ipsum Dolor...',
                'ingredients' => ['eggs', 'bacon', 'butter'],
            ],
            [
                'id' => 2,
                'title' => 'Spaghetti Pomodoro',
                'preparation' => 'Lorem Ipsum Dolor...',
                'ingredients' => ['tomato', 'saghetti pasta', 'olive oil'],
            ],
        ];

        $response = new Response(json_encode($recipes));
        $response->setStatusCode(200);
        $response->headers->set('Content-Type', 'application/json');
        return $response;
    }
}