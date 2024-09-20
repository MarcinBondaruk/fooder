<?php

function getConnection(): PDO {
    $username = 'postgres';
    $password = 'postgres';
    $host = "localhost";
    $port = 5432;
    $dbname = "fooder";
    $sslmode = "disable";

    $dsn = "pgsql:host=$host;port=$port;dbname=$dbname;sslmode=$sslmode";

    try {
        $pdo = new PDO($dsn, $username, $password);

        echo "connection established\n";
        return $pdo;
    } catch (\Throwable $e) {
        echo "error occured " . $e->getMessage() . "\n";
        exit();
    }
}

function createUsersTable(): void {

    try {
        $sql = <<<'SQL'
            CREATE TABLE IF NOT EXISTS users (
                id SERIAL PRIMARY KEY,
                username VARCHAR(128) NOT NULL,
                email VARCHAR(320) NOT NULL
            );
        SQL;

        $conn = getConnection();

        $conn->exec($sql);
        echo "users table created\n";
    } catch (\Throwable $e) {
        echo "error occured " . $e->getMessage() . "\n";
        exit();
    }
}

function insertUser(string $username, string $email): int {
    $sql = <<<'SQL'
        INSERT INTO users (username, email)
        VALUES (:username, :email)
        RETURNING id;
    SQL;
    $conn = getConnection();

    $stmt = $conn->prepare($sql);
    $stmt->execute(['username' => $username, 'email' => $email]);

    return $stmt->fetchColumn();
}

createUsersTable();
// echo insertUser('test', 'test');

var_dump($_REQUEST);
var_dump($_SERVER);
var_dump($_GET);