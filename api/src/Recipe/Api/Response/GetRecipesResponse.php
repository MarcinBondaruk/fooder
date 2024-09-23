<?php
namespace Recipe\Api\Response;

final class GetRecipesResponse
{
    public function __construct(
        public int $id,
        public string $title,
        public string $thumbnail,
        public array $tags,
    ) {}
}