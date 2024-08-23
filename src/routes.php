<?php

use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Route;
use Symfony\Component\Routing\RouteCollection;

$routes = new RouteCollection();
$routes->add('hello', new Route(
    path: '/hello/{name}',
    defaults: ['name' => 'Lock'],
    methods: ['GET']
));
$routes->add('goodbye', new Route(
    path: '/goodbye/{name}',
    defaults: [
        'name' => 'Sawyer',
        '_controller' => function (Request $request): Response {
            $request->attributes->set('foo', 'bar');

            $contents = renderTemplate($request);

            $response = new Response($contents);
            $response->headers->set('Content-Type', 'text/plain');

            return $response;
        }
    ],
    methods: ['GET'],
));
$routes->add('leap_year', new Route(
    path: '/leap-year/{year}',
    defaults: [
        'year' => 2000,
        '_controller' => 'Calendar\Controller\LeapYearController::getAnswer',
    ],
    methods: ['GET']
));
$routes->add('recipes_get', new Route(
    path: '/api/v1/recipes',
    defaults: [
        '_controller' => 'Controller\Recipe::getRecipes',
    ],
    methods: ['GET'],
));

return $routes;