<?php
namespace Recipe\Repository;

use PDO;
use PDOException;
use Recipe\Model\Recipe;

// [
//     "id" => 2,
//     "title" => "",
//     "thumbnail" => "",
//     "preparation" => "",
//     "ingredients" => [],
//     "tags" => [],
// ],

final class RecipeRepository
{
    public function __construct(
        private readonly PDO $pdo,
    ) {}

    public function add(Recipe $recipe): void
    {
        try {
            $stmt = $this->pdo->prepare('INSERT INTO recipes(title, preparation) VALUES (:title, :preparation)');
            $stmt->execute([
                'title' => $recipe->title(),
                'preparation' => $recipe->preparation(),
            ]);
        } catch(PDOException $e) {
            // TODO: log db error
            throw new FailedToAddRecipeException();
        }
    }

    public function retrieve(int $id): ?Recipe
    {
        try {
            $stmt = $this->pdo->prepare('SELECT * FROM recipes WHERE id = :id');
            $stmt->execute([
                'id' => $id,
            ]);
        } catch(PDOException $e) {
            // TODO: log db error
            return null;
        }

        $result = $stmt->fetch(PDO::FETCH_ASSOC);

        if (empty($result)) {
            return null;
        }

        return new Recipe(
            $result['id'],
            $result['title'],
            $result['preparation'],
            []
        );
    }
}