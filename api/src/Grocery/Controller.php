<?php
namespace Grocery;

use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

final class Controller
{
    public function __construct(
        private GroceryService $groceryService
    ) {}

    public function createShoppingList(Request $request): Response
    {
        $groceryList = $this->groceryService->newShoppingList(1);

        $response = new Response(json_encode($groceryList));
        $response->setStatusCode(201);
        $response->headers->set('Content-Type', 'application/json');
        return $response;
    }
}