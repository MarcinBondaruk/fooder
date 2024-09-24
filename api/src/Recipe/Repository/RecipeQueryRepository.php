<?php

namespace Recipe\Repository;

use PDO;
use Recipe\Api\Response\GetRecipeResponse;
use Recipe\Api\Response\GetRecipesResponse;

const RECIPES = [
    [
        "id" => 1,
        "title" => "SMOOTHIE BOWL",
        "thumbnail" => "/some/img/url.img",
        "preparation" => "
            1. Obrań pomarańcze, wybrać pestki. Obrać banana.
            2. Zblendować pomarańczę, banana, jogurt i sok/mleko.
            3. Włożyć do miseczki, ułożyć malint i borówki, posypać pestkami i płatkami owsianymi, skropić syropem klonowym.",
        "ingredients" => [
            "banan",
            "pomarańcza",
            "borówki",
            "maliny",
            "jogurt naturalny",
            "sok ananasowy",
            "pestki słonecznika",
            "płatki owsiane",
            "syrop klonowy"
        ],
        "tags" => [
            "śniadanie",
            "vege"
        ],
    ],
    [
        "id" => 2,
        "title" => "GRANOLA Z OWOCAMI",
        "thumbnail" => "",
        "preparation" => "Lorem ipsum",
        "ingredients" => [
            "płatki owsiane",
            "wiórki kokosowe",
            "płatki migdałowe",
            "miód",
            "masło klarowane",
        ],
        "tags" => [
            "śniadanie",
            "vege"
        ],
    ],
    [
        "id" => 3,
        "title" => "PLACKI TWAROGOWE",
        "thumbnail" => "",
        "preparation" => "
            1. Do miski włożyć twaróg, dodać jajko, rozgnieść praską.
            2. Dodać przesianą mąkę i wymieszać.
            3. Dodać erytrol/ksylitol, wymieszać.
            4. Rozgrzać olej na patelni.
            5. Ciasto nakładać na rozgrzaną patelnię (po 2 łyżki na placek).
            6. Zmniejszyć ogień, smażyć po 2-3 minuty, przewrócić na drugą stronę.
            7. Podawać posypane cynamonem, skropione syropem klonowym i z borówkami.
        ",
        "ingredients" => [
            "twaróg chudy",
            "jajko",
            "ksylitol",
            "mąka orkiszow",
            "olej",
            "cynamon",
            "borówki",
            "syrop klonowy"
        ],
        "tags" => [
            "śniadanie",
            "wegetariańskie"
        ],
    ],
    [
        "id" => 4,
        "title" => "NALEŚNIKI Z NUTĄ POMARAŃCZY",
        "thumbnail" => "",
        "preparation" => "
            1. Masło roztopić.
            2. Jajko, śmietankę i sól ubić trzepaczką.
            3. Rozgrzać patelnię i usmażyć 3 cieniutkie naleśniki.
            4. Wymieszać sos.
            5. Każdy z usmażonych naleśników złoyć na ćwiartkę, ułożyć na patelni, wlać sos. Gotować chwilę do wchłonięci.
            5. Podsmażyć pestki.
            6. Naleśniki podsawać posypane cukrem i pestkami.
        ",
        "ingredients" => [
            "jajko",
            "śmietanka",
            "sól",
            "masło klarowane",
            "cukier puder",
            "pestki słonecznika",
            "cynamon",
            "sok z pomarańczy",
            "skórka z pomarańczy",
        ],
        "tags" => [
            "śniadanie",
            "wegetariańskie"
        ],
    ],
];

final class RecipeQueryRepository
{
    public function __construct(
        private readonly PDO $pdo,
    ) {}
    /**
     * @return GetRecipesResponse[]
     */
    public function getRecipes(): array
    {
        return array_map(function ($item) {
            return new GetRecipesResponse($item['id'], $item['title'], $item['thumbnail'], $item['tags']);
        }, RECIPES);
    }

    public function getRecipe(int $id): ?GetRecipeResponse
    {
        foreach (RECIPES as $recipe) {
            if ($recipe['id'] === $id) {
                return new GetRecipeResponse(
                    $recipe['id'],
                    $recipe['title'],
                    $recipe['thumbnail'],
                    $recipe['preparation'],
                    $recipe['ingredients'],
                    $recipe['tags'],
                );
            }
        }

        return null;
    }
}