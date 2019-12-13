<?php
$name = $_POST["name"];
$pdo = new PDO("mysql:host=vulner_mysql;dbname=vulner;", "root", "root", [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]);
$stm = $pdo->query("insert into users(name) values('$name')");
try {
    $stm->fetchAll();
} catch (Exception $exception) { }
echo htmlspecialchars($_POST["password"], ENT_QUOTES);
