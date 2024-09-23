<?php

namespace Recipe\Dto;

final class RecipeDTO
{
    public function __construct(
        public readonly string $title,
        public readonly string $preparation,
    ) {}
}