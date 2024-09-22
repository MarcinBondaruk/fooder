<?php
namespace Recipe;

final class Recipe
{
    public function __construct(
        private int $id,
        private string $title,
        private string $preparation,
        private array $ingredients,
    ) {}
}