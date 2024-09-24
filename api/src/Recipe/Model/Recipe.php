<?php
namespace Recipe\Model;

final class Recipe
{
    public function __construct(
        private ?int $id,
        private string $title,
        private string $preparation,
        private array $ingredients,
    ) {}

    public function title(): string
    {
        return $this->title;
    }

    public function preparation(): string
    {
        return $this->preparation;
    }
}