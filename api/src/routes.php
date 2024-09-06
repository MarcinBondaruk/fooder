<?php

use Symfony\Component\Routing\Route;
use Symfony\Component\Routing\RouteCollection;

$routes = new RouteCollection();
$routes->add('cookbook_get_recipes', new Route(
    path: '/api/v1/recipes',
    defaults: [
        '_controller' => ['cookbook.ctrl', 'getRecipes'],
    ],
    methods: ['GET'],
));
$routes->add('cookbook_get_recipe', new Route(
    path: '/api/v1/recipes/{id}',
    defaults: [
        '_controller' => ['cookbook.ctrl', 'getRecipeById'],
    ],
    methods: ['GET'],
));
$routes->add('grocery_list_create', new Route(
    path: '/api/v1/grocery-list',
    defaults: [
        '_controller' => ['groceries.ctrl', 'createGroceryList'],
    ],
    methods: ['POST'],
));

return $routes;