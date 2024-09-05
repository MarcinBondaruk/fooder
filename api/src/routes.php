<?php

use Symfony\Component\Routing\Route;
use Symfony\Component\Routing\RouteCollection;

$routes = new RouteCollection();
$routes->add('recipes_get', new Route(
    path: '/api/v1/recipes',
    defaults: [
        '_controller' => ['cookbook.ctrl', 'getRecipes'],
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