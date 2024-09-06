<?php
namespace CookBook\Response;

final class GetRecipeResponse
{
    public function __construct(
        public int $id,
        public string $title,
        public string $thumbnail,
        public string $preparation,
        public array $ingredients,
        public array $tags,
    ) {}
}