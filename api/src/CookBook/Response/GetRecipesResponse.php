<?php
namespace CookBook\Response;

final class GetRecipesResponse
{
    public function __construct(
        public int $id,
        public string $title,
    ) {}
}