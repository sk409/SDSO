<?php
$name = $_GET["name"];
$pdo = new PDO("mysql:host=vulner_mysql;dbname=vulner;", "root", "root", [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]);
$stm = $pdo->query("select * from users where name=$name");
$stm->fetchAll();
echo htmlspecialchars($_GET["password"], ENT_QUOTES);
