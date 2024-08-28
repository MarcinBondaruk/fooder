<?php
namespace CookBook\Views;

final class PreviewRecipe
{
    public function __construct(
        public int $id,
        public string $title,
    ) {}
}