<?php

function maskify(string $input): string
{
    if (empty($input)) {
        return "";
    }

    if (strlen($input) < 6) {
        return $input;
    }

    for ($i = 1; $i < strlen($input) - 4; $i++) {
        if (ctype_digit($input[$i])) {
            $input[$i] = '#';
        }
    }

    return $input;
}

$examples = [
    '1234567890123456',
    'A1B2C3D4E5F6G7H8',
    '12345',
    '',
    '1-2345-6789-0123-4567',
    'ABCD1234EFGH5678',
];

foreach ($examples as $example) {
    echo "Original: $example\nMasked:   " . maskify($example) . "\n\n";
}