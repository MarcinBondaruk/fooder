<?php
namespace Groceries;

use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

final class Controller
{
    public function createGroceryList(Request $request): Response
    {
        return new Response('Here\'s your grocery list');
    }
}