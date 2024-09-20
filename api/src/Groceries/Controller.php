<?php
namespace Groceries;

use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

final class Controller
{
    public function createGroceryList(Request $request): Response
    {
        // wykonac serwis, zwrocic response
        // $groceryList = $this->groceryService->newShoppingList();
        $response = 'GroceryList';
        $response = 'asdas'
        return new Response($response);
    }
}